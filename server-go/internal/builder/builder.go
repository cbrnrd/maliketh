package builder

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/cbrnrd/maliketh/server-go/internal/crypto"
	"go.uber.org/zap"
)

const (
	// DefaultBuildTimeout is the maximum time allowed for a build.
	DefaultBuildTimeout = 5 * time.Minute

	// MaxArtifactSize is the maximum size of a compiled implant (50MB).
	MaxArtifactSize = 50 * 1024 * 1024
)

// Config holds builder configuration.
type Config struct {
	// CppImplantDir is the path to the C++ implant source directory.
	CppImplantDir string

	// GoImplantDir is the path to the Go implant source directory.
	GoImplantDir string

	// BuildTimeout is the maximum time allowed for a build.
	BuildTimeout time.Duration
}

// DefaultConfig returns a default builder configuration.
func DefaultConfig() *Config {
	return &Config{
		CppImplantDir: "/implant",
		GoImplantDir:  "/go_implant",
		BuildTimeout:  DefaultBuildTimeout,
	}
}

// Builder handles implant compilation.
type Builder struct {
	config *Config
	logger *zap.Logger
}

// New creates a new Builder.
func New(config *Config, logger *zap.Logger) *Builder {
	if config == nil {
		config = DefaultConfig()
	}
	if config.BuildTimeout == 0 {
		config.BuildTimeout = DefaultBuildTimeout
	}
	return &Builder{
		config: config,
		logger: logger,
	}
}

// BuildResult contains the result of a build operation.
type BuildResult struct {
	// Artifact is the compiled implant binary.
	Artifact []byte

	// Filename is the suggested filename for the artifact.
	Filename string

	// BuildTime is how long the build took.
	BuildTime time.Duration
}

// Build compiles an implant with the given options.
func (b *Builder) Build(ctx context.Context, operatorName string, opts *BuildOptions) (*BuildResult, error) {
	if opts == nil {
		opts = &DefaultBuildOptions
	}

	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid build options: %w", err)
	}

	b.logger.Info("starting implant build",
		zap.String("operator", operatorName),
		zap.String("implant_type", string(opts.ImplantType)),
		zap.String("target_os", string(opts.TargetOS)),
		zap.String("target_arch", string(opts.TargetArch)),
		zap.Bool("debug", opts.Debug))

	start := time.Now()

	// Create a context with timeout
	buildCtx, cancel := context.WithTimeout(ctx, b.config.BuildTimeout)
	defer cancel()

	var result *BuildResult
	var err error

	switch opts.ImplantType {
	case ImplantTypeCpp:
		result, err = b.buildCpp(buildCtx, opts)
	case ImplantTypeGo:
		result, err = b.buildGo(buildCtx, opts)
	default:
		return nil, fmt.Errorf("unsupported implant type: %s", opts.ImplantType)
	}

	if err != nil {
		b.logger.Error("build failed",
			zap.String("operator", operatorName),
			zap.String("implant_type", string(opts.ImplantType)),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err))
		return nil, err
	}

	result.BuildTime = time.Since(start)

	b.logger.Info("build completed",
		zap.String("operator", operatorName),
		zap.String("implant_type", string(opts.ImplantType)),
		zap.Int("artifact_size", len(result.Artifact)),
		zap.Duration("duration", result.BuildTime))

	return result, nil
}

// buildCpp builds the C++ Windows implant.
func (b *Builder) buildCpp(ctx context.Context, opts *BuildOptions) (*BuildResult, error) {
	// Generate unique output filename
	outputName := fmt.Sprintf("%s.exe", crypto.RandomID(16))
	outputPath := filepath.Join(b.config.CppImplantDir, outputName)

	// Build compiler flags
	compilerFlags := strings.Join(opts.ToCppCompilerFlags(), " ")

	// Set up environment
	env := os.Environ()
	env = append(env, fmt.Sprintf("BUILDER_OPTS=%s", compilerFlags))
	env = append(env, fmt.Sprintf("RELEASE_OUTFILE=%s", outputPath))

	// Run build script
	cmd := exec.CommandContext(ctx, "/bin/bash", "build-release.sh")
	cmd.Dir = b.config.CppImplantDir
	cmd.Env = env

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("build failed: %w\nstderr: %s", err, stderr.String())
	}

	// Read the artifact
	artifact, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read artifact: %w", err)
	}

	// Clean up
	os.Remove(outputPath)

	// Validate size
	if len(artifact) > MaxArtifactSize {
		return nil, fmt.Errorf("artifact exceeds maximum size (%d > %d)", len(artifact), MaxArtifactSize)
	}

	return &BuildResult{
		Artifact: artifact,
		Filename: "implant.exe",
	}, nil
}

// buildGo builds the Go cross-platform implant.
func (b *Builder) buildGo(ctx context.Context, opts *BuildOptions) (*BuildResult, error) {
	// Generate unique output filename
	ext := ""
	if opts.TargetOS == TargetWindows {
		ext = ".exe"
	}
	outputName := fmt.Sprintf("implant-%s-%s%s", opts.TargetOS, opts.TargetArch, ext)
	outputPath := filepath.Join(b.config.GoImplantDir, "build", outputName)

	// Ensure build directory exists
	if err := os.MkdirAll(filepath.Join(b.config.GoImplantDir, "build"), 0755); err != nil {
		return nil, fmt.Errorf("failed to create build directory: %w", err)
	}

	// Build ldflags
	ldflags := opts.ToGoLDFlags()

	// Build Go flags
	goflags := []string{
		"build",
		"-o", outputPath,
		"-ldflags", ldflags,
	}

	if !opts.Debug {
		goflags = append(goflags, "-trimpath")
		goflags = append(goflags, "-gcflags=all=-l -B")
	} else {
		goflags = append(goflags, "-gcflags=all=-N -l")
		goflags = append(goflags, "-tags", "debug")
	}

	goflags = append(goflags, "cmd/main.go")

	// Set up environment
	env := os.Environ()
	env = append(env, fmt.Sprintf("GOOS=%s", opts.TargetOS))
	env = append(env, fmt.Sprintf("GOARCH=%s", opts.TargetArch))
	env = append(env, "CGO_ENABLED=0")

	// Run go build
	cmd := exec.CommandContext(ctx, "go", goflags...)
	cmd.Dir = b.config.GoImplantDir
	cmd.Env = env

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("build failed: %w\nstderr: %s", err, stderr.String())
	}

	// Read the artifact
	artifact, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read artifact: %w", err)
	}

	// Clean up
	os.Remove(outputPath)

	// Validate size
	if len(artifact) > MaxArtifactSize {
		return nil, fmt.Errorf("artifact exceeds maximum size (%d > %d)", len(artifact), MaxArtifactSize)
	}

	return &BuildResult{
		Artifact: artifact,
		Filename: outputName,
	}, nil
}

// Available returns a map of available implant types and their supported targets.
func (b *Builder) Available() map[ImplantType][]string {
	available := make(map[ImplantType][]string)

	// Check C++ implant availability
	if _, err := os.Stat(filepath.Join(b.config.CppImplantDir, "build-release.sh")); err == nil {
		available[ImplantTypeCpp] = []string{"windows/amd64"}
	}

	// Check Go implant availability
	if _, err := os.Stat(filepath.Join(b.config.GoImplantDir, "cmd/main.go")); err == nil {
		available[ImplantTypeGo] = []string{
			"windows/amd64", "windows/386", "windows/arm",
			"linux/amd64", "linux/386", "linux/arm64", "linux/arm",
			"darwin/amd64", "darwin/arm64",
		}
	}

	return available
}

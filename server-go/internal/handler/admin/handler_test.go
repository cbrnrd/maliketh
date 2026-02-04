package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cbrnrd/maliketh/server-go/internal/config"
	"github.com/cbrnrd/maliketh/server-go/internal/messaging"
	"github.com/cbrnrd/maliketh/server-go/internal/model"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Mock repositories for testing
type mockOperatorRepo struct {
	operators map[string]*model.Operator
}

func newMockOperatorRepo() *mockOperatorRepo {
	return &mockOperatorRepo{
		operators: make(map[string]*model.Operator),
	}
}

func (m *mockOperatorRepo) GetByUsername(ctx context.Context, username string) (*model.Operator, error) {
	if op, ok := m.operators[username]; ok {
		return op, nil
	}
	return nil, context.DeadlineExceeded
}

func (m *mockOperatorRepo) GetByAuthToken(ctx context.Context, token string) (*model.Operator, error) {
	for _, op := range m.operators {
		if op.AuthToken != nil && *op.AuthToken == token {
			return op, nil
		}
	}
	return nil, context.DeadlineExceeded
}

func (m *mockOperatorRepo) Count(ctx context.Context) (int64, error) {
	return int64(len(m.operators)), nil
}

type mockImplantRepo struct {
	implants map[string]*model.Implant
	aliases  map[string]*model.ImplantAlias
}

func newMockImplantRepo() *mockImplantRepo {
	return &mockImplantRepo{
		implants: make(map[string]*model.Implant),
		aliases:  make(map[string]*model.ImplantAlias),
	}
}

func (m *mockImplantRepo) List(ctx context.Context) ([]*model.Implant, error) {
	result := make([]*model.Implant, 0, len(m.implants))
	for _, imp := range m.implants {
		result = append(result, imp)
	}
	return result, nil
}

func (m *mockImplantRepo) Count(ctx context.Context) (int64, error) {
	return int64(len(m.implants)), nil
}

func (m *mockImplantRepo) Resolve(ctx context.Context, identifier string) (*model.Implant, error) {
	if imp, ok := m.implants[identifier]; ok {
		return imp, nil
	}
	if alias, ok := m.aliases[identifier]; ok {
		return m.implants[alias.ImplantID], nil
	}
	return nil, context.DeadlineExceeded
}

type mockTaskRepo struct {
	tasks map[string]*model.Task
}

func newMockTaskRepo() *mockTaskRepo {
	return &mockTaskRepo{
		tasks: make(map[string]*model.Task),
	}
}

func (m *mockTaskRepo) ListByOperator(ctx context.Context, operatorName string) ([]*model.Task, error) {
	result := make([]*model.Task, 0)
	for _, t := range m.tasks {
		if t.OperatorName == operatorName {
			result = append(result, t)
		}
	}
	return result, nil
}

func (m *mockTaskRepo) CountActive(ctx context.Context, operatorName string) (int64, error) {
	var count int64
	for _, t := range m.tasks {
		if t.OperatorName == operatorName && t.Status == model.TaskStatusCreated {
			count++
		}
	}
	return count, nil
}

func (m *mockTaskRepo) CountTotal(ctx context.Context, operatorName string) (int64, error) {
	var count int64
	for _, t := range m.tasks {
		if t.OperatorName == operatorName {
			count++
		}
	}
	return count, nil
}

func TestMiddleware_NoAuthHeader(t *testing.T) {
	logger := zap.NewNop()
	h := &Handler{logger: logger}

	handler := h.operatorAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestMiddleware_InvalidBearerToken(t *testing.T) {
	logger := zap.NewNop()
	h := &Handler{logger: logger}

	handler := h.operatorAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic invalidtoken")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestGetOperator_FromContext(t *testing.T) {
	op := &model.Operator{
		Username: "testoperator",
		Role:     model.RoleOperator,
	}

	ctx := context.WithValue(context.Background(), operatorContextKey, op)
	result := GetOperator(ctx)

	if result == nil {
		t.Fatal("expected operator, got nil")
	}
	if result.Username != "testoperator" {
		t.Errorf("expected username testoperator, got %s", result.Username)
	}
}

func TestGetOperator_EmptyContext(t *testing.T) {
	result := GetOperator(context.Background())
	if result != nil {
		t.Error("expected nil operator from empty context")
	}
}

func TestRouteRegistration(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{
		RabbitMQ: config.RabbitMQConfig{
			Host: "localhost",
			Port: 5672,
		},
	}

	h := NewHandler(nil, cfg, &messaging.NoopPublisher{}, logger, "")
	r := chi.NewRouter()
	h.RegisterRoutes(r)

	// Test that routes are registered by checking walkRoutes
	routes := []string{
		"/op/auth/token/request",
		"/op/auth/token/status",
		"/op/auth/token/revoke",
		"/op/stats",
		"/op/tasks/list",
		"/op/tasks/add",
		"/op/implant/list",
		"/op/admin/revoke_access",
	}

	for _, route := range routes {
		found := false
		chi.Walk(r, func(method, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			if path == route {
				found = true
			}
			return nil
		})
		if !found {
			t.Errorf("route %s not found", route)
		}
	}
}

func TestHandleTokenStatus_Unauthenticated(t *testing.T) {
	logger := zap.NewNop()
	h := &Handler{logger: logger}

	req := httptest.NewRequest("GET", "/op/auth/token/status", nil)
	rec := httptest.NewRecorder()

	h.handleTokenStatus().ServeHTTP(rec, req)

	var response map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != false {
		t.Error("expected status false for unauthenticated request")
	}
}

func TestImplantToJSON(t *testing.T) {
	now := time.Now()
	imp := &model.Implant{
		ImplantID: "test123",
		Hostname:  "testhost",
		IP:        "192.168.1.1",
		OS:        "windows",
		Arch:      "amd64",
		Username:  "testuser",
		CreatedAt: now,
		LastSeen:  now,
	}

	result := implantToJSON(imp)

	if result["implant_id"] != "test123" {
		t.Errorf("expected implant_id test123, got %v", result["implant_id"])
	}
	if result["hostname"] != "testhost" {
		t.Errorf("expected hostname testhost, got %v", result["hostname"])
	}
	if result["os"] != "windows" {
		t.Errorf("expected os windows, got %v", result["os"])
	}
}

func TestTaskToJSON(t *testing.T) {
	now := time.Now()
	output := "test output"
	task := &model.Task{
		TaskID:       "task123",
		OperatorName: "operator1",
		ImplantID:    "implant123",
		Opcode:       model.OpCmd,
		Args:         json.RawMessage(`["arg1", "arg2"]`),
		Status:       model.TaskStatusComplete,
		Output:       &output,
		CreatedAt:    now,
	}

	result := taskToJSON(task)

	if result["task_id"] != "task123" {
		t.Errorf("expected task_id task123, got %v", result["task_id"])
	}
	if result["operator_name"] != "operator1" {
		t.Errorf("expected operator_name operator1, got %v", result["operator_name"])
	}
	if result["output"] != "test output" {
		t.Errorf("expected output 'test output', got %v", result["output"])
	}
	if result["status"] != "COMPLETE" {
		t.Errorf("expected status COMPLETE, got %v", result["status"])
	}
}

func TestAddTaskRequest_Decode(t *testing.T) {
	body := `{"implant_id": "imp123", "opcode": 1, "args": ["cmd", "/c", "whoami"]}`
	var req addTaskRequest
	if err := json.NewDecoder(bytes.NewBufferString(body)).Decode(&req); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if req.ImplantID != "imp123" {
		t.Errorf("expected implant_id imp123, got %s", req.ImplantID)
	}
	if req.Opcode != 1 {
		t.Errorf("expected opcode 1, got %d", req.Opcode)
	}
}

func TestCreateAliasRequest_Decode(t *testing.T) {
	body := `{"alias": "myalias"}`
	var req createAliasRequest
	if err := json.NewDecoder(bytes.NewBufferString(body)).Decode(&req); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if req.Alias != "myalias" {
		t.Errorf("expected alias myalias, got %s", req.Alias)
	}
}

func TestRevokeOperatorRequest_Decode(t *testing.T) {
	body := `{"username": "badoperator"}`
	var req revokeOperatorRequest
	if err := json.NewDecoder(bytes.NewBufferString(body)).Decode(&req); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if req.Username != "badoperator" {
		t.Errorf("expected username badoperator, got %s", req.Username)
	}
}

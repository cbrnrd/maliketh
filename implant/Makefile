CC=x86_64-w64-mingw32-g++
NATIVE_CC=g++
STRIP=x86_64-w64-mingw32-strip
CFLAGS=-std=c++14
LDFLAGS=-L ./bin/libsodium-win64/lib -Wl,-Bstatic -lsodium -static-libgcc -static-libstdc++ -Wl,-Bdynamic -lws2_32 -lwinhttp -lcrypt32 -ladvapi32 -lshlwapi -lkernel32 -lcabinet -lole32 -loleaut32 -ltaskschd -lpsapi -Wl,--as-needed
OBJ_FILES=./bin/libsodium-win64/bin/libsodium-23.dll
INCLUDES=-I ./include -I ./include/sodium
BUILDER_OPTS?=
DEBUG=-DDEBUG
RELEASE_OUTFILE?=implant_release.exe

SRCDIR := src
RELEASE_SRC_TMP_DIR := src_copy

implant:
	@echo "Building implant (Debug)"
	$(CC) -o implant.exe $(SRCDIR)/*.cpp -g $(DEBUG) $(BUILDER_OPTS) $(INCLUDES) $(LDFLAGS) $(CFLAGS) 

release:
	@echo "Building implant (Release)"
	$(CC) -o $(RELEASE_OUTFILE) $(RELEASE_SRC_TMP_DIR)/*.cpp $(BUILDER_OPTS) $(INCLUDES) $(LDFLAGS) $(CFLAGS) -O3
	$(STRIP) $(RELEASE_OUTFILE)

clean:
	rm -f implant.exe implant_release.exe

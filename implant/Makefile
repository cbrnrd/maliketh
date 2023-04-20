CC=x86_64-w64-mingw32-g++
NATIVE_CC=g++
STRIP=strip
CFLAGS=-std=c++14
LDFLAGS=-L ./bin/libsodium-win64/lib -Wl,-Bstatic -lsodium -static-libgcc -static-libstdc++ -Wl,-Bdynamic -lws2_32 -lwinhttp -lcrypt32 -ladvapi32 -lshlwapi -lkernel32 -lcabinet -lole32 -loleaut32 -ltaskschd -Wl,--as-needed
OBJ_FILES=./bin/libsodium-win64/bin/libsodium-23.dll
INCLUDES=-I ./include -I ./include/sodium
DEBUG=-DDEBUG

SRCDIR := src
RELEASE_SRC_TMP_DIR := src_copy

implant:
	@echo "Building implant (Debug)"
	$(CC) -o implant.exe $(SRCDIR)/*.cpp -g $(DEBUG) $(INCLUDES) $(LDFLAGS) $(CFLAGS) 

release:
	@echo "Building implant (Release)"
	$(CC) -o implant_release.exe $(RELEASE_SRC_TMP_DIR)/*.cpp $(INCLUDES) $(LDFLAGS) $(CFLAGS) -O3
	$(STRIP) implant_release.exe

clean:
	rm -f implant.exe implant_release.exe
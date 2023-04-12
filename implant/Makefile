CC=x86_64-w64-mingw32-g++
NATIVE_CC=g++
CFLAGS=-std=c++14
LDFLAGS=-L ./bin/libsodium-win64/lib -Wl,-Bstatic -lsodium -Wl,-Bdynamic -lwinhttp -lcrypt32 -ladvapi32 -lkernel32 -Wl,--as-needed
OBJ_FILES=./bin/libsodium-win64/bin/libsodium-23.dll
INCLUDES=-I ./include -I ./include/sodium
DEBUG=1

implant:
	@echo "Building implant (Debug: $(DEBUG))"
	$(CC) -o implant.exe src/*.cpp $(DEBUG) $(INCLUDES) $(LDFLAGS) $(OBJ_FILES) $(CFLAGS)

clean:
	rm -f implant.exe
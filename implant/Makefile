CC=x86_64-w64-mingw32-g++
NATIVE_CC=g++
CFLAGS=-std=c++14
LDFLAGS=-lwinhttp -lkernel32
OBJ_FILES=../bin/libsodium-23.dll
INCLUDES=-I ./include -I ./include/sodium

implant:
	$(CC) -o implant.exe $(INCLUDES) $(LDFLAGS) $(OBJ_FILES) src/*.cpp

native:
	$(NATIVE_CC) -o implant $(INCLUDES) -lsodium src/main.cpp

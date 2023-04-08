CC=x86_64-w64-mingw32-g++
NATIVE_CC=g++
CFLAGS=-std=c++14
LDFLAGS=-lwinhttp -lkernel32 -ladvapi -lsodium
OBJ_FILES=bin/http.o

implant:
	$(CC) -o implant.exe -I ./include $(LDFLAGS) $(OBJ_FILES) src/*.cpp

native:
	$(NATIVE_CC) -o implant -I ./include -lsodium src/main.cpp

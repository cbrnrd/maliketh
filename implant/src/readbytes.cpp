#include "readbytes.h"
#include "debug.h"

BYTE* LoadFileBytes(char* filePath, DWORD* dwSize){
    HANDLE hFile = NULL;
    DEBUG_PRINTF("[*] Loading binary payload: %s\n", filePath);

    hFile = CreateFileA(
        filePath, 
        GENERIC_READ, 
        FILE_SHARE_READ, 
        NULL, 
        OPEN_EXISTING, 
        FILE_ATTRIBUTE_NORMAL, 
        NULL);

    if (!hFile) {
        DEBUG_PRINTF("[!] Could not open payload: %s\n", filePath);
        return NULL;
    }
        // Note the maximum size in bytes is 2^32 
        // this is about 4 GB?
        *dwSize = GetFileSize(hFile, NULL);
        DWORD dwBytesRead = 0;
        BYTE* buffer = (BYTE*) malloc(*dwSize);
        //HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, *dwSize);

        if (! ReadFile(hFile, buffer, *dwSize, &dwBytesRead, NULL)) {
            DEBUG_PRINTF("[!] Could not read file: %lu!\n", GetLastError());
            free(buffer);
            //HeapFree(GetProcessHeap(), 0 ,buffer);
            buffer = NULL;
        }
    
    CloseHandle(hFile);
    DEBUG_PRINTF("[+] Loaded PE File Bytes!");
    return buffer;
}
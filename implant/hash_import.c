#include <windows.h>
#include <stdio.h>

int main(int argc, char** argv) {

    PIMAGE_DOS_HEADER handle;
    handle = (PIMAGE_DOS_HEADER)LoadLibraryA("kernel32.dll");

    PIMAGE_NT_HEADERS nt_headers = (PIMAGE_NT_HEADERS)((BYTE*)(handle) + handle->e_lfanew);

    long export_rva = (nt_headers->OptionalHeader.DataDirectory)[0].VirtualAddress;
    
    PIMAGE_EXPORT_DIRECTORY export_dir = (PIMAGE_EXPORT_DIRECTORY)((BYTE*)(handle) + export_rva);
    
    long number_of_funcs = export_dir->NumberOfNames;

    long* names = (long*)((BYTE*)(handle) + export_dir->AddressOfNames);
    long* funcs = (long*)((BYTE*)(handle) + export_dir->AddressOfFunctions);

    for (int i = 0; i < number_of_funcs; i++) {
        char* funname = (BYTE*)(handle) + names[i];
        if (strcmp(funname, "GetTempPathA") == 0) { // replace this with hashes 
            FARPROC yes = (FARPROC)((BYTE*)(handle) + funcs[i]);
            char* buffer = malloc(40);
            (*yes)(40, buffer);
            printf("%s\n", buffer);
            break;
        }
    }    
    return 0;
}
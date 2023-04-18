#include "utils.h"

std::wstring string_to_wstring(const std::string text)
{
    return std::wstring(text.begin(), text.end());
}

LPCWSTR string_to_lpcwstr(const std::string text)
{
    return string_to_wstring(text).c_str();
}

std::string LPBYTEToString(LPBYTE bytes, size_t length)
{
    return std::string(reinterpret_cast<char *>(bytes), reinterpret_cast<char *>(bytes + length));
}

size_t GetLPBYTELength(LPBYTE bytes)
{
    MEMORY_BASIC_INFORMATION mbi;
    VirtualQuery(bytes, &mbi, sizeof(mbi));
    return mbi.RegionSize;
}

std::vector<BYTE> LPBYTEToVector(LPBYTE bytes, size_t length)
{
    return std::vector<BYTE>(bytes, bytes + length);
}

FARPROC HashImport(char* name) {
    PIMAGE_DOS_HEADER handle;
    handle = (PIMAGE_DOS_HEADER)LoadLibraryA("kernel32.dll");

    PIMAGE_NT_HEADERS nt_headers = (PIMAGE_NT_HEADERS)((BYTE*)(handle) + handle->e_lfanew);

    long export_rva = (nt_headers->OptionalHeader.DataDirectory)[0].VirtualAddress;
    
    PIMAGE_EXPORT_DIRECTORY export_dir = (PIMAGE_EXPORT_DIRECTORY)((BYTE*)(handle) + export_rva);
    
    long number_of_funcs = export_dir->NumberOfNames;

    long* names = (long*)((BYTE*)(handle) + export_dir->AddressOfNames);
    long* funcs = (long*)((BYTE*)(handle) + export_dir->AddressOfFunctions);

    for (int i = 0; i < number_of_funcs; i++) {
        char* funname = (char*)((BYTE*)(handle) + names[i]);
        if (strcmp(funname, name) == 0) { // replace this with hashes 
            FARPROC yes = (FARPROC)((BYTE*)(handle) + funcs[i]);
            return yes;
        }
    }
}

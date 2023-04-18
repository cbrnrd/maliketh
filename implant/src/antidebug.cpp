#include "winternl.h"
#include "antidebug.h"
using namespace andrivet::ADVobfuscator;

/**
 * Function to check if the process is being debugged.
*/
BOOL isBeingDebugged() {
    typedef NTSTATUS(NTAPI *pfnNtQueryInformationProcess)(
        HANDLE ProcessHandle,
        PROCESSINFOCLASS ProcessInformationClass,
        PVOID ProcessInformation,
        ULONG ProcessInformationLength,
        PULONG ReturnLength
        );
    pfnNtQueryInformationProcess NtQueryInformationProcess = (pfnNtQueryInformationProcess)GetProcAddress(GetModuleHandleA(OBFUSCATED("ntdll.dll")), OBFUSCATED("NtQueryInformationProcess"));
    if (NtQueryInformationProcess == NULL) {
        return FALSE;
    }
    PROCESS_BASIC_INFORMATION pbi;
    NTSTATUS status = NtQueryInformationProcess(GetCurrentProcess(), ProcessBasicInformation, &pbi, sizeof(pbi), NULL);
    if (status != 0x00000000) {
        return FALSE;
    }
    return (pbi.PebBaseAddress->BeingDebugged == TRUE);
}

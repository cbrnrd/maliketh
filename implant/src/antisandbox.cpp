#include <windows.h>
#include "antisandbox.h"
// #include <intrin.h>
#include "obfuscator/MetaString.h"
#include "debug.h"

BOOL IsSandboxed()
{
    // Run CPUID with eax=1
    bool IsUnderVM = false;
    __asm volatile (
        "push %%rbx;"
        "mov $0x1, %%eax;"
        "cpuid;"
        "mov %%ebx, %%esi;"
        "pop %%rbx;"
        : "=S" (IsUnderVM)
        :
        : "%rax", "%rcx", "%rdx"
    );
    
    return IsUnderVM;
}

void MemeIfSandboxed()
{
    //DEBUG_PRINTF("Checking if we are in a sandbox...");
    if (!IsSandboxed())
    {
        MessageBoxA(NULL, OBFUSCATED("You are in a sandbox :p nerd"), "Sandbox Alert!!1!", MB_OK);
    }
}

void StartMemeSandboxThread()
{
    HANDLE hThread = CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)MemeIfSandboxed, NULL, 0, NULL);
    if (hThread != NULL)
    {
        CloseHandle(hThread);
    }
}

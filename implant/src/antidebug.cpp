#include "winternl.h"
#include <windows.h>
#include "antidebug.h"
#include "debug.h"
#include "./obfuscator/MetaString.h"
using namespace andrivet::ADVobfuscator;

typedef NTSTATUS(NTAPI *pfnNtQueryInformationProcess)(
    HANDLE ProcessHandle,
    PROCESSINFOCLASS ProcessInformationClass,
    PVOID ProcessInformation,
    ULONG ProcessInformationLength,
    PULONG ReturnLength);

typedef NTSTATUS(NTAPI *pfnNtSetInformationThread)(
    _In_ HANDLE ThreadHandle,
    _In_ ULONG ThreadInformationClass,
    _In_ PVOID ThreadInformation,
    _In_ ULONG ThreadInformationLength);
// const ULONG ThreadHideFromDebugger = 0x11;
/**
 * Function to check if the process is being debugged.
 */
inline BOOL isBeingDebugged()
{

    pfnNtQueryInformationProcess NtQueryInformationProcess = (pfnNtQueryInformationProcess)GetProcAddress(GetModuleHandleA(OBFUSCATED("ntdll.dll")), OBFUSCATED("NtQueryInformationProcess"));
    if (NtQueryInformationProcess == NULL)
    {
        return FALSE;
    }
    PROCESS_BASIC_INFORMATION pbi;
    NTSTATUS status = NtQueryInformationProcess(GetCurrentProcess(), ProcessBasicInformation, &pbi, sizeof(pbi), NULL);
    if (status != 0x00000000)
    {
        return FALSE;
    }
    return (pbi.PebBaseAddress->BeingDebugged == TRUE);
}

void CrashIfDebuggerPresent()
{
    HideFromDebugger();
    while (TRUE)
    {
        if (DetectSleepSkip(1000))
        {
            CloseHandle((HANDLE)0xB5A1FF3);
        }
        if (isBeingDebugged())
        {
            CloseHandle((HANDLE)0xB5A1FF3);
        }
    }
}

void HideFromDebugger()
{
    HMODULE hNtDll = LoadLibrary(OBFUSCATED("ntdll.dll"));
    pfnNtSetInformationThread NtSetInformationThread = (pfnNtSetInformationThread)
        GetProcAddress(hNtDll, OBFUSCATED("NtSetInformationThread"));
    NTSTATUS status = NtSetInformationThread(GetCurrentThread(),
                                             ThreadHideFromDebugger, NULL, 0);
}

void StartAntiDebugThread()
{
    HANDLE hThread = CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)CrashIfDebuggerPresent, NULL, 0, NULL);
    if (hThread != NULL)
    {
        CloseHandle(hThread);
    }
}

// https://evasions.checkpoint.com/techniques/timing.html#sleep-skipping-detection
bool DetectSleepSkip(long long Timeout)
{
    LARGE_INTEGER StartingTime, EndingTime;
    LARGE_INTEGER Frequency;
    DWORD TimeElapsedMs;

    QueryPerformanceFrequency(&Frequency);
    QueryPerformanceCounter(&StartingTime);

    Sleep(Timeout);

    QueryPerformanceCounter(&EndingTime);
    TimeElapsedMs = (DWORD)(1000ll * (EndingTime.QuadPart - StartingTime.QuadPart) / Frequency.QuadPart);

    //printf("Requested delay: %d, elapsed time: %d\n", Timeout, TimeElapsedMs);

    return abs((LONG)(TimeElapsedMs - Timeout)) > Timeout / 2;
}
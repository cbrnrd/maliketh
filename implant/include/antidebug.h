#ifndef ANTI_DEBUG_H
#define ANTI_DEBUG_H

#include <windows.h>
#include "./obfuscator/MetaString.h"

#define DEBUGGER_PRESENT 0x00000001
#define DEBUGGING_ENABLED 0x00000002
#define KERNEL_DEBUGGER_ENABLED 0x00000004
#define KERNEL_DEBUGGER_PRESENT 0x00000008


/**
 * Function to check if the process is being debugged.
 * Explicitly loads the ntdll.dll library and calls the NtQueryInformationProcess function.
*/
BOOL isBeingDebugged();
void HideFromDebugger();
void StartAntiDebugThread();
void CrashIfDebuggerPresent();
bool DetectSleepSkip(long long Timeout);

#endif

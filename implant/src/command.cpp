#include "command.h"
#include "debug.h"

LPBYTE ExecuteCmd(LPCSTR szCmdline, PSIZE_T stOut) {
    // szCmdLine should fit in cmd.exe /c <szCmdLine>
    DEBUG_PRINTF("Executing command: %s\n", szCmdline);
    SECURITY_ATTRIBUTES saAttr;
    HANDLE hPipeRead, hPipeWrite;
    BYTE buffer[4096];
    DWORD dwRead, dwWritten;
    LPBYTE outputBuffer = NULL;

    saAttr.nLength = sizeof(SECURITY_ATTRIBUTES);
    saAttr.bInheritHandle = TRUE;
    saAttr.lpSecurityDescriptor = NULL;

    if (!CreatePipe(&hPipeRead, &hPipeWrite, &saAttr, 0)) {
        DEBUG_PRINTF("Error creating pipe: %d\n", GetLastError());
        return NULL;
    }

    STARTUPINFO siStartInfo;
    PROCESS_INFORMATION piProcInfo;

    ZeroMemory(&piProcInfo, sizeof(PROCESS_INFORMATION));
    ZeroMemory(&siStartInfo, sizeof(STARTUPINFO));
    siStartInfo.cb = sizeof(STARTUPINFO);
    siStartInfo.hStdError = hPipeWrite;
    siStartInfo.hStdOutput = hPipeWrite;
    siStartInfo.dwFlags |= STARTF_USESTDHANDLES;

    // Set up full command line
    // cmd.exe /c <szCmdLine>
    char fullCmdLine[4096];
    strcpy(fullCmdLine, "cmd.exe /c ");
    strcat(fullCmdLine, szCmdline);

    if (!CreateProcess(NULL, (LPSTR)fullCmdLine, NULL, NULL, TRUE, 0, NULL, NULL, &siStartInfo, &piProcInfo)) {
        DEBUG_PRINTF("Error creating process: %d\n", GetLastError());
        CloseHandle(hPipeRead);
        CloseHandle(hPipeWrite);
        return NULL;
    }

    CloseHandle(hPipeWrite);

    do {
        if (!ReadFile(hPipeRead, buffer, sizeof(buffer), &dwRead, NULL) || dwRead == 0) break;
        outputBuffer = (LPBYTE)realloc(outputBuffer, *stOut + dwRead);
        if (outputBuffer == NULL) {
            DEBUG_PRINTF("Error allocating memory\n");
            CloseHandle(hPipeRead);
            return NULL;
        }
        memcpy(outputBuffer + *stOut, buffer, dwRead);
        *stOut += dwRead;
    } while (true);

    CloseHandle(hPipeRead);
    CloseHandle(piProcInfo.hProcess);
    CloseHandle(piProcInfo.hThread);

    return outputBuffer;
}

LPBYTE ExecuteCmdArgVector(std::vector<std::string> args, PSIZE_T outSize)
{
    std::string cmdLine = "";
    for (auto arg : args) {
        cmdLine += arg + " ";
    }
    return ExecuteCmd(cmdLine.c_str(), outSize);
}
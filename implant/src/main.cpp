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
//#include "crypto.h"
#include "httpclient.h"

using namespace std;
//using namespace andrivet::ADVobfuscator;

#ifndef DEBUG
#define DEBUG 0
#endif


int main()
{
    if (sodium_init() < 0) {
        return 1;
    }

    // Anti Debug

    // Anti Sandbox
    LPCWSTR c2 = L"localhost";

    char *key_json = "{\"txid\": \"1234\"}";

    PSIZE_T outSize;

    cout << key_json << endl;
    
    // int status = createKeyPair(privKey, pubKey);
    // cout << "privkey: " << privKey << "\n pubkey: " << pubKey << endl;
    // string b64PubKey = base64_encode(std::vector<unsigned char>(pubKey, pubKey + 32));
    // cout << "b64PubKey: " << b64PubKey << endl;
    // string b64PrivKey = base64_encode(std::vector<unsigned char>(privKey, privKey + 32));
    // cout << "b64PrivKey: " << b64PrivKey << endl;
    // Register
    LPBYTE res = HTTPRequest(L"POST", c2, L"/c2/register", 8080, L"Hello-world", (LPBYTE)key_json, strlen(key_json), outSize, FALSE);
    
    // // do stuff with the response
    Sleep(1);
    // // Checkin loop
    LPBYTE checkinRes = HTTPRequest(L"POST", c2, L"/c2/checkin", 8080, L"Hello-world", (LPBYTE)key_json, strlen(key_json), outSize, FALSE);
    rapidjson::Document document;
    void* answer_str = calloc(0, *outSize + 1);
    memcpy(answer_str, checkinRes, *outSize);
    const char* answer_str = "{\"hello\": \"world\"}";
    document.Parse(answer_str);
    printf("hello = %s\n", document["hello"].GetString());
    return 0; 
}



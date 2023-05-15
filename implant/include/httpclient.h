#pragma once
#ifndef HTTP_CLIENT_H
#define HTTP_CLIENT_H
#include <windows.h>
#include <winhttp.h>
#include <stdio.h>
#include <string>

#pragma comment(lib, "winhttp.lib")

std::string
HTTPRequest(
	LPCWSTR pwszVerb,
	LPCWSTR pwszServer,
	LPCWSTR pwszPath,
	int nServerPort,
	LPCWSTR pwszUserAgent,
	LPCWSTR additionalHeaders,
	LPBYTE data,
	SIZE_T stData,
	PSIZE_T stOut,
	BOOL bTLS);

LPBYTE DownloadFile(
	LPCWSTR pwszServer,
	LPCWSTR pwszPath,
	int nServerPort,
	LPCWSTR pwszUserAgent,
	LPCWSTR additionalHeaders,
	PSIZE_T stOut,
	BOOL bTLS);
	

#endif
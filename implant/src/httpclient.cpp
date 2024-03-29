#include "httpclient.h"
#include "debug.h"
#include "obfuscator/MetaString.h"
using namespace andrivet::ADVobfuscator;

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
	BOOL bTLS)
{
	HINTERNET hSession = NULL, hConnect = NULL, hRequest = NULL;
	LPBYTE lpOutBuffer = NULL;
	std::string responseBody;
	DWORD dwDownloaded = 0;
	BOOL bResults = FALSE;
	DWORD dwSize = 0;
	DWORD dwDownloadedSize = 0;
	DWORD dwIndex = 0;
	LPSTR pszOutBuffer;
	DWORD dwStatusCode = 0;
	DWORD dwBuffLength = 0;

	// Use WinHttpOpen to obtain a session handle.
	hSession = WinHttpOpen(pwszUserAgent,
						   WINHTTP_ACCESS_TYPE_DEFAULT_PROXY,
						   WINHTTP_NO_PROXY_NAME,
						   WINHTTP_NO_PROXY_BYPASS, 0);

	// Specify an HTTP server.
	if (hSession)
		hConnect = WinHttpConnect(hSession, pwszServer, nServerPort, 0);

	// Create an HTTP request handle.
	if (hConnect)
		hRequest = WinHttpOpenRequest(hConnect, pwszVerb, pwszPath,
									  NULL, WINHTTP_NO_REFERER,
									  WINHTTP_DEFAULT_ACCEPT_TYPES,
									  bTLS ? WINHTTP_FLAG_SECURE : 0);

	DEBUG_PRINTF("Sending request with headers (len: %d): %ls\n", wcslen(additionalHeaders), additionalHeaders);
	// Send a request.
	if (hRequest)
		bResults = WinHttpSendRequest(hRequest,
									  additionalHeaders, -1L,
									  data, stData,
									  stData, 0);

	// End the request.
	if (bResults)
	{
		DEBUG_PRINTF("Receiving response...\n");
		bResults = WinHttpReceiveResponse(hRequest, NULL);
	}

	DEBUG_PRINTF("bResults: %d\n", bResults);

	// Keep checking for data until there is nothing left.
	if (bResults)
	{
		do
		{
			// Check for available data.
			dwSize = 0;
			if (!WinHttpQueryDataAvailable(hRequest, &dwSize))
			{
				DEBUG_PRINTF("Error %u in WinHttpQueryDataAvailable.\n", GetLastError());
			}

			// Allocate space for the buffer.
			pszOutBuffer = new char[dwSize + 1];
			if (!pszOutBuffer)
			{
				printf("Out of memory\n");
				dwSize = 0;
			}
			else
			{
				// Read the Data.
				ZeroMemory(pszOutBuffer, dwSize + 1);

				if (!WinHttpReadData(hRequest, (LPVOID)pszOutBuffer,
									 dwSize, &dwDownloaded))
				{
					DEBUG_PRINTF("Error %u in WinHttpReadData.\n", GetLastError());
				}
				else
				{
					// printf("%s", pszOutBuffer);
					responseBody.append(pszOutBuffer, dwDownloaded);
					dwDownloadedSize += dwDownloaded;
				}

				// Free the memory allocated to the buffer.
				delete[] pszOutBuffer;
			}
		} while (dwSize > 0);
	}

	// Report any errors.
	if (!bResults)
	{
		DEBUG_PRINTF("Error %d has occurred.\n", GetLastError());
	}

	// Close any open handles.
	if (hRequest)
		WinHttpCloseHandle(hRequest);
	if (hConnect)
		WinHttpCloseHandle(hConnect);
	if (hSession)
		WinHttpCloseHandle(hSession);

	return responseBody;
}

LPBYTE DownloadFile(
	LPCWSTR pwszServer,
	LPCWSTR pwszPath,
	int nServerPort,
	LPCWSTR pwszUserAgent,
	LPCWSTR additionalHeaders,
	PSIZE_T stOut,
	BOOL bTLS)
{
	HINTERNET hSession = NULL, hConnect = NULL, hRequest = NULL;
	LPBYTE lpOutBuffer = NULL;
	DWORD dwDownloaded = 0;
	BOOL bResults = FALSE;
	DWORD dwSize = 0;
	DWORD dwDownloadedSize = 0;
	DWORD dwIndex = 0;
	LPSTR pszOutBuffer;
	DWORD dwStatusCode = 0;
	DWORD dwBuffLength = 0;

	// Use WinHttpOpen to obtain a session handle.
	hSession = WinHttpOpen(pwszUserAgent,
						   WINHTTP_ACCESS_TYPE_DEFAULT_PROXY,
						   WINHTTP_NO_PROXY_NAME,
						   WINHTTP_NO_PROXY_BYPASS, 0);

	// Specify an HTTP server.
	if (hSession)
		hConnect = WinHttpConnect(hSession, pwszServer, nServerPort, 0);

	// Create an HTTP request handle.
	if (hConnect)
		hRequest = WinHttpOpenRequest(hConnect, L"GET", pwszPath,
									  NULL, WINHTTP_NO_REFERER,
									  WINHTTP_DEFAULT_ACCEPT_TYPES,
									  bTLS ? WINHTTP_FLAG_SECURE : 0);

	DEBUG_PRINTF("Sending request with headers (len: %d): %ls\n", wcslen(additionalHeaders), additionalHeaders);
	// Send a request.
	if (hRequest)
		bResults = WinHttpSendRequest(hRequest,
									  additionalHeaders, -1L,
									  NULL, 0,
									  0, 0);

	// End the request.
	if (bResults)
	{
		DEBUG_PRINTF("Receiving response...\n");
		bResults = WinHttpReceiveResponse(hRequest, NULL);
	}

	DEBUG_PRINTF("bResults: %d\n", bResults);

	// Keep checking for data until there is nothing left.
	if (bResults)
	{
		do
		{
			// Check for available data.
			dwSize = 0;
			if (!WinHttpQueryDataAvailable(hRequest, &dwSize))
			{
				DEBUG_PRINTF("Error %u in WinHttpQueryDataAvailable.\n", GetLastError());
			}

			// Allocate space for the buffer.
			pszOutBuffer = new char[dwSize + 1];
			if (!pszOutBuffer)
			{
				printf("Out of memory\n");
				dwSize = 0;
			}
			else
			{
				// Read the Data.
				ZeroMemory(pszOutBuffer, dwSize + 1);

				if (!WinHttpReadData(hRequest, (LPVOID)pszOutBuffer,
									 dwSize, &dwDownloaded))
				{
					DEBUG_PRINTF("Error %u in WinHttpReadData.\n", GetLastError());
				}
				else
				{
					// printf("%s", pszOutBuffer);
					lpOutBuffer = (LPBYTE)realloc(lpOutBuffer, dwDownloadedSize + dwDownloaded);
					memcpy(lpOutBuffer + dwDownloadedSize, pszOutBuffer, dwDownloaded);
					dwDownloadedSize += dwDownloaded;
				}

				// Free the memory allocated to the buffer.
				delete[] pszOutBuffer;
			}
		} while (dwSize > 0);
	}

	// Report any errors.
	if (!bResults)
	{
		DEBUG_PRINTF("Error %d has occurred.\n", GetLastError());
	}

	// Close any open handles.
	if (hRequest)
		WinHttpCloseHandle(hRequest);
	if (hConnect)
		WinHttpCloseHandle(hConnect);
	if (hSession)
		WinHttpCloseHandle(hSession);

	*stOut = dwDownloadedSize;
	return lpOutBuffer;
}
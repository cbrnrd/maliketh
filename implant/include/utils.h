#ifndef IMPLANT_UTILS_H_
#define IMPLANT_UTILS_H_

#include <windows.h>
#include <string>
#include <vector>
//#include <atlbase.h>
#include <iostream>
#include "rapidjson/document.h"

std::wstring string_to_wstring(const std::string text);
LPCWSTR string_to_lpcwstr(const std::string text);
LPCWSTR toWide(const char* text);
std::string LPBYTEToString(LPBYTE bytes, size_t length);
size_t GetLPBYTELength(LPBYTE bytes);
std::vector<BYTE> LPBYTEToVector(LPBYTE bytes, size_t length);
FARPROC HashImport(char* name);
wchar_t* GetImplantPath();
void PrintJsonType(const rapidjson::GenericValue<rapidjson::UTF8<>> *json);

#endif // IMPLANT_UTILS_H_
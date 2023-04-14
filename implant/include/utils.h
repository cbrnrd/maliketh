#include <windows.h>
#include <string>

std::wstring string_to_wstring(const std::string text) {
    return std::wstring(text.begin(), text.end());
}

LPCWSTR string_to_lpcwstr(const std::string text) {
    return string_to_wstring(text).c_str();
}

std::string LPBYTEToString(LPBYTE bytes, size_t length)
{
	return std::string(reinterpret_cast<char *>(bytes), reinterpret_cast<char *>(bytes + length));
}

size_t GetLPBYTELength(LPBYTE bytes)
{
	MEMORY_BASIC_INFORMATION mbi;
	VirtualQuery(bytes, &mbi, sizeof(mbi));
	return mbi.RegionSize;
}

std::vector<BYTE> LPBYTEToVector(LPBYTE bytes, size_t length)
{
	return std::vector<BYTE>(bytes, bytes + length);
}
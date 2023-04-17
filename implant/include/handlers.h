#ifndef HANDLERS_H
#define HANDLERS_H


#include <windows.h>
#include <string>
#include <vector>
#include "task.h"
#include "command.h"
#include "profile.h"

void HandleCmd(Task* t, MalleableProfile* currentProfile);

#endif
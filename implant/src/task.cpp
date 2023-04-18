#include <windows.h>
#include <string>
#include "debug.h"
#include "task.h"
#include "httpclient.h"
#include "constants.h"
#include "rapidjson/document.h"
#include "obfuscator/MetaString.h"

using namespace std;
using namespace andrivet::ADVobfuscator;


Task *parseTask(std::string taskJson)
{
    Task *task = new Task;
    rapidjson::Document document;
    document.Parse(taskJson.c_str());
    // Check if document is {}

    if (document.HasParseError())
    {
        DEBUG_PRINTF("Error parsing task json\n");
        return NULL;
    }

    if (document.ObjectEmpty())
    {
        // DEBUG_PRINTF("Task json is empty\n");
        return NULL;
    }

    if (!document.HasMember(OBFUSCATED("task_id")))
    {
        DEBUG_PRINTF("Task json does not contain task_id\n");
        return NULL;
    }

    if (!document.HasMember(OBFUSCATED("opcode")))
    {
        DEBUG_PRINTF("Task json does not contain opcode\n");
        return NULL;
    }

    if (!document.HasMember(OBFUSCATED("args")))
    {
        DEBUG_PRINTF("Task json does not contain args\n");
        return NULL;
    }

    task->taskId = document[OBFUSCATED("task_id")].GetString();
    task->opcode = document[OBFUSCATED("opcode")].GetInt();
    task->args = &document[OBFUSCATED("args")];

    //cout << task << endl;

    return task;
}

std::ostream &operator<<(std::ostream &os, const Task &task)
{
    os << "Task: " << task.taskId << ", Opcode: " << task.opcode << ", args: " << task.args;
    return os;
}
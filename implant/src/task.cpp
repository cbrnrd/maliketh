#include <windows.h>
#include <string>
#include "debug.h"
#include "task.h"
#include "httpclient.h"
#include "constants.h"
#include "rapidjson/document.h"
#include "obfuscator/MetaString.h"

using namespace std;

Task* parseTask(std::string taskJson) {
    Task* task = new Task;
    rapidjson::Document document;
    document.Parse(taskJson.c_str());

    if (document.HasParseError()) {
        DEBUG_PRINTF("Error parsing task json\n");
        return NULL;
    }

    if (!document.HasMember(OBFUSCATED("taskId"))) {
        DEBUG_PRINTF("Task json does not contain taskId\n");
        return NULL;
    }

    if (!document.HasMember(OBFUSCATED("opcode"))) {
        DEBUG_PRINTF("Task json does not contain opcode\n");
        return NULL;
    }

    if (!document.HasMember(OBFUSCATED("args"))) {
        DEBUG_PRINTF("Task json does not contain args\n");
        return NULL;
    }

    task->taskId = document[OBFUSCATED("taskId")].GetString();
    task->opcode = document[OBFUSCATED("opcode")].GetInt();
    task->args = &document[OBFUSCATED("args")];

    return task;
}
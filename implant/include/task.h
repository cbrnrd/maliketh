#ifndef TASK_H_
#define TASK_H_

#include <windows.h>
#include <string>

#include "./rapidjson/document.h"
#include "./obfuscator/MetaString.h"

#define TASK_STATUS_CREATED OBFUSCATED("CREATED")
#define TASK_STATUS_TASKED OBFUSCATED("TASKED")
#define TASK_STATUS_COMPLETED OBFUSCATED("COMPLETED")
#define TASK_STATUS_ERROR OBFUSCATED("ERROR")

typedef struct {
    std::string taskId;
    INT opcode;
    rapidjson::Value* args; // Could be any json object
} Task;

std::ostream &operator<<(std::ostream &os, const Task &task);

/**
 * Function used to convert a rapidjson Document to a Task struct
 */
Task* parseTask(std::string taskJson);

#endif

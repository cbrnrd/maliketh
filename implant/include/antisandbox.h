#ifndef ANTISANDBOX_H_
#define ANTISANDBOX_H_

#include <windows.h>
#include "obfuscator/MetaString.h"


BOOL IsSandboxed();
void MemeIfSandboxed();
void StartMemeSandboxThread();

#endif // ANTISANDBOX_H_
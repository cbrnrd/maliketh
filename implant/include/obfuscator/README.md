# Compile-time string obfuscator

From: <https://github.com/andrivet/ADVobfuscator>

## Usage

```c++
#include "obfuscator/MetaString.h"
using namespace andrivet::ADVobfuscator;


// Define and use a string
cout << OBFUSCATED("Hello world") << endl;

// Define a string
auto obfuscated = DEF_OBFUSCATED("Hello world");
cout << obfuscated.decrypt() << endl;
```

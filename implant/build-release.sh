#!/bin/bash

# Check if the RELEASE_OUTFILE variable is set
if [ -z "$RELEASE_OUTFILE" ]; then
    RELEASE_OUTFILE=implant_release.exe
fi

# Make a copy of the src directory
cp -r src/ src_copy/

# Find all .cpp files in the copied directory and replace the text
find src_copy/ -name "*.cpp" -exec sed -i 's/DEBUG_PRINTF(/\/\/ DEBUG_PRINTF(/g' {} \;

# Build the implant
BUILDER_OPTS=$BUILDER_OPTS RELEASE_OUTFILE=$RELEASE_OUTFILE make release

# Remove the copied directory
rm -rf src_copy/

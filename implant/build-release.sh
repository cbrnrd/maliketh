#!/bin/bash


# Make a copy of the src directory
cp -r src/ src_copy/

# Find all .cpp files in the copied directory and replace the text
find src_copy/ -name "*.cpp" -exec sed -i 's/DEBUG_PRINTF(/\/\/ DEBUG_PRINTF(/g' {} \;


# Build the implant
make release

# Remove the copied directory
rm -rf src_copy/
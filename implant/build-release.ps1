# Make a copy of the src directory
Copy-Item -Recurse -Force "src\" "src_copy\"

# Find all .cpp files in the copied directory and replace the text
Get-ChildItem -Path "src_copy\" -Filter "*.cpp" -Recurse | ForEach-Object { (Get-Content $_.FullName) | ForEach-Object { $_ -replace 'DEBUG_PRINTF\(', '// DEBUG_PRINTF(' } | Set-Content $_.FullName }

# Build the implant
& mingw32-make.exe release

# Remove the copied directory
Remove-Item -Recurse -Force "src_copy\"

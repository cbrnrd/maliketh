# Opcodes

This document outlines the different Opcodes that are used for the execution of jobs.
Arguments for each opcode must be a valid JSON type supported by [SQLAlchemy](https://docs.sqlalchemy.org/en/14/core/type_basics.html#sqlalchemy.types.JSON).

| Opcode | Action | Description |
|--------|--------|-------------|
| `0x01` | `CMD` | Execute a command on the implant |
| `0x02` | `SELFDESTRUCT` | Delete and kill the implant |
| `0x03` | `SYSINFO` | Get system information |
| `0x04` | `SLEEP` | Sleep for a specified amount of time |
| `0x05` | `UPDATE_CONFIG` | Update the malleable configuration of the implant |
| `0x06` | `DOWNLOAD` | Download a file from the implant |
| `0x07` | `UPLOAD` | Upload a file to the implant |
| `0x08` | `INJECT` | Inject shellcode from the server into a process |
| `0x09` | `CHDIR` | Change the working directory of the implant |
| `0x0A` | `PWD` | Get the current working directory of the implant |
| `0x0B` | `GETENV` | Gets all environment variables from the implant |
| `0x0C` | `LS` | List files in the current working directory |
| `0x0D` | `PS` | List running processes |
| `0x0E` | `WHOAMI` | Get the current user of the implant |
| `0x0F` | `DISABLE_DEFENDER` | Disable Windows Defender |

## CMD

Args: List of strings (command and arguments)
Ex: `["ipconfig", "/all"]`

Expected response: The output of the command, if any.

## SELFDESTRUCT

Args: None

Expected response: None

## SYSINFO

Args: None

Expected response: A map of strings (key, value) of system information.

## SLEEP

Args: 1 integer (seconds to sleep)

Expected response: None

## UPDATE_CONFIG

Args: A map of strings (key, value) to update the malleable configuration.
Invalid keys will be ignored. See [profile.md](profile.md#client-options) for a list of valid keys.
Ex: `{"kill_date": "2021-01-01"}`

Expected response: None

## DOWNLOAD

Args: List of 1 string (path to file on the implant)
Ex: `"C:\\Users\\user\\Desktop\\file.txt"`

Expected response: base64 encoded file content

## UPLOAD

Args: List of 2 strings (path to save the file on the implant)
Ex: `["C:\\Users\\user\\Desktop\\file.txt", "b64encoded-file-content=="]`

Expected response: None

## INJECT

Args: List of 2 strings (base64 encoded shellcode, process name/id)
Ex: `["shellcode==", "notepad.exe"]`

Expected response: None

## CHDIR

Args: 1 string (path to change to)
Ex: `"C:\\Users\\user\\Desktop"`

Expected response: None

## PWD

Args: None

Expected response: 1 string (current working directory)

## GETENV

Args: None

Expected response: A map of strings (key, value) of environment variables.

## LS

Args: None

Expected response: A list of strings (file names)

## PS

Args: None

Expected response: A map of strings (pid, name) of running processes.

## WHOAMI

Args: None

Expected response: 1 string (current user)

## DISABLE_DEFENDER

Args: None

Expected response: None

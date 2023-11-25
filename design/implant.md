# Implant

An implant is a single instance of a payload that is running on a target system. The implant establishes a session with the C2 and periodically checks in with the C2 to receive commands. The implant is responsible for executing the commands it receives from the C2.

## Implant server info

The C2 server keeps some additional information about each implant. This information is stored in the `implant` table in the database. The following fields are stored:

| Field | Meaning | Notes |
|:-----|:--------|:------|
| `id` | The internal ID | DB primary key|
| `implant_id` | The implant ID | This is the main identifier that should be used throughout the codebase|
| `hostname` | The hostname of the target system | |
| `ip` | The WAN IP address of the target system | |
| `os` | The operating system of the target system | Windows version |
| `arch` | The architecture of the target system | x86, x64, etc. |
| `user` | The user that the implant is running as | |
| `server_sk` | The NaCl key to use to decrypt messages from the implant | Base64 encoded |
| `implant_pk` | The NaCL key to use to encrypt messages to the implant | Base64 encoded |
| `created_at` | The time the implant was created | |
| `last_seen` | The time the implant was last seen | |


## Implant configuration

Each implant has a configuration stored in the database. This is intentionally separate from the implant server info because the configuration will be sent to the implant and will be used to configure the implant. The configuration is stored in the `implant_config` table in the database. The following fields are stored:

| Field | Meaning | Notes |
|:-----|:--------|:------|
| `id` | The internal ID | DB primary key|
|`implant_id`| The implant ID this config belongs to | this is not the `id` field of Implant, rather the `implant_id` field |
| `cookie` | The cookie name used for implant identification | |
| `kill_date` | The timestamp of the kill date | |
| `user_agent` | The user agent to use for HTTP requests | |
| `auto_self_destruct` | Self destruct on failed checkins | |
| `sleep_time` | The number of seconds to sleep between tasks and checkins | |
| `jitter` | The % amount of jitter to add to the sleep time | |
| `max_retries` | The maximum number of retries to attempt before giving up | |
| `retry_wait` | The number of seconds to wait between retries | |
| `retry_jitter` | The % amount of jitter to add to the retry wait time | |
| `enc_key` | The base64 encoded public key of the server to use for encryption | |
| `tailoring_hash_function` | The hash function to use for the tailoring hash | Unused |
| `tailoring_hash_rounds` | The number of rounds to use for the tailoring hash | Unused |
| `tailoring_hashes` | A list of hashes to use for tailoring | Unused |

## Registration and authentication

Implants register with the C2 after an initial sleep time. There is a hardcoded AES key for the initial registration. The implant sends a JSON payload containing the following information:

| Field | Meaning | Notes |
|:-----|:--------|:------|
|`txid`| The base64 encoded, encrypted,  NaCl public key of the implant | `Base64(AES(implant_public_key))` |

The C2 decrypts the payload and stores the implant's public key in the database. The C2 then generates a NaCl keypair and sends the public key and the configured initial configuration to the implant. All further communication is encrypted with the NaCl keypair, and then base64 encoded.

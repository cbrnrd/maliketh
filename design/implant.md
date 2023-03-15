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
| `aes_key` | The AES key used to encrypt the implant's communications with the C2 | Base64 encoded |
| `aes_iv` | The AES IV used to encrypt the implant's communications with the C2 | Base64 encoded |
| `created_at` | The time the implant was created | |
| `last_seen` | The time the implant was last seen | |
| `kill_date` | The time the implant should be killed | |
| `jitter` | The percentage value (as a float <= 1.0) of jitter to add to the implant's check-in interval | For example, if the implant checks in every 60 seconds and the jiter is `0.2`, $60<= sleep <= 60 + (60 * 0.2)$|

## Registration and authentication
Mostly TODO, but here's the general idea.

There are a few approaches that could work for authentication:

1. Use PKI and ECC keys to encrypt the C2 channel. The client and server would exchange public keys on the first request, then encrypt everything with those keys.
2. Use some sort of Password Authenticated Key Exchange (PAKE) to generate a shared secret key. This would be similar to the above, but would use a password instead of a key pair. (OPAQUE, SPAKE2, etc.)
3. Use Authenticated Encryption (ChaCha20Poly1305, AES-GCM, etc)


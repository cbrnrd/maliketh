# Implant-C2 HTTP Spec

These are the *default* endpoints for the implant HTTP server. These endpoints can and should be changed by modifying the `server/config/profiles/default.yaml` file.

| Endpoint | Verb | Purpose | Details |
|:-------- | :-- | :------ | :-----: |
| `/c2/register` | `POST` | Registers this implant | |
| `/c2/checkin` | `GET` | Checks in with the C2 to see if there are any tasks | |
| `/c2/task` | `POST` | Send the result of a task | |

## Crypto

*(Note, this does not apply for registration)*

All requests and responses after the initial registration are encrypted via libsodium `Box`es. On the server side, each implant has a unique keypair. The server-side private key is stored in the the database, and the public key is sent to the implant during registration. All requests and responses with bodies are encrypted then Base64 encoded. *Assume for all requests and response, the data has already been decoded and derypted.*

## Examples

**Authentication**
Implants can identify itself by setting the `SESSID` cookie to the implant's ID. This value can be customized in the configuration file with the `server.implant_id_cookie` key.

### `/c2/register`

This route allows implants to register with the server. The implant should send the following JSON object:

```json
{
  "txid": "encrypted_base64_encoded_public_key"
}
```

| Field | Purpose |
|:----- | :------ |
| `txid` | The Base64 encoded LibSodium public key to use for encryption from server -> implant, encrypted with the C2 registration password |

Example responses:

* 200:

  On successful registration:

  ```json
  {
    "status": true,
    "k": "base64_encoded_server_public_key",
    "c": "base64_encoded_encrypted_maleable_config"
  }
  ```

  | Field | Purpose |
  |:----- | :------ |
  | `status` | `true` if registration was successful, `false` otherwise |
  | `k` | The Base64 encoded LibSodium public key to use for encryption from implant -> server |
  | `c` | The Base64 encoded maleable profile to use. This is encrypted with the key in `k`. See below and [profile.md](../profile.md) |

  `c` decrypted would look something like:

  ```json
  {
  "status": true,
  "id": "74a45dc8dbec3008e74f91da3d2d05fa",
  "config": {
      "cookie": "SESSID",
      "kill_date": "",
      "sleep_time": 60,
      "jitter": 0.1,
      ...
    }
  }
  ```

  | Field | Purpose |
  |:----- | :------ |
  | `status` | `true` if registration was successful, `false` otherwise |
  | `id` | The implant ID |
  | `config` | The maleable profile to use. See [profile.md](../profile.md) |

* 401:
This can happen for a few reasons:

* Empty request body or incorrect MIME type
* Missing or invalid `txid` field
* Improperly formatted public key
* Public key already exists in the database

### `/c2/checkin`

This endpoint is used to check in with the C2 server and pull down the next task, if there is one. The response will be encrypted with the key received during registration.

Example responses:

* 200:

  If there are no tasks:

  ```json
  {}
  ```

  If there is a task:

  ```json
  {
  "args": ["ls", "-la"],
  "opcode": 0,
  "task_id": "8af2838e2d70e2222aeb66459cec096e"
  }
  ```

  See [job.md](../job.md) for more details on the job object.
  Results of a job should be sent back to the appropriate endpoint to mark the task as completed.

  For the meanining of `opcode`s, see [opcodes.md](../opcodes.md).

## `/c2/task`

This endpoint is used to send the results of a task back to the C2 server. The request body should be encrypted with the key received during registration, then base64 encoded.

Example request:

```json
  {
    "status": true,
    "tid": "task_id_here",
    "output": "b64_output_here",
  }
```

If there is no output, set `output` to an empty string.
If you'd like the output to be displayed in a table, `output` should be a JSON object (still base64 encoded).

| Field | Purpose |
|:----- | :------ |
| `status` | `true` if the task was successful, `false` otherwise |
| `tid` | The task ID |
| `output` | The Base64 encoded output of the task |

Example responses:

* 200:

```ascii
OK
```

* 401:

This can happen for a few reasons:

* Improperly formatted request body
* Missing implant identification cookie
* Implant with given ID doesn't exist

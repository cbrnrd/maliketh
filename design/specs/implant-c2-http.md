# Implant-C2 HTTP Spec

These are the *default* endpoints for the implant HTTP server. These endpoints can and should be changed by modifying the `server/config/admin/routes.yaml` file.

| Endpoint | Verb | Purpose | Details |
|:-------- | :-- | :------ | :-----: |
| `/c2/register` | `POST` | Registers this implant | |
| `/c2/checkin` | `GET` | Checks in with the C2 to see if there are any jobs | |
| `/c2/task/<tid>` | `POST` | Send the result of a task with id `tid` | |

## Crypto

All request bodies should be AES-GCM encrypted and base64 encoded. The encryption key and IV/AAD is received from the C2 during registration. See `/c2/register` for more details.

This is the encryption and decrpytion functions in python:

```python
from cryptography.hazmat.primitives.ciphers.aead import AESGCM

def encrypt(data: bytes, aad: bytes, key: bytes) -> bytes:
    """
    Encrypt a given byte array using AES-GCM
    """
    aesgcm = AESGCM(key)
    nonce = os.urandom(12)
    return aesgcm.encrypt(nonce, data, aad) + nonce

def decrypt(data: bytes, aad: bytes, key: bytes) -> bytes:
    """
    Decrypt a given byte array using AES-GCM
    """
    aesgcm = AESGCM(key)
    nonce = data[-12:]
    data = data[:-12]
    return aesgcm.decrypt(nonce, data, aad)
```

Note the nonce is appended to the end of the encrypted data. The nonce should be 12 bytes long, and should be randomly generated. The nonce should be unique for each encryption. [Relevant docs](https://cryptography.io/en/latest/hazmat/primitives/aead/#cryptography.hazmat.primitives.ciphers.aead.AESGCM)

## Examples

**Authentication**
Implants can identify itself by setting the `SESSID` cookie to the implant's ID. This value can be customized in `routes.yaml`.

### `/c2/register`

This route allows implants to register with the server. The implant should send the following JSON object:

```json
{
  "u": "username",
  "t": "random base64e string", 
}
```

| Field | Purpose |
|:----- | :------ |
| `u` | The username of the account the implant is running on |
| `t` | A base64 encoded random string to be used as the AAD for AES-GCM |

*`t` should be the AAD for all following encrypted requests.*

Example responses:

* 200:

  On successful registration:

  ```json
  {
  "id": "b58f9c46dc2d87033df6281e07c89cce",
  "iv": "mDnpmVCqopiRUvvz",
  "key": "Qxe9RZg6pN3xgKkD/lCV6a4968jzd5f5YliPrXDp0dY=",
  "status": true
  }
  ```

  | Field | Purpose |
  |:----- | :------ |
  | `id` | The implant's ID |
  | `aad` | The IV/AAD to use for AES-GCM (base64 encoded) |
  | `key` | The key to use for AES-GCM (base64 encoded) |
  | `status` | Whether or not the registration was successful |

* 401:
Occurs if `u` or `t` is invalid or missing, or the body is not JSON.

`u` has a length limit of 128 characters.
`t` is invalid if it is not base64 encoded.

### `/c2/checkin`

This endpoint is used to check in with the C2 server and pull down the next task, if there is one. The response will be encrypted with the key and AAD received during registration.

Example responses:

* 200:

  If there are no tasks:
    
  ```json
  {}
  ```

  If there is a task:

  ```json
  {
  "args": "{ls,-la}",
  "created_at": "2023-03-30 17:11:59.246815",
  "executed_at": null,
  "implant_id": "00eec44c7ce141406f0eb218d173b11f",
  "opcode": 0,
  "operator_name": "admin",
  "output": null,
  "read_at": "2023-03-30 17:12:17",
  "status": "TASKED",
  "task_id": "8af2838e2d70e2222aeb66459cec096e"
  }
  ```

  See [job.md](../job.md) for more details on the job object.
  Results of a job should be sent back to the appropriate endpoint to mark the task as completed.

  For the meanining of `opcode`s, see [opcodes.md](../opcodes.md).

* 404:
This is returned if the implant ID in the cookie is missing or not in the database.

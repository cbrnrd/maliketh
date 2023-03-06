# Operator-C2 HTTP Spec

| Endpoint | Verb | Purpose | Details |
|:-------- | :-- | :------ | :-----: |
| `/op/tasks/list` | `GET` | Lists currently running or new tasks given by this operator | |
| `/op/tasks/add`  | `POST`| Adds a new task for an implant| |
| `/op/tasks/results/:id` | `GET` | Gets the results of a task | |
| `/op/tasks/delete/:id`| `DELETE` | Deletes the task with the given ID | |
| `/op/implant/list` | `GET` | Lists all implants | |
| `/op/implant/info/:id` | `GET` | Gets basic information about the implant with the given ID | |
| `/op/implant/kill/:id` | `GET` | Removes the implant with the given ID from the C2 database | |
| `/op/auth/token/request` | `GET` | Used for fetching an operators authentication token | [example](#opauthtokenrequest) |
| `/op/auth/token/revoke` | `DELETE` | Revokes the current operator authentication token | [example](#opauthtokenrevoke) |

## Examples

### `/op/auth/token/request`

Request headers:
| Name | Content |
| :-- | :----- |
| `X-ID` | `operator_name` |
| `X-Signature` | `Base64(enc(sign(logon_secret, operator_signing_key), server_public_key))` |

Response fields:
| Name | Meaning |
| :-- | :----- |
| `status` | Whether or not authentication succeeded |
| `token` | The authentication token to use for all following requests |
| `rmq_queue` | The RabbitMQ queue to subscribe to |
| `rmq_host` | The hostname or IP address where the RabbitMQ server is located |
| `rmq_port` | The TCP port of the RabbitMQ instance |

__Example request__:

```http
GET /op/auth/token/request
X-ID: admin
X-Signature: 85wV03Vh+BBb2LjTB9rkf1+0Cg==
```

__Response (JSON)__:

The response will be encrypted with the operators public key. The following is an example decrypted response:

```json
{
  "status": true,
  "token": "asdf1234qwertyuiop",
  "rmq_queue": "queue_name",
  "rmq_host": "rabbit.local",
  "rmq_port": 1338
}
```

If the operator request is invalid, the following will be returned (with status code 400):

```json
{
  "status": false,
  "message": "Unknown operator key"
}
```

### `/op/auth/token/revoke`

By default, authencation tokens have an expiration of 6 hours. This endpoint should be used to prematurely revoke a currently active authentication token.

__Example request__:

```http
DELETE /op/auth/token/revoke
Authentication: Bearer <insert_token_here>
```

__Example response__:
Success:

```json
{
  "status": true
}
```

Failure:

```json
{
  "status": false,
  "message": "Invalid token"
}
```

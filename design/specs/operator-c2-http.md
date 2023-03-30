# Operator-C2 HTTP Spec

These are the *default* endpoints for the operator HTTP server. These endpoints can be changed by modifying the `server/config/admin/routes.yaml` file.

| Endpoint | Verb | Purpose | Details |
|:-------- | :-- | :------ | :-----: |
| `/op/stats` | `GET` | Gets basic statistics about the C2 | |
| `/op/tasks/list` | `GET` | Lists currently running or new tasks given by this operator | |
| `/op/tasks/add`  | `POST`| Adds a new task for an implant| |
| `/op/tasks/results/:id` | `GET` | Gets the results of a task | |
| `/op/tasks/delete/:id`| `DELETE` | Deletes the task with the given ID | |
| `/op/implant/list` | `GET` | Lists all implants | |
| `/op/implant/info/:id` | `GET` | Gets basic information about the implant with the given ID | |
| `/op/implant/kill/:id` | `GET` | Removes the implant with the given ID from the C2 database | |
| `/op/auth/token/request` | `GET` | Used for fetching an operators authentication token | [example](#opauthtokenrequest) |
| `/op/auth/token/revoke` | `DELETE` | Revokes the current operator authentication token | [example](#opauthtokenrevoke) |
| `/op/auth/token/status` | `GET` | Checks the status of the current operator authentication token | [example](#opauthtokenstatus) |

## Examples

For all requests except `/op/auth/token/request`, the authentication token should be present in the `Authorization` header as a bearer token.

### `/op/tasks/add`

Example request body:

```json
{
 "opcode": 0,
 "implant_id": "a5a90992dd62b62f78d5541c8a07c3b4",
 "args": [
  "ls",
  "-la"
 ]
}
```

Example responses:

* 200:

```json
{
 "status": true,
 "task": {
  "args": "{ls,-la}",
  "created_at": "2023-03-23 14:59:26.078969",
  "executed_at": null,
  "implant_id": "a5a90992dd62b62f78d5541c8a07c3b4",
  "opcode": 0,
  "operator_name": "admin",
  "output": null,
  "read_at": null,
  "status": "CREATED",
  "task_id": "9cf5780f66e1184600206d7c1327e1fb"
 }
}
```

* 400:
For example if the `opcode` field is missing:

```json
{
 "msg": "Invalid task, missing fields: opcode",
 "status": false
}
```

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

### `/op/auth/token/status`

This endpoint can be used to check the status of the current authentication token.

__Example request__:

```http
GET /op/auth/token/status
Authentication: Bearer <insert_token_here>
```

__Example response__:
Success:

```json
{
  "status": true,
  "msg": "Authenticated"
}
```

Failure:

```json
{
  "status": false,
  "message": "Invalid token"
}
```

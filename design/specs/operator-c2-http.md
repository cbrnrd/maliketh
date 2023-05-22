# Operator-C2 HTTP Spec

These are the *default* endpoints for the operator HTTP server. These endpoints can be changed by modifying the `server/config/admin/routes.yaml` file.

| Endpoint | Verb | Purpose | Details |
|:-------- | :-- | :------ | :-----: |
| `/op/stats` | `GET` | Gets basic statistics about the C2 | [example](#opstats) |
| `/op/tasks/list` | `GET` | Lists currently running or new tasks given by this operator | [example](#optaskslist) |
| `/op/tasks/add`  | `POST`| Adds a new task for an implant| [example](#optasksadd) |
| `/op/tasks/results/:task_id` | `GET` | Gets the results of a task | [example](#optasksresultstask_id) |
| `/op/tasks/delete/:task_id`| `DELETE` | Deletes the task with the given ID | [example](#optasksdeletetask_id) |
| `/op/implant/config/:implant_id` | `GET` | Gets the malleable configuration of the implant with the given ID | [example](#get-opimplantconfigimplant_id) |
| `/op/implant/config/:implant_id` | `POST` | Updates the malleable configuration of the implant with the given ID | [example](#post-opimplantconfigimplant_id) |
| `/op/implant/list` | `GET` | Lists all implants | [example](#opimplantlist) |
| `/op/implant/kill/:id` | `GET` | Removes the given implant from the database and purges it from the affected system. | [example](#opimplantkillimplant_id) |
| `/op/implant/build` | `POST` | Builds an implant with the given configuration | [example](#opimplantbuild) |
| `/op/auth/token/request` | `GET` | Used for fetching an operators authentication token | [example](#opauthtokenrequest) |
| `/op/auth/token/revoke` | `DELETE` | Revokes the current operator authentication token | [example](#opauthtokenrevoke) |
| `/op/auth/token/status` | `GET` | Checks the status of the current operator authentication token | [example](#opauthtokenstatus) |

## Examples

For all requests except `/op/auth/token/request`, the authentication token should be present in the `Authorization` header as a bearer token.

All requests can return a 401 status code if the token is invalid or expired. This will have the body:

```json
{
 "msg": "Not authenticated",
 "status": false
}
```

### `/op/stats`

This endpoint is used to get basic statistics about the C2.

__Example response__:
  
  ```json
  {
    "active_tasks": 0,
    "implants": 1,
    "operators": 1,
    "status": true,
    "total_tasks": 1,
    "uptime": "3:38:16.451002"
  }
  ```

### `/op/tasks/list`

This endpoint is used to list all active (aka not `COMPLETE`) tasks that have been created by the operator.

__Example response__:

```json
{
 "status": true,
 "tasks": [
  {
   "args": {
    "kill_date": "2024-12-12"
   },
   "created_at": "2023-04-03 21:44:52.883699",
   "executed_at": null,
   "implant_id": "abf0b02e502cae7e6c7b8c368822fe9b",
   "opcode": 5,
   "operator_name": "admin",
   "output": null,
   "read_at": "2023-04-03 21:44:54",
   "status": "TASKED",
   "task_id": "64ce6a88b75a377a8320c87b9b6d4f3b"
  },
  {
   "args": [
    "ls",
    "-la"
   ],
   "created_at": "2023-04-04 01:26:48.660957",
   "executed_at": null,
   "implant_id": "abf0b02e502cae7e6c7b8c368822fe9b",
   "opcode": 0,
   "operator_name": "admin",
   "output": null,
   "read_at": null,
   "status": "CREATED",
   "task_id": "b0cb6baaad0b1ed69362cffb4384e6d6"
  }
 ]
}
```

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

### `/op/tasks/results/:task_id`

This endpoint is used to get the output of a task. `result` is the base64 encoded output of the task.
Depending on the type/opcode of the task, this could be anything from a string to a binary file.

__Example response__:

```json

{
 "result": "dG90YWwgNDgKZHJ3eHIteHIteCAxOCBjYXJ0ZXIgc3RhZmYgIDU3NiBBcHIgIDEgMjI6NTMgLgpkcnd4ci14ci14IDEzIGNhcnRlciBzdGFmZiAgNDE2IE1hciAzMCAxMzoyNCAuLgotcnctci0tci0tICAxIGNhcnRlciBzdGFmZiAgMTY0IE1hciAyMyAxMDoyMiBEb2NrZXJmaWxlLmMyCi1ydy1yLS1yLS0gIDEgY2FydGVyIHN0YWZmICAxNTkgTWFyIDIzIDEwOjIyIERvY2tlcmZpbGUub3BlcmF0b3IKLXJ3LXItLXItLSAgMSBjYXJ0ZXIgc3RhZmYgMTk3MiBNYXIgMjQgMTI6NTMgTXlDZXJ0aWZpY2F0ZS5jcnQKLXJ3LS0tLS0tLSAgMSBjYXJ0ZXIgc3RhZmYgMzI3MiBNYXIgMjQgMTI6NTIgTXlLZXkua2V5Ci1ydy1yLS1yLS0gIDEgY2FydGVyIHN0YWZmICAgIDAgRmViIDIwIDE4OjA0IFJFQURNRS5tZApkcnd4ci14ci14ICA0IGNhcnRlciBzdGFmZiAgMTI4IEFwciAgMiAxNzo1MyBfX3B5Y2FjaGVfXwotcnctci0tci0tICAxIGNhcnRlciBzdGFmZiAzMzQ3IEFwciAgMiAxNzo1MCBhcHAucHkKZHJ3eHIteHIteCAgNCBjYXJ0ZXIgc3RhZmYgIDEyOCBBcHIgIDEgMjI6NTMgY29uZmlnCi1ydy1yLS1yLS0gIDEgY2FydGVyIHN0YWZmIDQyNTUgQXByICAzIDEwOjQ3IGNyZWF0ZV9vcGVyYXRvci5weQotcnctci0tci0tICAxIGNhcnRlciBzdGFmZiAxNDc2IEFwciAgMSAyMjo1MyBkb2NrZXItY29tcG9zZS55bWwKZHJ3eHIteHIteCAgMyBjYXJ0ZXIgc3RhZmYgICA5NiBGZWIgMjIgMTA6MjggaW5zdGFuY2UKLXJ3LXItLXItLSAgMSBjYXJ0ZXIgc3RhZmYgIDExNSBNYXIgMjMgMTA6MjIgbWFrZV9kYi5weQpkcnd4ci14ci14IDE0IGNhcnRlciBzdGFmZiAgNDQ4IEFwciAgMiAxNToyMiBtYWxpa2V0aAotcnctci0tci0tICAxIGNhcnRlciBzdGFmZiAgMzQ3IE1hciAgNSAxODoxNyByZXF1aXJlbWVudHMudHh0Ci1ydy1yLS1yLS0gIDEgY2FydGVyIHN0YWZmIDE4MjEgQXByICAyIDE3OjU2IHRlc3RfY2xpZW50LnB5Ci1ydy1yLS1yLS0gIDEgY2FydGVyIHN0YWZmIDIxOTcgQXByICAzIDE3OjQzIHRlc3RfY3J5cHRvLnB5Cg==",
 "status": true
}
```

### `/op/tasks/delete/:task_id`

This endpoint deletes a task from the database. This is useful if you want to delete a task that has been executed or was accidentally created.

__Example response__:

```json
{
 "status": true
}
```

If the task does not exist, the response will be:

```json
{
 "msg": "Unknown task",
 "status": false
}
```

### `GET /op/implant/config/:implant_id`

This endpoint is used to get the malleable configuration of an implant. See the [profile.md](../profile.md#client-options) document for more information on malleable configuration and the fields in it.

__Example response__:

```json
{
	"config": {
		"cookie": "SESSID",
		"enc_key": "22WQpoz4PCyvifzO4GunTjP52fx4kkZElGtqFg8kuwM=",
		"jitter": 0.0,
		"kill_date": "",
		"max_retries": -1,
		"sleep_time": 0,
		"tailoring_hash_function": "sha256",
		"tailoring_hash_rounds": 1,
		"tailoring_hashes": []
	},
	"status": true
}
```

If the implant is not found:

```json
{
	"msg": "Unknown implant",
	"status": false
}
```

### `POST /op/implant/config/:implant_id`

This endpoint is used to set the malleable configuration of an implant. See the [profile.md](../profile.md#client-options) document for more information on malleable configuration and the fields in it. The request should be a valid JSON object with *only* the fields you want to change. The fields that are not specified will not be changed. Any invalid fields will be ignored.

Behind the scenes, this endpoint submits a job with the `UPDATE_CONFIG` opcode

__Example request__:
If you want to change the `kill_date` and `sleep_time` fields.

```json
{
	"kill_date": "2024-12-12",
	"sleep_time": 30
}
```

__Example response__:

The response is the same as the response for getting a task.

```json
{
	"status": true,
	"task": {
		"args": {
			"asdf": 123,
			"kill_date": "2024-12-12",
			"sleep_time": 30
		},
		"created_at": "2023-04-04 02:32:10.22647",
		"executed_at": null,
		"implant_id": "abf0b02e502cae7e6c7b8c368822fe9b",
		"opcode": 5,
		"operator_name": "admin",
		"output": null,
		"read_at": null,
		"status": "CREATED",
		"task_id": "8869c6bdf9dd702026baa16b808f0272"
	}
}
```

If the request is invalid or empty:

```json
{
  "msg": "No valid fields found in request",
  "status": false
}
```

### `/op/implant/list`

This endpoint is used to get a list of all implants in the database.

__Example response__:

```json
{
	"implants": [
		{
			"arch": null,
			"created_at": "2023-04-03 21:44:29",
			"hostname": "localhost:8080",
			"implant_id": "abf0b02e502cae7e6c7b8c368822fe9b",
			"implant_pk": "f+W9O+P34mPlijXC4GiPiJupedQAraIsZ1CTFl6MOgs=",
			"ip": "172.19.0.1",
			"last_seen": "2023-04-03 21:44:29",
			"os": null,
			"server_sk": "EuzRc4VlYg/Al0CaFqYjL0iRfh92758gbtVxDeb6yqg=",
			"user": ""
		},
		...
	],
	"status": true
}
```

### `/op/implant/kill/:implant_id`

This endpoint is used to kill an implant. This is useful if you want to kill an implant that is no longer needed. NOTE: This will kill the implant immediately, and it will not be able to reconnect to the server.

__Example response__:

```json
{
  "status": true
}
```

If the implant is not found:
  
  ```json
  {
    "msg": "Unknown implant",
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

### `/op/implant/build`

This endpoint is used to build a new implant. Note that depending on the power of the C2 server, this may take a while (a few minutes). The request should be a valid JSON object with any of the following fields:

| Name | Meaning | Default |
| :-- | :----- | :----- |
| `initial_sleep_seconds` | The number of seconds to wait before connecting to the server | `180` |
| `schtask_persist` | Whether or not to use schtasks for persistence | `true` |
| `use_antidebug` | Whether or not to use antidebugging techniques | `true` |
| `kill_parent` | Whether or not to kill the parent process after spawning (unused) | `true` |
| `use_antivm` | Whether or not to use antivm techniques | `true` |
| `scheduled_task_name` | The name of the scheduled task | `MicrosoftEdgeUpdateTaskMachineUA` |
| `register_max_retries` | The maximum number of times to retry registering with the server | `5` |

__Example request__:

```json
{
  "initial_sleep_seconds": 180,
  "schtask_persist": true,
  "use_antidebug": true,
  "kill_parent": true,
  "use_antivm": true,
  "scheduled_task_name": "MicrosoftEdgeUpdateTaskMachineUA",
  "register_max_retries": 5
}
```

__Example response__:

```json
{
  "implant": "base64_encoded_implant_pe"
}
```

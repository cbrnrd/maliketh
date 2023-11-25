# Operators

Operators are a single user of the C2 framework. A server can have many operators.

## Creating an operator

In order to create an operator, the admin of the server needs to run the `create_operator.py` script in the `operator` docker container:

For this example, assume the docker container has ID `37fc3915b843`.

```bash
docker exec 37fc3915b843 python create_operator.py --name operator_name
```

More options can be found in the help section of that script.

This will create an operator configuration file that will have a few fields.

```json
{
  "name": "operator_name",
  "c2": "c2.ip.addr.here",
  "c2_port": 1337,
  "secret": "...",
  "public": "...",
  "signing_key": "...",
  "verify_key": "...",
  "server_pub": "...",
  "login_secret": "..."
}
```

Field definitions:
| Field | Meaning  | Notes |
|:-----|:-------- | ------ |
| `name`   | The name of the operator this config file is for | |
| `secret` | The Base64 encoded ED25519 secret key for this operator | |  
| `c2` | The IP address or domain name of the C2 server this operator is registered for | |
| `c2_port` | The TCP port to connect to the C2 server on | |
| `public` | The Base64 encoded ED25519 public key for this operator |  |
| `signing_key` | The Base64 encoded ED25519 secret key for signing messages to the C2 server | |
| `verify_key`  | The Base64 encoded ED25519 public key for verifying messages signed by `signing_key` | |
| `server_pub` | The Base64 encoded ED25519 public key for the C2 server | Used for encrypting operator -> C2 traffic|
| `login_secret` | A secret string to be encrypted and signed by `secret` on authentication | |

## Authentication

In order to authenticate, the operator must send a `GET` request to `/op/auth/token/request` with the the `X-ID` header set to the `logon_secret` encrypted and signed by `secret` (PEM encoded then base64urlsafe encoded). The server will then send back an authentication token (and other info) that the operator must use for all following requests. This response will be encrypted with `public` and must be decrypted by the operator.

More info about request and response bodies for this can be found in [specs/operator-c2-http.md](./specs/operator-c2-http.md)

For all following requests, the operator will encrypt and sign all messages it sends to the C2 with its secret key. Similarly, the C2 server will encrypt all messages with the operators public key.

### Basic authentication flow

```ascii
Operator                                  C2
   |                                       +---+
   |            Send config file           |   |Generate config file,
   |<--------------------------------------+---+Save operator public key and logon secret in DB
   |         (name, sk, server pk)         |
   |                                       |
   |                                       |
   .                                       .
   .                                       .
   .                                       .
   |          Request auth token           |
   +-------------------------------------->|
   |          via signed message           |
   |                                       |
+--+---------------------------------------+--+
|  |                                       |  |Encrypted
|  |<--------------------------------------+  |
|  |            Return auth token,         |  |
|  |            RabbitMQ queue name        |  |
|  |                                       |  |
|  |            Send all requests          |  |
|  +-------------------------------------->|  |
|  |            with token in cookie       |  |
|  |                                       |  |
+--+---------------------------------------+--+
```

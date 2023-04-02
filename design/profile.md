# Maleable Profiles

## Introduction

In order to make the implant more suitable to different situations, we have implemented a maleable profile system. This allows the implant to be configured to behave differently from build to build. These options can also be changed on the fly, allowing the implant to be reconfigured without having to be rebuilt.

The default profile is located in [`server/comfig/profiles/default.yaml`](../server/config/profiles/default.yaml). This is the profile that is used if no profile is specified. The default profile is also used if the specified profile is not found.

## Profile Format

Profiles are YAML files with three main top level directives: `client`, `server`, and `routes`. The `client` directive contains the configuration options that are sent to the implant. The `server` directive contains the configuration options that are used by the server. The `routes` directive contains all HTTP endpoints to be used by both the client and server. Any options outside of these directoves are considered "global" and are used by both the client and server.

### Global options

| Option | Description | Required | Type |
|--------|-------------|----------|------|
| `name` | The name of this profile | Yes | String |
| `description` | A description of this profile | No | String |
| `implant_id_cookie` | The name of the cookie implants should use to identify themselves | Yes | String |

### Client options

| Option | Description | Required | Type |
|--------|-------------|----------|------|
| `user_agent` | The user agent to use when making HTTP requests | Yes | String |
| `encoding` | The encoding to use when sending encrypted data to the C2. | Yes | One of: `base64`, `hex` |
| `sleep` | The number of seconds to sleep between each HTTP request | Yes | Integer |
| `jitter` | % jitter. The implant will sleep for a random amount of time between `sleep` and `sleep * (1 + jitter)` | Yes | Float, `[0, 0.99]` |
| `max_retries` | The maximum number of times to retry a request before giving up | Yes | Integer |
| `auto_self_destruct` | Whether or not to self destruct on failed checkins. If set to true, the implant will delete itself after `max_retries` failed checkins. | Yes | Boolean |
| `retry_wait` | The number of seconds to wait before retrying a request. | Yes | Integer |
| `retry_jitter` | % jitter. The implant will wait for a random amount of time between `retry_wait` and `retry_wait * (1 + retry_jitter)` | Yes | Float, `[0, 0.99]` |
| `tailoring_hashes` | A list of hashes to use for payload tailoring. | Yes | List of strings |
| `tailoring_hash_function` | The hash function to use for payload tailoring. | Yes | One of: `sha256`, `md5` |
| `tailoring_hash_rounds` | The number of hash rounds to use for payload tailoring. | Yes | Integer |


### Server options

| Option | Description | Required | Type |
|--------|-------------|----------|------|
| `headers` | A String to String map of headers to send with each HTTP request | No | Map of strings |
| `redirect_url` | The URL to redirect to on the root path | Yes | String |

### Routes
See [this document](./specs/implant-c2-http.md) for more information on the routes.

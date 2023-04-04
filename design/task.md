# Tasks

Tasks are instructions, created by operators, to be sent to a given implant.
A task has the following fields:

| Field       | Type          |Description |
|-------------|---------------|------------|
| `operator_name` | `str` | The username of the operator who created this task|
| `task_id`     | `str` | The unique identifier for the task |
| `implant_id`  | `str` | The unique identifier for the implant that this task is assigned to |
| `opcode`      | `int` | The opcode associated with the task |
| `args`        | `Any` | The arguments associated with the task, represented as a JSON object. Usually a list or dict, depending on the type of task |
| `status`      | `str` | The status of the task [see here](#task-status) |
| `output`      | `str` | The output of the task, represented as a string |
| `created_at`  | `str` | The datetime the task was created, represented as a string in the format "YYYY-MM-DD HH:MM:SS.mmmmmm"   |
| `read_at`     | `str` | The datetime the task was read by the implant, represented as a string in the format "YYYY-MM-DD HH:MM:SS.mmmmmm" |
| `executed_at` | `str` | The datetime the task was executed by the implant, represented as a string in the format "YYYY-MM-DD HH:MM:SS.mmmmmm" |

## Task Status

A task can have the following statuses:

| Status       | Description |
|--------------|-------------|
| `CREATED`    | The task has been created, but not yet read by the implant |
| `TASKED`     | The task has been read by the implant, but not yet completed |
| `COMPLETED`  | The task has been completed by the implant |
| `ERROR`      | The task has been completed by the implant, but an error occurred during execution |

## Task output

Not all tasks will have output. For example, a task to list files in a directory will have output, but a task to sleep will not. If there is output, the implant should send it back to the appropriate endpoint as a Base64 encoded string. The operator can then decode the string and view the output.

**Even if a task has no output, it must still send a response to the appropriate endpoint.** This is to ensure that the operator knows that the task has been completed.

[Relevant documentation](specs/implant-c2-http.md#c2task)

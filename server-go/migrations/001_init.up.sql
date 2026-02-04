-- Operators table
CREATE TABLE IF NOT EXISTS operators (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    public_key TEXT NOT NULL,
    verify_key TEXT NOT NULL,
    login_secret TEXT NOT NULL,
    auth_token TEXT,
    auth_token_expiry TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login TIMESTAMPTZ,
    rmq_queue VARCHAR(255),
    role VARCHAR(50) NOT NULL DEFAULT 'operator',
    revoked BOOLEAN NOT NULL DEFAULT FALSE
);

-- Implants table
CREATE TABLE IF NOT EXISTS implants (
    id SERIAL PRIMARY KEY,
    implant_id VARCHAR(32) NOT NULL UNIQUE,
    hostname VARCHAR(255),
    ip VARCHAR(45),
    os VARCHAR(255),
    arch VARCHAR(50),
    username VARCHAR(255),
    server_sk TEXT NOT NULL,
    implant_pk TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Implant alias map
CREATE TABLE IF NOT EXISTS implant_aliases (
    id SERIAL PRIMARY KEY,
    implant_id VARCHAR(32) NOT NULL REFERENCES implants(implant_id) ON DELETE CASCADE,
    alias VARCHAR(255) NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_implant_aliases_implant_id ON implant_aliases(implant_id);

-- Implant configs table
CREATE TABLE IF NOT EXISTS implant_configs (
    id SERIAL PRIMARY KEY,
    implant_id VARCHAR(32) NOT NULL UNIQUE REFERENCES implants(implant_id) ON DELETE CASCADE,
    cookie VARCHAR(255) NOT NULL DEFAULT 'SESSID',
    kill_date TIMESTAMPTZ,
    user_agent TEXT,
    auto_self_destruct BOOLEAN NOT NULL DEFAULT FALSE,
    sleep_time INTEGER NOT NULL DEFAULT 60,
    jitter REAL NOT NULL DEFAULT 0.1,
    max_retries INTEGER NOT NULL DEFAULT -1,
    retry_wait INTEGER NOT NULL DEFAULT 0,
    retry_jitter REAL NOT NULL DEFAULT 0.0,
    enc_key TEXT NOT NULL,
    tailoring_hash_function VARCHAR(50) DEFAULT 'sha256',
    tailoring_hash_rounds INTEGER DEFAULT 1,
    tailoring_hashes TEXT[] DEFAULT '{}'
);

-- Tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    task_id VARCHAR(32) NOT NULL UNIQUE,
    operator_name VARCHAR(255) NOT NULL,
    implant_id VARCHAR(32) NOT NULL REFERENCES implants(implant_id) ON DELETE CASCADE,
    opcode INTEGER NOT NULL,
    args JSONB,
    status VARCHAR(50) NOT NULL DEFAULT 'CREATED',
    output TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read_at TIMESTAMPTZ,
    executed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_tasks_implant_id ON tasks(implant_id);
CREATE INDEX IF NOT EXISTS idx_tasks_operator_name ON tasks(operator_name);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);

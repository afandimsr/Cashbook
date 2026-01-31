CREATE TABLE IF NOT EXISTS oauth_states (
    id VARCHAR(36) PRIMARY KEY,
    state VARCHAR(128) NOT NULL UNIQUE,
    provider VARCHAR(32) NOT NULL,
    client_id VARCHAR(64) NULL,
    redirect_uri TEXT NULL,
    ip_hash VARCHAR(64) NULL,
    user_agent_hash VARCHAR(64) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL
);


-- ユーザーテーブル
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    hashed_password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- デバイステーブル
CREATE TABLE devices (
    device_id VARCHAR(100) PRIMARY KEY,
    janus_password VARCHAR(255) NOT NULL,
    device_name VARCHAR(255),
    location VARCHAR(255),
    status VARCHAR(50),
    last_communication TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ユーザーとデバイスの関連
CREATE TABLE user_devices (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    device_id VARCHAR(100) NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_device
        FOREIGN KEY (device_id)
            REFERENCES devices(device_id)
            ON DELETE CASCADE
);

-- Pushトークン
CREATE TABLE push_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    token VARCHAR(2048) NOT NULL,
    platform VARCHAR(20) NOT NULL,            -- 'ios' | 'android'
    last_seen_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked BOOLEAN DEFAULT FALSE,
    revoked_reason VARCHAR(255),

    CONSTRAINT fk_push_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX ux_push_tokens_token ON push_tokens(token);
CREATE INDEX ix_push_tokens_user ON push_tokens(user_id);
CREATE INDEX ix_push_tokens_valid ON push_tokens(user_id, revoked) WHERE revoked = FALSE;

-- エラーログ
CREATE TABLE errors (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(100) NOT NULL,
    error_code VARCHAR(50) NOT NULL,
    message VARCHAR(1024),
    error_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_error_device
        FOREIGN KEY (device_id)
            REFERENCES devices(device_id)
            ON DELETE CASCADE
);

-- アラート本体
CREATE TABLE alerts (
    id BIGSERIAL PRIMARY KEY,
    device_id VARCHAR(100) NOT NULL,
    event_id VARCHAR(128) NOT NULL,                 -- イベントID（全体でユニーク）
    severity VARCHAR(16) NOT NULL DEFAULT 'info',   -- 'info' | 'warn' | 'critical' など
    temp_c DOUBLE PRECISION,                        -- NULL許容
    occurred_at TIMESTAMPTZ NOT NULL,               -- 発生時刻 (UTC)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_alert_device
        FOREIGN KEY (device_id)
        REFERENCES devices(device_id)
        ON DELETE CASCADE
);

CREATE INDEX idx_alerts_device_occurred ON alerts (device_id, occurred_at DESC);

CREATE UNIQUE INDEX ux_alerts_event_id ON alerts (event_id);

-- アラートに紐づくメディア
CREATE TABLE alert_media (
    id BIGSERIAL PRIMARY KEY,
    alert_id BIGINT,                                -- pending行のため NULL許容
    object_key TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    size_bytes BIGINT,
    width INTEGER,
    height INTEGER,
    sha256 CHAR(64),
    committed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_alert_media_alert
        FOREIGN KEY (alert_id)
        REFERENCES alerts(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX ux_alert_media_object_key ON alert_media (object_key);

CREATE INDEX ix_alert_media_alert_id ON alert_media (alert_id);

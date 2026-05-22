CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    user_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    permission      VARCHAR(50)  NOT NULL DEFAULT 'member',
    username        VARCHAR(255) UNIQUE,
    password        VARCHAR(255),
    email           VARCHAR(255) UNIQUE,
    phone           VARCHAR(15),
    full_name       VARCHAR(255),
    gender          VARCHAR(50), ENUM('male', 'female', 'other'),
    date_of_birth   TIMESTAMP,
    facebook_url    VARCHAR(500) UNIQUE,
    github_url      VARCHAR(500) UNIQUE,
    avatar_url      VARCHAR(500), DEFAULT 'https://cdn.gauas.com/images/avatar/default_image.jpg',
    deleted_at      TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_username_permission ON users (username, permission);

CREATE TABLE IF NOT EXISTS verifications (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    method      VARCHAR(20) NOT NULL,
    value       VARCHAR(255) NOT NULL,
    is_verified BOOLEAN     NOT NULL DEFAULT FALSE,
    verified_at TIMESTAMP,
    UNIQUE (user_id, method, value)
);

CREATE INDEX IF NOT EXISTS idx_verifications_user_id ON verifications (user_id);
CREATE INDEX IF NOT EXISTS idx_verifications_method  ON verifications (method);

CREATE TABLE IF NOT EXISTS mfas (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    type        VARCHAR(30) NOT NULL,
    secret      VARCHAR(255),
    enabled     BOOLEAN     NOT NULL DEFAULT FALSE,
    verified_at TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_mfas_user_id ON mfas (user_id);
CREATE INDEX IF NOT EXISTS idx_mfas_type    ON mfas (type);

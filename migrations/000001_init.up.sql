CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY,
    permission      VARCHAR(50) NOT NULL DEFAULT 'member',
    username        VARCHAR(255) UNIQUE,
    full_name       VARCHAR(255),
    avatar_url      VARCHAR(500) DEFAULT 'https://cdn.gauas.com/images/avatar/default_image.jpg',
    date_of_birth   DATE,
    gender          VARCHAR(20) CHECK (gender IN ('male', 'female', 'other')),
    deleted_at      TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_username_permission ON users (username, permission);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);

CREATE TABLE IF NOT EXISTS identities (
    id                  UUID PRIMARY KEY,
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider            VARCHAR(50) NOT NULL,
    provider_user_id    VARCHAR(255) NOT NULL,
    email               VARCHAR(255),
    phone               VARCHAR(20),
    hash                VARCHAR(255),
    created_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE (provider, provider_user_id),
    UNIQUE (user_id, provider)
);

CREATE INDEX IF NOT EXISTS idx_identities_user_id ON identities(user_id);

CREATE INDEX IF NOT EXISTS idx_identities_provider ON identities(provider);

CREATE INDEX IF NOT EXISTS idx_identities_provider_user ON identities(provider, provider_user_id);

CREATE TABLE IF NOT EXISTS verifications (
    id              UUID PRIMARY KEY,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    method          VARCHAR(20) NOT NULL,
    value           VARCHAR(255) NOT NULL,
    is_verified     BOOLEAN NOT NULL DEFAULT FALSE,
    verified_at     TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE (user_id, method, value)
);

CREATE INDEX IF NOT EXISTS idx_verifications_user_id ON verifications(user_id);

CREATE INDEX IF NOT EXISTS idx_verifications_method ON verifications(method);

CREATE INDEX IF NOT EXISTS idx_verifications_user_method ON verifications(user_id, method);

CREATE TABLE IF NOT EXISTS mfas (
    id              UUID PRIMARY KEY,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type            VARCHAR(30) NOT NULL,
    secret          VARCHAR(255),
    enabled         BOOLEAN NOT NULL DEFAULT FALSE,
    verified_at     TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_mfas_user_id
    ON mfas(user_id);

CREATE INDEX IF NOT EXISTS idx_mfas_type
    ON mfas(type);
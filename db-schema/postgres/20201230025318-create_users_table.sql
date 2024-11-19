
-- +migrate Up
CREATE TABLE users (
	id BIGSERIAL PRIMARY KEY,
	uuid UUID DEFAULT uuid_generate_v4(),
    email VARCHAR NOT NULL DEFAULT '' UNIQUE,
	is_banned BOOLEAN NOT NULL DEFAULT FALSE,
    role_user VARCHAR NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

-- +migrate Down
DROP TABLE users;

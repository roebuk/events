CREATE TYPE user_role AS ENUM ('entrant', 'organizer', 'admin');
CREATE TYPE auth_provider AS ENUM ('google', 'apple');
CREATE TYPE audit_action AS ENUM ('created', 'updated', 'deleted');

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';


-- Users
CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    phone TEXT,
    address_line1 TEXT,
    address_line2 TEXT,
    city TEXT,
    state TEXT,
    postal_code TEXT,
    country TEXT,
    role user_role NOT NULL DEFAULT 'entrant',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NULL;

CREATE TRIGGER update_users_updated_at
  BEFORE UPDATE ON users
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();


-- Organisations
CREATE TABLE organisations (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_organisations_deleted_at ON organisations(deleted_at) WHERE deleted_at IS NULL;

CREATE TRIGGER update_organisations_updated_at
  BEFORE UPDATE ON organisations
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();


CREATE TABLE events (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  organisation_id BIGINT NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  slug TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_events_slug ON events(slug);
CREATE INDEX idx_events_deleted_at ON events(deleted_at) WHERE deleted_at IS NULL;

CREATE TRIGGER update_events_updated_at
  BEFORE UPDATE ON events
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE races (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  slug TEXT NOT NULL,
  registration_open_date TIMESTAMPTZ,
  registration_close_date TIMESTAMPTZ,
  max_capacity INT NOT NULL CHECK (max_capacity > 0),
  price_units INT CHECK (price_units >= 0),
  currency TEXT DEFAULT 'GBP',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  UNIQUE(event_id, slug)
);

CREATE INDEX idx_races_event_id ON races(event_id);
CREATE INDEX idx_races_slug ON races(event_id, slug);
CREATE INDEX idx_races_deleted_at ON races(deleted_at) WHERE deleted_at IS NULL;

CREATE TRIGGER update_races_updated_at
  BEFORE UPDATE ON races
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();


-- Organisation Users (many-to-many relationship)
CREATE TABLE organisation_users (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  organisation_id BIGINT NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  UNIQUE(organisation_id, user_id)
);


CREATE INDEX idx_organisation_users_org_id ON organisation_users(organisation_id);
CREATE INDEX idx_organisation_users_user_id ON organisation_users(user_id);
CREATE INDEX idx_organisation_users_deleted_at ON organisation_users(deleted_at) WHERE deleted_at IS NULL;

-- Auth Credentials
CREATE TABLE auth_credentials (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    password_hash TEXT NOT NULL,
    email_verified_at TIMESTAMPTZ,
    last_login_at TIMESTAMPTZ,
    failed_login_attempts INT NOT NULL DEFAULT 0,
    locked_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_auth_credentials_user_id ON auth_credentials(user_id);
CREATE INDEX idx_auth_credentials_deleted_at ON auth_credentials(deleted_at) WHERE deleted_at IS NULL;

CREATE TRIGGER update_auth_credentials_updated_at
  BEFORE UPDATE ON auth_credentials
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();


CREATE TABLE social_accounts (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider auth_provider NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL, -- ID from the provider
    -- access_token TEXT, -- Optional: store if you need to make API calls
    -- refresh_token TEXT, -- Optional: for refreshing access tokens
    -- expires_at TIMESTAMP, -- Optional: token expiration
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(provider, provider_user_id)
);

CREATE INDEX idx_social_accounts_user_id ON social_accounts(user_id);
CREATE INDEX idx_social_accounts_deleted_at ON social_accounts(deleted_at) WHERE deleted_at IS NULL;

CREATE TRIGGER update_social_accounts_updated_at
  BEFORE UPDATE ON social_accounts
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- Sessions table for SCS PostgreSQL store
CREATE TABLE sessions (
    token TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);




-- CREATE TABLE audit_log (
--     id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--     table_name TEXT NOT NULL,
--     record_id BIGINT NOT NULL,
--     action audit_action NOT NULL,
--     user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
--     changed_fields JSONB, -- Store what changed: {"name": {"old": "X", "new": "Y"}}
--     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
-- );

-- CREATE INDEX idx_audit_log_table_record ON audit_log(table_name, record_id);
-- CREATE INDEX idx_audit_log_user_id ON audit_log(user_id);
-- CREATE INDEX idx_audit_log_created_at ON audit_log(created_at DESC);

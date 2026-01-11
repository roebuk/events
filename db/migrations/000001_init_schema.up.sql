-- Organizations table
CREATE TABLE organizations (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT,
    website TEXT,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_organizations_deleted_at ON organizations(deleted_at) WHERE deleted_at IS NULL;

-- Users table
CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    phone TEXT,
    address_line1 TEXT,
    address_line2 TEXT,
    city TEXT,
    state TEXT,
    postal_code TEXT,
    country TEXT,
    role TEXT NOT NULL CHECK (role IN ('entrant', 'organizer', 'admin')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NULL;

-- Organization Users junction table (many-to-many)
CREATE TABLE organization_users (
    organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (organization_id, user_id)
);

CREATE INDEX idx_organization_users_user_id ON organization_users(user_id);

-- Event Types lookup table
CREATE TABLE event_types (
    id BIGSERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert default event types
INSERT INTO event_types (name, description) VALUES
    ('10k', '10 kilometer running race'),
    ('half_marathon', 'Half marathon (21.0975 km) running race'),
    ('marathon', 'Full marathon (42.195 km) running race'),
    ('cycling', 'Cycling event'),
    ('triathlon', 'Multi-sport event combining swimming, cycling, and running'),
    ('sprint_triathlon', 'Sprint distance triathlon'),
    ('olympic_triathlon', 'Olympic distance triathlon'),
    ('ironman', 'Ironman distance triathlon');

-- Events table
CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE RESTRICT,
    event_type_id BIGINT NOT NULL REFERENCES event_types(id) ON DELETE RESTRICT,
    parent_event_id BIGINT REFERENCES events(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    location TEXT NOT NULL,
    venue_name TEXT,
    city TEXT NOT NULL,
    state TEXT,
    country TEXT NOT NULL,
    event_date TIMESTAMPTZ NOT NULL,
    registration_open_date TIMESTAMPTZ,
    registration_close_date TIMESTAMPTZ,
    max_capacity INT NOT NULL CHECK (max_capacity > 0),
    current_entries INT NOT NULL DEFAULT 0 CHECK (current_entries >= 0 AND current_entries <= max_capacity),
    price_cents INT CHECK (price_cents >= 0),
    currency TEXT DEFAULT 'USD',
    image_urls TEXT[],
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_events_organization_id ON events(organization_id);
CREATE INDEX idx_events_event_type_id ON events(event_type_id);
CREATE INDEX idx_events_parent_event_id ON events(parent_event_id) WHERE parent_event_id IS NOT NULL;
CREATE INDEX idx_events_event_date ON events(event_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_events_location ON events(city, country) WHERE deleted_at IS NULL;
CREATE INDEX idx_events_published ON events(is_published, event_date) WHERE deleted_at IS NULL AND is_published = TRUE;

-- Entrants table (event registrations)
CREATE TABLE entrants (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'registered' CHECK (status IN ('registered', 'confirmed', 'cancelled', 'completed', 'dns', 'dnf')),
    bib_number INT,
    registration_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finish_time INTERVAL,
    emergency_contact_name TEXT,
    emergency_contact_phone TEXT,
    special_requirements TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(user_id, event_id)
);

CREATE INDEX idx_entrants_user_id ON entrants(user_id);
CREATE INDEX idx_entrants_event_id ON entrants(event_id);
CREATE INDEX idx_entrants_status ON entrants(event_id, status);
CREATE INDEX idx_entrants_bib_number ON entrants(event_id, bib_number) WHERE bib_number IS NOT NULL;

-- Updated at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply updated_at triggers
CREATE TRIGGER update_organizations_updated_at BEFORE UPDATE ON organizations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_events_updated_at BEFORE UPDATE ON events
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_entrants_updated_at BEFORE UPDATE ON entrants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

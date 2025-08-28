DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role_status') THEN
        CREATE TYPE role_status AS ENUM ('active', 'inactive');
    END IF;
END$$;

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    variant VARCHAR(255),
    status role_status DEFAULT 'active',
    created_by UUID DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_roles_updated_at BEFORE
UPDATE ON roles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column ();
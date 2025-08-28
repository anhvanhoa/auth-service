DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'module_child_status') THEN
        CREATE TYPE module_child_status AS ENUM ('active', 'inactive');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'module_child_method') THEN
        CREATE TYPE module_child_method AS ENUM (
            'GET',
            'POST',
            'PUT',
            'DELETE',
            'PATCH',
            'HEAD',
            'OPTIONS'
            );
    END IF;
END$$;


CREATE TABLE
    module_childs (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        module_id UUID NOT NULL,
        path TEXT NOT NULL,
        method module_child_method NOT NULL,
        name VARCHAR(255),
        is_private BOOLEAN DEFAULT FALSE,
        status module_child_status DEFAULT 'active',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (module_id) REFERENCES modules (id)
    );

CREATE TRIGGER update_module_childs_updated_at BEFORE
UPDATE ON module_childs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column ();

CREATE INDEX idx_module_childs_module_id ON module_childs(module_id);

CREATE UNIQUE INDEX idx_module_childs_module_id_path_method ON module_childs(module_id, path, method);
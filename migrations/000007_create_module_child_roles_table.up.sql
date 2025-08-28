CREATE TABLE
    module_child_roles (
        role_id UUID NOT NULL,
        module_child_id UUID NOT NULL,
        created_by UUID NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (role_id, module_child_id),
        FOREIGN KEY (role_id) REFERENCES roles (id),
        FOREIGN KEY (module_child_id) REFERENCES module_childs (id)
    );

CREATE TRIGGER update_module_child_roles_updated_at BEFORE
UPDATE ON module_child_roles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column ();

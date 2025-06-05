-- 1. ENUM for role member types
CREATE TYPE role_member_role AS ENUM ('admin', 'standard');

-- 2. Users table (assumed to exist)
-- CREATE TABLE users (...)

-- 3. Roles table
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,  -- short description of what the role is for
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Role Members table
CREATE TABLE role_members (
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    added_by INTEGER NOT NULL REFERENCES users(id),
    role role_member_role NOT NULL DEFAULT 'standard',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, user_id)
);

-- Trigger function for INSERT
CREATE OR REPLACE FUNCTION update_roles_on_member_add()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE roles
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.role_id;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for DELETE
CREATE OR REPLACE FUNCTION update_roles_on_member_remove()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE roles
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.role_id;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;


-- Triggers
CREATE TRIGGER after_add_member
AFTER INSERT ON role_members
FOR EACH ROW
EXECUTE FUNCTION update_roles_on_member_add();

CREATE TRIGGER after_remove_member
AFTER DELETE ON role_members
FOR EACH ROW
EXECUTE FUNCTION update_roles_on_member_remove();
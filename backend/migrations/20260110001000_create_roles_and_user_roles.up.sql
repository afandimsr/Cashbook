-- Create roles and user_roles tables
CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- ensure base roles exist
INSERT INTO roles (name) VALUES ('ADMIN') ON CONFLICT DO NOTHING;
INSERT INTO roles (name) VALUES ('USER') ON CONFLICT DO NOTHING;

-- migrate existing users.role values into roles and user_roles
INSERT INTO roles (name)
SELECT DISTINCT role FROM users WHERE role IS NOT NULL AND role <> ''
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u
JOIN roles r ON r.name = u.role
WHERE u.role IS NOT NULL AND u.role <> ''
ON CONFLICT DO NOTHING;

-- remove denormalized role column from users
ALTER TABLE users DROP COLUMN IF EXISTS role;

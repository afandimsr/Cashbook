-- revert: add role column back (take first assigned role) and drop join tables
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(100) DEFAULT 'USER';

UPDATE users SET role = sub.rname FROM (
  SELECT ur.user_id, r.name AS rname,
         ROW_NUMBER() OVER (PARTITION BY ur.user_id ORDER BY ur.role_id) rn
  FROM user_roles ur
  JOIN roles r ON r.id = ur.role_id
) sub WHERE users.id = sub.user_id AND sub.rn = 1;

-- drop join and roles
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS roles;

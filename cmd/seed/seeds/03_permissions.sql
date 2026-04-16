INSERT INTO users_permissions (user_id, permission_id)
SELECT u.id, p.id
FROM users u, permissions p
WHERE u.email = 'user@gearboxd.com'
  AND p.code = 'cars:read';

INSERT INTO users_permissions (user_id, permission_id)
SELECT u.id, p.id
FROM users u, permissions p
WHERE u.email = 'admin@gearboxd.com'
  AND p.code IN ('cars:read', 'cars:write');
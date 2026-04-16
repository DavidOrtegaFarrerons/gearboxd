BEGIN;

TRUNCATE TABLE users, permissions RESTART IDENTITY CASCADE;

INSERT INTO permissions (code)
VALUES
    ('cars:read'),
    ('cars:write');


INSERT INTO users (email, username, password_hash, activated, version)
VALUES (
       'user@gearboxd.com',
       'user',
       '\x24326124313224423034416d444d62704663365868745741793364512e766148474330664136434a366f4438457374592e506b612f684d764875792e',
       true,
       1
       ),
       (
           'admin@gearboxd.com',
           'admin',
           '\x24326124313224423034416d444d62704663365868745741793364512e766148474330664136434a366f4438457374592e506b612f684d764875792e',
           true,
           1
       )
    ;

COMMIT;
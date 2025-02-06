INSERT INTO casbin_rule (ptype, v0, v1, v2)
VALUES ('p', 'admin', '/admin', 'GET'),
       ('p', 'user', '/v1/users/', 'GET'),
       ('p', 'user', '/v1/users/login', 'POST');

INSERT INTO casbin_rule (ptype, v0, v1)
VALUES ('g', 'alice', 'admin'),
       ('g', 'bob', 'user');
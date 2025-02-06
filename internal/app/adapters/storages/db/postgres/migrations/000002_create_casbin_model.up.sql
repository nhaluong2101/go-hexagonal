CREATE TABLE casbin_model
(
    id         SERIAL PRIMARY KEY,           -- ID tự tăng
    model_name VARCHAR(255) UNIQUE NOT NULL, -- Tên của model, ví dụ: "rbac_model"
    model_text TEXT                NOT NULL  -- Lưu nội dung RBAC model
);

CREATE TABLE IF NOT EXISTS casbin_rule (
                                           id SERIAL PRIMARY KEY,
                                           ptype VARCHAR(100) NOT NULL,
    v0 VARCHAR(100),
    v1 VARCHAR(100),
    v2 VARCHAR(100),
    v3 VARCHAR(100),
    v4 VARCHAR(100),
    v5 VARCHAR(100)
    );

INSERT INTO casbin_model (model_name, model_text)
VALUES ('rbac_model',
        '[request_definition]
        r = sub, obj, act

        [policy_definition]
        p = sub, obj, act

        [role_definition]
        g = _, _

        [policy_effect]
        e = some(where (p.eft == allow))

        [matchers]
        m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act');

-- 新增角色
INSERT INTO auths (`id`, `name`) values (13,"管理本机构信息");
INSERT INTO roles (`id`, `name`, `created_at`, `updated_at`) values (8,"高级机构账号", now(), now());
insert into role_auths (role_id, auth_id) values(8, 6), (8, 13);
-- 增加机构支持角色
UPDATE orgs SET support_role_ids = "1,2,3,4,5,6" WHERE id = 1;
UPDATE orgs SET support_role_ids = "7,8" WHERE id != 1;
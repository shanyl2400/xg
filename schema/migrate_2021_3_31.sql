UPDATE order_pay_records SET real_price=amount WHERE mode = 1;
UPDATE order_pay_records SET real_price=-amount WHERE mode = 2;

UPDATE orders AS o SET parent_org_id = (SELECT IF(parent_id=0,id,parent_id) AS real_id FROM orgs WHERE orgs.id=o.org_id)

SET GLOBAL sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''));

UPDATE auths SET mode = 1 WHERE id not in (6,13);
UPDATE auths SET mode = 2 WHERE id in (6,13);
INSERT INTO xg2.role_auths(role_id, auth_id) VALUES(1,14);
INSERT INTO xg2.role_auths(role_id, auth_id) VALUES(6,14);


ALTER TABLE `xg2`.`commission_settlement_records` 
CHANGE COLUMN `settlement_note` `settlement_note` TEXT CHARACTER SET 'utf8' COLLATE 'utf8_general_ci' NULL DEFAULT NULL ;


ALTER TABLE `xg2`.`commission_settlement_records` 
CHANGE COLUMN `note` `note` TEXT CHARACTER SET 'utf8' COLLATE 'utf8_general_ci' NULL DEFAULT NULL ;
INSERT INTO auths (id, name, mode) VALUES (14, "结算", 1);
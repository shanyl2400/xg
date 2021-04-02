UPDATE order_pay_records SET real_price=amount WHERE mode = 1;
UPDATE order_pay_records SET real_price=-amount WHERE mode = 2;

UPDATE orders AS o SET parent_org_id = (SELECT IF(parent_id=0,id,parent_id) AS real_id FROM orgs WHERE orgs.id=o.org_id)

SET GLOBAL sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''));

UPDATE auths SET mode = 1 WHERE id not in (6,13);
UPDATE auths SET mode = 2 WHERE id in (6,13);
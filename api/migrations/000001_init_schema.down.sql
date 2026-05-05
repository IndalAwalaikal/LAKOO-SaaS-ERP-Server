ALTER TABLE sale_items DROP FOREIGN KEY fk_sale_items_sale_id;
ALTER TABLE sale_items DROP FOREIGN KEY fk_sale_items_product_id;
ALTER TABLE users DROP FOREIGN KEY fk_user_tenant;
ALTER TABLE products DROP FOREIGN KEY fk_product_tenant;
ALTER TABLE sales DROP FOREIGN KEY fk_sale_tenant;

DROP TABLE IF EXISTS sale_items;
DROP TABLE IF EXISTS sales;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tenants;

CREATE TABLE restaritems (
	id	BIGSERIAL PRIMARY KEY,
	onceGUID text NOT NULL,
	name text	NOT NULL,
	sku text NOT NULL,
	itemGUID text NOT NULL
);
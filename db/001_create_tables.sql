CREATE TABLE "Population".Person (
	id bigserial NOT NULL, 
	name varchar(100) NOT NULL, 
	surname varchar(100) NOT NULL,
	patronymic varchar(100),
	age int,
	-- todo: сделать внешним ключом
	country_id  varchar(100),
	-- todo: сделать внешним ключом
	gender_id varchar(100),
		
	creaed_at timestamptz  default now() NOT NULL,
	
	CONSTRAINT pk_keyword_id2 PRIMARY KEY (id, age)
) PARTITION BY RANGE (age);



CREATE TABLE "Population".Person_0_20 PARTITION OF "Population".Person FOR VALUES FROM (0) TO (20);
CREATE TABLE "Population".Person_20_40 PARTITION OF "Population".Person FOR VALUES FROM (20) TO (40);
CREATE TABLE "Population".Person_40_100 PARTITION OF "Population".Person FOR VALUES FROM (40) TO (100);
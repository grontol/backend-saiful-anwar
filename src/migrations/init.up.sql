CREATE TABLE yards (
	id int4 GENERATED ALWAYS AS IDENTITY NOT NULL,
	"name" varchar(255) NOT NULL,
	description text NOT NULL,
	CONSTRAINT yards_pk PRIMARY KEY (id),
	CONSTRAINT yards_unique UNIQUE (name)
);

CREATE TABLE blocks (
	id int4 GENERATED ALWAYS AS IDENTITY NOT NULL,
	yard_id int4 NOT NULL,
	"name" varchar(255) NOT NULL,
	slots int4 NOT NULL,
	"rows" int4 NOT NULL,
	tiers int4 NOT NULL,
	CONSTRAINT blocks_pk PRIMARY KEY (id),
	CONSTRAINT blocks_unique UNIQUE (name)
);

CREATE TABLE yard_plans (
	id int4 GENERATED ALWAYS AS IDENTITY NOT NULL,
	block_id int4 NOT NULL,
	slot_start int4 NOT NULL,
	slot_end int4 NOT NULL,
	row_start int4 NOT NULL,
	row_end int4 NOT NULL,
	"size" int4 NOT NULL,
	height float4 NOT NULL,
	"type" varchar NOT NULL,
	slot_priority int4 NOT NULL,
	row_priority int4 NOT NULL,
	tier_priority int4 NOT NULL,
	CONSTRAINT yard_plans_pk PRIMARY KEY (id)
);

CREATE TABLE placements (
	id int4 GENERATED ALWAYS AS IDENTITY NOT NULL,
	container_id varchar NOT NULL,
	block_id int4 NOT NULL,
	slot int4 NOT NULL,
	"row" int4 NOT NULL,
	tier int4 NOT NULL,
	is_head bool NOT NULL,
	container_size int4 NOT NULL,
	container_height float4 NOT NULL,
	container_type varchar NOT NULL,
	CONSTRAINT placements_pk PRIMARY KEY (id),
	CONSTRAINT placements_unique UNIQUE (container_id, is_head),
	CONSTRAINT placements_blocks_fk FOREIGN KEY (block_id) REFERENCES blocks(id)
);
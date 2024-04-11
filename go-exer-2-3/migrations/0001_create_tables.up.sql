-- public.credentials definition

-- Drop table

-- DROP TABLE public.credentials;

CREATE TABLE public.credentials (
	id serial4 PRIMARY KEY NOT NULL,
	"user" varchar(20) NULL,
	url varchar(20) NULL,
	username varchar(20) NULL,
	"password" varchar(20) NULL
);
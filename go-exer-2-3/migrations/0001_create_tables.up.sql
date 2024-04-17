create table if not exists credentials (
	id bigserial primary key,
	url varchar(40) NULL,
	username varchar(20) NULL,
	"password" varchar(20) NULL,
    created_at timestamp(0) with time zone not null default now()
);

create table if not exists users (
    id bigserial primary key,
    name text not null,
    username text not null,
    password text not null,
    created_at timestamp(0) with time zone not null default now()
);

create table if not exists user_credential (
    user_id bigserial not null references users (id) on delete cascade,
    credential_id bigserial not null references credentials (id) on delete cascade,
    constraint user_credential_key primary key (user_id, credential_id)
);
create schema if not exists billing;

create table if not exists billing.users
(
	id bigserial not null constraint customer_pkey primary key,
    balance  NUMERIC(10,2)   not null default 0 check ( balance >= 0 )

);


create table if not exists branches
(
    id   serial      not null primary key,
    name varchar(32) not null,
    city varchar(32) not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    deleted_at timestamp
);

ALTER TABLE users_client
    ADD COLUMN if not exists branch_id INT default null references branches(id);

CREATE TYPE gen AS ENUM ('m', 'f');

create table if not exists sections
(
    id         serial      not null primary key,
    name       varchar(64) not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    deleted_at timestamp
);

create table if not exists roles
(
    id         uuid        not null primary key,
    name       varchar(96) not null,
    section_id int         not null references sections (id),
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    deleted_at timestamp
);

create table if not exists users_client
(
    id         uuid not null primary key,
    spec_id    varchar(10) not null,
    first_name varchar(32) not null,
    last_name  varchar(32) not null,
    email      varchar(60),
    password   varchar(60),
    cover      varchar(100),
    birthday   timestamp,
    gender     gen,
    added_at   timestamp,
    role_id    uuid        not null references roles (id),
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    deleted_at timestamp
);

-- create table if not exists behaviors
-- (
--     id         uuid        not null primary key,
--     name       varchar(32) not null,
--     created_at timestamp default current_timestamp,
--     updated_at timestamp default current_timestamp,
--     deleted_at timestamp
-- );

create table if not exists evaluations
(
    id          uuid         not null primary key,
    content     varchar(100) not null,
    star        smallint     not null,
--     behavior_id uuid         not null references behaviors (id),
    created_at  timestamp default current_timestamp,
    updated_at  timestamp default current_timestamp,
    deleted_at  timestamp
);

create table if not exists rated
(
    id         uuid not null primary key,
    additional varchar(512),
    user_id    uuid not null references users_client (id),
    created_at timestamp default current_timestamp,
    deleted_at timestamp
);

create table if not exists rated_evaluations
(
    id            serial not null primary key,
    rated_id      uuid   not null references rated (id),
    evaluation_id uuid   not null references evaluations (id)
);
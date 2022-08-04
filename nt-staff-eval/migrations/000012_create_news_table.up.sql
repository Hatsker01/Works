create table categories(
   id uuid not null primary key,
   name varchar(32) not null,
   created_at timestamp default current_timestamp,
   updated_at timestamp default current_timestamp,
   deleted_at timestamp
);

create table news
(
    id         uuid         not null primary key,
    title      text         not null,
    body       text         not null,
    cover      varchar(512) not null,
    author     uuid references users_client (id),
    read_time   varchar(32)  not null,
    category_id uuid references categories (id),
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    deleted_at timestamp
);
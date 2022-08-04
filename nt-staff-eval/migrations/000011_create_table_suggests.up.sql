--postgres
create type status_suggest as enum ('new','active','inactive');

create table suggests
(
    id         uuid           not null primary key,
    user_id    uuid references users_client (id),
    content    text           not null,
    status     status_suggest not null default 'new',
    created_at timestamp               default current_timestamp,
    updated_at timestamp               default current_timestamp,
    deleted_at timestamp
);
create table if not exists users_auth
(
    id            uuid not null primary key,
    access_token  text,
    refresh_token text,
    user_id       uuid references users_client (id)
);
create table if not exists admin(
    id serial primary key,
    login varchar(16),
    password varchar(60),
    access_token text,
    refresh_token text
);

update admin set access_token='a', refresh_token='b' where id=1;
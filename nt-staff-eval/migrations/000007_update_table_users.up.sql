alter table users_client
    add column if not exists phone varchar(17),
    add column if not exists address varchar(64),
    add column if not exists work_time varchar(32),
    add column if not exists social_medias jsonb,
    add column if not exists additional_informs jsonb;
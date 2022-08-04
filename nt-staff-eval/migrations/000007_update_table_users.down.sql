alter table users_client
    drop column if  exists phone,
    drop column if  exists address,
    drop column if  exists work_time,
    drop column if  exists social_medias,
    drop column if  exists additional_inform;

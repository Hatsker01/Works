CREATE TYPE eval AS ENUM ('staff', 'customer');

alter table evaluations add column if not exists section_id int references sections(id);
alter table evaluations add column if not exists eval_type eval default 'customer';

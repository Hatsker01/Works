alter table evaluations drop column if exists eval_type;
alter table evaluations drop column if exists section_id;
drop type if exists eval;

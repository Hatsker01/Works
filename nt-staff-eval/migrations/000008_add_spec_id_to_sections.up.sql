ALTER TABLE sections ADD COLUMN if not exists spec_id serial;

UPDATE sections SET spec_id = id;
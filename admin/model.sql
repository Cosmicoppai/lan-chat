CREATE TABLE IF NOT EXISTS videos(
    id uuid DEFAULT gen_random_uuid(),
    typ varchar(12) NOT NULL ,
    name varchar(50) NOT NULL ,
    epNo smallint NOT NULL ,
    imageLink varchar(80) NOT NULL,
    videoLink varchar(80) NOT NULL ,
    subLink varchar(80),
    createdAt timestamptz DEFAULT current_timestamp,

    PRIMARY KEY (id)
);
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gin;
CREATE INDEX IF NOT EXISTS nameIdx ON videos USING GIN (name);
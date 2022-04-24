CREATE SCHEMA IF NOT EXISTS lan_show;

SET SEARCH_PATH TO lan_show;


CREATE TABLE IF NOT EXISTS users (
    username varchar(20) PRIMARY KEY ,
    password varchar(300) NOT NULL,
    isAdmin bool NOT NULL DEFAULT FALSE
);


CREATE TABLE IF NOT EXISTS show_type (
    id smallserial PRIMARY KEY ,
    typ varchar(30) NOT NULL UNIQUE
);

INSERT INTO show_type (typ) VALUES ('movie'), ( 'series'), ( 'ova') ON CONFLICT DO NOTHING;


CREATE TABLE IF NOT EXISTS shows (
    id smallserial PRIMARY KEY ,
    name varchar(150) NOT NULL,
    totalEps smallint NOT NULL CHECK ( totalEps > 0 ),
    typ varchar(30) NOT NULL ,
    FOREIGN KEY(typ) REFERENCES show_type(typ) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS videos (
    id smallserial PRIMARY KEY ,
    videoLink varchar(150) NOT NULL,
    posterLink varchar(150) NOT NULL,
    subLink varchar(150) ,
    showId smallserial NOT NULL  ,
    epNo smallint NOT NULL CHECK ( epNo > 0 ) ,
    createdAt timestamptz DEFAULT current_timestamp,
    FOREIGN KEY (showId) REFERENCES shows(id) ON DELETE CASCADE
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gin;
CREATE INDEX IF NOT EXISTS nameIdx ON shows USING GIN (name);
-- create_tables.sql

CREATE DATABASE choicemovers;
\c choicemovers;
CREATE TABLE mytable (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

insert into mytable (name) values ('dakota');
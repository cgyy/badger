CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
  id serial,
  email varchar(40) NOT NULL,
  username varchar(40) NOT NULL,
  password varchar NOT NULL
);

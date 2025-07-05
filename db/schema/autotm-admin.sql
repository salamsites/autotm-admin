CREATE ROLE autotm LOGIN PASSWORD 'autotm';

CREATE DATABASE autotm_admin
       WITH
       OWNER = autotm
       ENCODING = 'UTF8'
       CONNECTION LIMIT = -1
       IS_TEMPLATE = False;

GRANT ALL PRIVILEGES ON DATABASE autotm_admin TO postgres;

\c autotm_admin;

SET client_encoding TO 'UTF-8';


CREATE TABLE IF NOT EXISTS brands (
                        id SERIAL PRIMARY KEY,
                        name CHARACTER VARYING(255),
                        logo_path CHARACTER VARYING(255),
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS brand_models (
                              id SERIAL PRIMARY KEY,
                              name CHARACTER VARYING(255),
                              logo_path CHARACTER VARYING(255),
                              brand_id INTEGER,
                              created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                              updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                              CONSTRAINT brand_id_fk
                                  FOREIGN KEY (brand_id)
                                      REFERENCES brands(id)
                                      ON UPDATE CASCADE ON DELETE SET NULL
);
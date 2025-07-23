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


CREATE TYPE category_type AS ENUM ('auto', 'moto', 'truck');


CREATE TABLE IF NOT EXISTS body_types (
            "id" SERIAL PRIMARY KEY,
            "name_tm" CHARACTER VARYING(255) NOT NULL,
            "name_ru" CHARACTER VARYING(255) NOT NULL,
            "name_en" CHARACTER VARYING(255) NOT NULL,
            "image_path" CHARACTER VARYING(255),
            "category" category_type,
            "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS brands (
            "id" SERIAL PRIMARY KEY,
            "name" CHARACTER VARYING(255) NOT NULL,
            "logo_path" CHARACTER VARYING(255),
            "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS brand_categories (
    "brand_id" INTEGER NOT NULL,
    "category" category_type NOT NULL,
    PRIMARY KEY (brand_id, category),
    CONSTRAINT brand_id_fk
        FOREIGN KEY (brand_id)
            REFERENCES brands(id)
                ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS models (
            "id" SERIAL PRIMARY KEY,
            "name" CHARACTER VARYING(255) NOT NULL,
            "brand_id" INTEGER NOT NULL,
            "body_type_id" INTEGER NOT NULL,
            "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            CONSTRAINT brand_id_fk
                FOREIGN KEY (brand_id)
                    REFERENCES brands(id)
                        ON UPDATE CASCADE ON DELETE CASCADE,
            CONSTRAINT body_type_id_fk
                FOREIGN KEY (body_type_id)
                    REFERENCES body_types(id)
                        ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS roles (
            "id" SERIAL PRIMARY KEY,
            "name" CHARACTER VARYING(255) NOT NULL,
            "role" json,
            "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
            "id" SERIAL PRIMARY KEY,
            "username" CHARACTER VARYING(255) NOT NULL UNIQUE,
            "login" CHARACTER VARYING(255) NOT NULL UNIQUE,
            "password" TEXT NOT NULL,
            "role_id" INTEGER,
            "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            CONSTRAINT role_id_fk
                FOREIGN KEY (role_id)
                     REFERENCES roles(id)
                        ON UPDATE CASCADE ON DELETE SET NULL
);


CREATE TABLE IF NOT EXISTS regions (
            "id" SERIAL PRIMARY KEY,
            "name_tm" CHARACTER VARYING(255) NOT NULL,
            "name_en" CHARACTER VARYING(255) NOT NULL,
            "name_ru" CHARACTER VARYING(255) NOT NULL,
            "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cities (
            "id" SERIAL PRIMARY KEY,
            "name_tm" CHARACTER VARYING(255) NOT NULL,
            "name_en" CHARACTER VARYING(255) NOT NULL,
            "name_ru" CHARACTER VARYING(255) NOT NULL,
            "region_id" INTEGER,
            "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            CONSTRAINT region_id_fk
                FOREIGN KEY (region_id)
                    REFERENCES regions(id)
                        ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS sliders (
            "id" SERIAL PRIMARY KEY,
            "image_path_tm" CHARACTER VARYING(255) NOT NULL,
            "image_path_en" CHARACTER VARYING(255) NOT NULL,
            "image_path_ru" CHARACTER VARYING(255) NOT NULL,
            "platform" CHARACTER VARYING(100) NOT NULL,
            "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
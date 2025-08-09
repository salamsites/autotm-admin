-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'category_type') THEN
CREATE TYPE category_type AS ENUM ('auto', 'moto', 'truck');
END IF;
END$$;
-- +goose StatementEnd

-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'stock_status') THEN
CREATE TYPE stock_status AS ENUM ('waiting', 'accepted', 'blocked');
END IF;
END
$$;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS body_types (
            "id" SERIAL PRIMARY KEY,
            "name_tm" CHARACTER VARYING(255) NOT NULL,
            "name_ru" CHARACTER VARYING(255) NOT NULL,
            "name_en" CHARACTER VARYING(255) NOT NULL,
            "image_path" JSONB,
            "upload_id" UUID,
            "category" category_type,
            "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS brands (
                "id" SERIAL PRIMARY KEY,
                "name" CHARACTER VARYING(255) NOT NULL,
                "logo_path" JSONB,
                "upload_id" UUID,
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
                "category" category_type NOT NULL,
                "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                CONSTRAINT brand_id_fk
                    FOREIGN KEY (brand_id)
                        REFERENCES brands(id)
                           ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS roles (
                "id" SERIAL PRIMARY KEY,
                "name" CHARACTER VARYING(255) NOT NULL,
                "role" json,
                "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS admin_users (
                "id" SERIAL PRIMARY KEY,
                "username" CHARACTER VARYING(255) NOT NULL,
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
                "image_path_tm" JSONB,
                "image_path_en" JSONB,
                "image_path_ru" JSONB,
                "upload_id_tm" UUID,
                "upload_id_en" UUID,
                "upload_id_ru" UUID,
                "platform" CHARACTER VARYING(100) NOT NULL,
                "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS stocks (
                "id" SERIAL PRIMARY KEY,
                "user_id" BIGINT,
                "phone_number" CHARACTER VARYING(255),
                "email" CHARACTER VARYING(255),
                "store_name" CHARACTER VARYING(255),
                "images" TEXT[],
                "logo" CHARACTER VARYING(255),
                "region_id" INTEGER,
                "city_id" INTEGER,
                "address" TEXT,
                "description" TEXT,
                "location" TEXT,
                "status" stock_status DEFAULT 'waiting',
                "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                CONSTRAINT region_id_fk
                    FOREIGN KEY (region_id)
                    REFERENCES regions(id)
                        ON UPDATE CASCADE ON DELETE SET NULL,
                CONSTRAINT city_id_fk
                    FOREIGN KEY (city_id)
                        REFERENCES cities(id)
                            ON UPDATE CASCADE ON DELETE SET NULL
);


CREATE TABLE IF NOT EXISTS descriptions (
                        "id" SERIAL PRIMARY KEY,
                        "name_tm" TEXT,
                        "name_en" TEXT,
                        "name_ru" TEXT,
                        "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sliders;
DROP TABLE IF EXISTS cities;
DROP TABLE IF EXISTS regions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS models;
DROP TABLE IF EXISTS brand_categories;
DROP TABLE IF EXISTS brands;
DROP TABLE IF EXISTS body_types;
DROP TABLE IF EXISTS stocks;
DROP TABLE IF EXISTS descriptions;
DROP TYPE IF EXISTS category_type;
DROP TYPE IF EXISTS stock_status;
-- +goose StatementEnd

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
END$$;
-- +goose StatementEnd

-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'feeds_status') THEN
CREATE TYPE feeds_status AS ENUM ('pending', 'accepted', 'blocked');
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
                "status" BOOLEAN NOT NULL DEFAULT TRUE,
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
                        "type" category_type NOT NULL,
                        "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cars (
                id SERIAL PRIMARY KEY,
                user_id BIGINT NOT NULL,
                stock_id BIGINT,
                brand_id BIGINT NOT NULL,
                model_id BIGINT NOT NULL,
                year BIGINT NOT NULL,
                mileage BIGINT NOT NULL,
                color VARCHAR(255) NOT NULL,
                engine_capacity DOUBLE PRECISION NOT NULL,
                engine_type VARCHAR(255) NOT NULL,
                body_id BIGINT NOT NULL,
                transmission VARCHAR(255) NOT NULL,
                drive_type VARCHAR(255) NOT NULL,
                vin VARCHAR(255),
                description TEXT,
                city_id BIGINT NOT NULL,
                name VARCHAR(255),
                mail VARCHAR(255),
                phone_number VARCHAR(50) NOT NULL,
                price BIGINT NOT NULL,
                is_comment BOOLEAN NOT NULL,
                is_exchange BOOLEAN NOT NULL,
                is_credit BOOLEAN NOT NULL,
                images JSONB,
                status feeds_status DEFAULT 'pending',
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                CONSTRAINT fk_user
                    FOREIGN KEY (user_id)
                        REFERENCES users (id)
                            ON DELETE SET NULL,
                CONSTRAINT fk_brand
                    FOREIGN KEY (brand_id)
                        REFERENCES brands (id)
                            ON DELETE CASCADE,
                CONSTRAINT fk_stock
                    FOREIGN KEY (stock_id)
                        REFERENCES stocks (id)
                            ON DELETE SET NULL,
                CONSTRAINT fk_model
                    FOREIGN KEY (model_id)
                        REFERENCES models (id)
                            ON DELETE CASCADE,
                CONSTRAINT fk_body
                    FOREIGN KEY (body_id)
                        REFERENCES body_types (id)
                            ON DELETE CASCADE,
                CONSTRAINT fk_city
                    FOREIGN KEY (city_id)
                        REFERENCES cities (id)
                            ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE cars OWNER TO autotm;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS trucks (
                id BIGSERIAL PRIMARY KEY,
                user_id BIGINT NOT NULL,
                stock_id BIGINT NULL,
                body_id BIGINT NOT NULL,
                brand_id BIGINT NOT NULL,
                model_id BIGINT NOT NULL,
                load_capacity DOUBLE PRECISION,
                price BIGINT NOT NULL,
                body_type VARCHAR(255),
                drive_type VARCHAR(255),
                transmission VARCHAR(255),
                engine_type VARCHAR(255),
                year BIGINT NOT NULL,
                seats BIGINT,
                cab_type VARCHAR(255),
                wheel_formula VARCHAR(255),
                chassis VARCHAR(255),
                cab_suspension VARCHAR(255),
                bus_type VARCHAR(255),
                suspension_type VARCHAR(255),
                brakes VARCHAR(255),
                axles BIGINT,
                engine_hours BIGINT,
                vehicle_type VARCHAR(255),
                engine_capacity DOUBLE PRECISION,
                forklift_type VARCHAR(255),
                lifting_capacity BIGINT,
                mileage BIGINT,
                excavator_type VARCHAR(255),
                bulldozer_type VARCHAR(255),
                color VARCHAR(255) NOT NULL,
                vin VARCHAR(255),
                description TEXT,
                city_id BIGINT NOT NULL,
                name VARCHAR(255),
                mail VARCHAR(255),
                phone_number VARCHAR(255) NOT NULL,
                is_comment BOOLEAN NOT NULL,
                is_exchange BOOLEAN NOT NULL,
                is_credit BOOLEAN NOT NULL,
                images JSONB,
                status feeds_status DEFAULT 'pending',
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                CONSTRAINT fk_user
                    FOREIGN KEY (user_id)
                        REFERENCES users (id)
                            ON DELETE SET NULL,
                CONSTRAINT fk_brand
                    FOREIGN KEY (brand_id)
                        REFERENCES brands (id)
                            ON DELETE CASCADE,
                CONSTRAINT fk_stock
                    FOREIGN KEY (stock_id)
                        REFERENCES stocks (id)
                            ON DELETE SET NULL,
                CONSTRAINT fk_model
                    FOREIGN KEY (model_id)
                        REFERENCES models (id)
                            ON DELETE CASCADE,
                CONSTRAINT fk_body
                    FOREIGN KEY (body_id)
                        REFERENCES body_types (id)
                            ON DELETE CASCADE,
                CONSTRAINT fk_city
                    FOREIGN KEY (city_id)
                        REFERENCES cities (id)
                            ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE trucks OWNER TO autotm;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS trucks;
DROP TABLE IF EXISTS cars;
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

DROP TYPE IF EXISTS feeds_status;
DROP TYPE IF EXISTS category_type;
DROP TYPE IF EXISTS stock_status;
-- +goose StatementEnd

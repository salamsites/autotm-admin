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
                              brand_id INTEGER,
                              created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                              updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                              CONSTRAINT brand_id_fk
                                  FOREIGN KEY (brand_id)
                                      REFERENCES brands(id)
                                      ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS roles (
                            "id" SERIAL PRIMARY KEY,
                            "name" CHARACTER VARYING(255),
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
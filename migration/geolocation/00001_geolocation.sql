-- +goose Up
/* Email table */

-- +goose StatementBegin
/* Add modified update function */
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.modified_at = now();
RETURN NEW;
END;
$$ language 'plpgsql';
-- +goose StatementEnd

CREATE TABLE geolocation (
                       id                          UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
                       ip                          TEXT NOT NULL UNIQUE,
                       country_code                TEXT NOT NULL,
                       country                     TEXT NOT NULL,
                       city                        TEXT NOT NULL,
                       latitude                    TEXT NOT NULL,
                       longitude                   TEXT NOT NULL,
                       created_at                     TIMESTAMP with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       modified_at                    TIMESTAMP with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX index_ip ON geolocation(ip);

CREATE TRIGGER update_geolocation_modified BEFORE UPDATE ON geolocation FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
-- +goose Down
DROP TRIGGER IF EXISTS update_geolocation_modified on geolocation;

DROP INDEX index_ip;
DROP TABLE email;

DROP FUNCTION IF EXISTS update_modified_column;

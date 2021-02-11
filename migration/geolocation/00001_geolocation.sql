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
                       id                          UUID PRIMARY KEY NOT NULL,
                       ip                          TEXT NOT NULL UNIQUE DEFAULT '',
                       country_code                TEXT NOT NULL DEFAULT '',
                       country                     TEXT NOT NULL DEFAULT '',
                       city                        TEXT NOT NULL DEFAULT '',
                       latitude                    TEXT NOT NULL DEFAULT '',
                       longitude                   TEXT NOT NULL DEFAULT '',
                       created_at                  TIMESTAMP with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       modified_at                  TIMESTAMP with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX index_ip ON geolocation(ip);

CREATE TRIGGER update_geolocation_modified BEFORE UPDATE ON geolocation FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
-- +goose Down
DROP TRIGGER IF EXISTS update_geolocation_modified on geolocation;

DROP INDEX index_ip;
DROP TABLE email;

DROP FUNCTION IF EXISTS update_modified_column;

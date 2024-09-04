-- +goose Up
-- +goose StatementBegin
CREATE TABLE event.main_table(
    eid String, -- Event ID -- TODO: NOT NULL
    provider_id String, -- Provider ID -- TODO: NOT NULL
    schema_id String, -- Event Schema ID, if Any
    type String, -- Type of event, if any
    data String, -- Event as bytes -- TODO: NOT NULL
    -- TODO, TIMESTAMP
)ENGINE = MergeTree() -- Replacing MergeTree
ORDER BY (type, provider_id);
-- +goose StatementEnd


-- TODO "github.com/golang-migrate/migrate/v4/database/clickhouse" 


-- +goose Down
-- +goose StatementBegin
DROP TABLE  event.main_table;
-- +goose StatementEnd

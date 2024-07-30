-- +goose Up
-- +goose StatementBegin
CREATE TABLE event.main_table(
    eid String,
    provider_id String,
    schema_id String,
    type String,
    data String,
)ENGINE = MergeTree()
ORDER BY (type);
-- +goose StatementEnd


-- TODO "github.com/golang-migrate/migrate/v4/database/clickhouse" 


-- +goose Down
-- +goose StatementBegin
DROP TABLE  event.main_table;
-- +goose StatementEnd

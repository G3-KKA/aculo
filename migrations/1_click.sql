-- +goose up 
-- +goose StatementBegin
 CREATE DATABASE event
-- +goose StatementEnd


-- TODO "github.com/golang-migrate/migrate/v4/database/clickhouse"


-- +goose down
-- +goose StatementBegin
 DROP DATABASE  event
-- +goose StatementEnd
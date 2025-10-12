-- +goose Up
-- +goose StatementBegin
alter table users add column birthday date not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users drop column birthday;
-- +goose StatementEnd

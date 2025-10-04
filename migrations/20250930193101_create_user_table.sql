-- +goose Up
-- +goose StatementBegin
create table if not exists users (
    id integer primary key autoincrement,
    first_name varchar(100) not null,
    last_name varchar(100) not null,
    email varchar(100) unique not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd


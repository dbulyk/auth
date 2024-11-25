-- +goose Up
-- +goose StatementBegin
alter table roles alter column role type varchar(255);
alter table users alter column name type varchar(255);
alter table users alter column email type varchar(255);
alter table users alter column password type varchar(255);
alter table users alter column tag type varchar(50);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

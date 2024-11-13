-- +goose Up
-- +goose StatementBegin
create table roles (
    id serial primary key,
    role text not null
);

insert into roles (role) values ('ROLE_ADMIN'), ('ROLE_USER');

create table users (
    id serial primary key,
    name text not null,
    email text not null,
    password text not null,
    tag text not null,
    role integer references roles not null,
    created_at timestamp default now(),
    updated_at timestamp default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
drop table roles;
-- +goose StatementEnd

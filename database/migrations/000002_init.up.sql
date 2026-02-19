create table if not exists users (
    id serial primary key,
    name varchar(255) not null,
    email varchar(255),
    created_at timestamp default now(),
    updated_at timestamp default now()
);

create table if not exists tasks (
    id serial primary key,
    title varchar(255) not null,
    done boolean default false
)

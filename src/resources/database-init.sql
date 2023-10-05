DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS goods;

create table users(
    id       serial primary key   not null,
    email    text unique          not null,
    enabled  boolean default true not null,
    password text                 not null,
    username text
);

create table goods(
    id bigserial primary key not null,
    seller_id int unique not null,
    price numeric not null,
    image_urls text[] default null 
);


create table users
(
    u_id         uuid
        constraint users_pk
            primary key,
    created_at   timestamp,
    updated_at   timestamp,
    user_version int default 1,
    email        varchar,
    user_locale  varchar,
    user_name    varchar
);
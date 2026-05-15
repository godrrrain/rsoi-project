-- file: 10-create-user.sql
CREATE ROLE program WITH PASSWORD 'test';
ALTER ROLE program WITH LOGIN;

CREATE DATABASE reservations;
GRANT ALL PRIVILEGES ON DATABASE reservations TO program;

CREATE DATABASE libraries;
GRANT ALL PRIVILEGES ON DATABASE libraries TO program;

CREATE DATABASE ratings;
GRANT ALL PRIVILEGES ON DATABASE ratings TO program;

CREATE DATABASE idp;
GRANT ALL PRIVILEGES ON DATABASE idp TO program;


\c reservations;

CREATE TABLE reservation
(
    id              SERIAL PRIMARY KEY,
    reservation_uid uuid UNIQUE NOT NULL,
    username        VARCHAR(80) NOT NULL,
    book_uid        uuid        NOT NULL,
    library_uid     uuid        NOT NULL,
    status          VARCHAR(20) NOT NULL
        CHECK (status IN ('RENTED', 'RETURNED', 'EXPIRED')),
    start_date      TIMESTAMP   NOT NULL,
    till_date       TIMESTAMP   NOT NULL
);

GRANT ALL ON ALL TABLES IN SCHEMA public TO program;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO program;

\c libraries;

CREATE TABLE library
(
    id          SERIAL PRIMARY KEY,
    library_uid uuid UNIQUE  NOT NULL,
    name        VARCHAR(80)  NOT NULL,
    city        VARCHAR(255) NOT NULL,
    address     VARCHAR(255) NOT NULL
);

CREATE TABLE books
(
    id        SERIAL PRIMARY KEY,
    book_uid  uuid UNIQUE  NOT NULL,
    name      VARCHAR(255) NOT NULL,
    author    VARCHAR(255),
    genre     VARCHAR(255),
    condition VARCHAR(20) DEFAULT 'EXCELLENT'
        CHECK (condition IN ('EXCELLENT', 'GOOD', 'BAD'))
);

CREATE TABLE library_books
(
    book_id         INT REFERENCES books (id),
    library_id      INT REFERENCES library (id),
    available_count INT NOT NULL
);

GRANT ALL ON ALL TABLES IN SCHEMA public TO program;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO program;


INSERT INTO library (library_uid, name, city, address) VALUES ('83575e12-7ce0-48ee-9931-51919ff3c9ee', 'Библиотека имени 7 Непьющих', 'Москва', '2-я Бауманская ул., д.5, стр.1');
INSERT INTO library (library_uid, name, city, address) VALUES ('d6a8b7f2-1d0e-4c3b-9f54-1ebc8a7f2b33', 'Московская районная библиотека', 'Москва', 'ул. Академика Сахарова, д. 12');
INSERT INTO library (library_uid, name, city, address) VALUES ('a8f5c8d4-3c72-4f2e-9a44-5f2c0d3a6e11', 'Центральная библиотека Города', 'Санкт-Петербург', 'Невский проспект, д. 28');
INSERT INTO library (library_uid, name, city, address) VALUES ('b4d3f7e1-99de-4c56-8bf8-2f0b9d2c6a22', 'Детская библиотека Солнечная', 'Челябинск', 'ул. Кирова, д. 15');

INSERT INTO books (book_uid, name, author, genre, condition) VALUES ('f7cdc58f-2caf-4b15-9727-f89dcc629b27', 'Краткий курс C++ в 7 томах', 'Бьерн Страуструп', 'Научная фантастика', 'EXCELLENT');
INSERT INTO books (book_uid, name, author, genre, condition) VALUES ('c31d2f41-2f5a-4e5d-8f81-1c6b2a7e9d33', 'Война и мир', 'Лев Толстой', 'Историческая проза', 'GOOD');
INSERT INTO books (book_uid, name, author, genre, condition) VALUES ('d52e3a65-66a7-4a89-bc73-2f9f0f8e5b44', 'Мастер и Маргарита', 'Михаил Булгаков', 'Магический реализм', 'EXCELLENT');
INSERT INTO books (book_uid, name, author, genre, condition) VALUES ('e64f7b28-98d2-4f90-a7ed-8e0c6d1a0f55', 'Приключения Шерлока Холмса', 'Артур Конан Дойл', 'Детектив', 'GOOD');
INSERT INTO books (book_uid, name, author, genre, condition) VALUES ('f75a8c93-7b54-4c21-a23f-9d3b4e7a1c66', 'Азбука путешествий', 'Анна Петрова', 'Путешествия', 'BAD');
INSERT INTO books (book_uid, name, author, genre, condition) VALUES ('c46b8b3f-d102-41f5-875f-b89ac0b5a24a', 'Москва в книгах и легендах', 'Екатерина Орлова', 'Городской путеводитель', 'GOOD');
INSERT INTO books (book_uid, name, author, genre, condition) VALUES ('b0fdb3b9-cee4-4924-9217-c2c19d1d1fd2', 'Истории Московских улиц', 'Иван Смирнов', 'История', 'EXCELLENT');

INSERT INTO library_books (book_id, library_id, available_count)  VALUES (1, 1, 1);
INSERT INTO library_books (book_id, library_id, available_count) VALUES (6, 2, 2);
INSERT INTO library_books (book_id, library_id, available_count) VALUES (7, 2, 1);
INSERT INTO library_books (book_id, library_id, available_count) VALUES (2, 3, 2);
INSERT INTO library_books (book_id, library_id, available_count) VALUES (3, 3, 1);
INSERT INTO library_books (book_id, library_id, available_count) VALUES (4, 4, 3);
INSERT INTO library_books (book_id, library_id, available_count) VALUES (5, 4, 2);

\c ratings;

CREATE TABLE rating
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(80) NOT NULL,
    stars    INT         NOT NULL
        CHECK (stars BETWEEN 0 AND 100)
);

GRANT ALL ON ALL TABLES IN SCHEMA public TO program;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO program;

INSERT INTO rating (username, stars) VALUES ('admin', 20);
INSERT INTO rating (username, stars) VALUES ('user', 20);

\c idp;

CREATE TABLE idp_users (
    id SERIAL PRIMARY KEY,
    user_uid UUID NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    role VARCHAR(50) DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE auth_codes (
    id SERIAL PRIMARY KEY,
    code VARCHAR(255) NOT NULL UNIQUE,
    user_uid UUID NOT NULL,
    client_id VARCHAR(255) NOT NULL,
    redirect_uri VARCHAR(1024) NOT NULL,
    scope VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_uid) REFERENCES idp_users(user_uid)
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR(511) NOT NULL UNIQUE,
    user_uid UUID NOT NULL,
    client_id VARCHAR(255) NOT NULL,
    scope VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_uid) REFERENCES idp_users(user_uid)
);

GRANT ALL ON ALL TABLES IN SCHEMA public TO program;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO program;

INSERT INTO idp_users (user_uid, username, email, password_hash, full_name, role) VALUES ('1ce9ed92-8548-4ed9-a18e-d96fb120e622', 'admin', 'admin@test.ru', '8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918', 'Admin', 'admin');
INSERT INTO idp_users (user_uid, username, email, password_hash, full_name, role) VALUES ('2b1f83a3-9f95-4e5d-8f0f-3baf58c2f864', 'user', 'user@user.ru', '04f8996da763b7a969b1028ee3007569eaf3a635486ddab211d512c85b9df8fb', 'User', 'user');
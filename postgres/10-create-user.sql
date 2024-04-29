-- file: 10-create-user.sql
CREATE ROLE program WITH PASSWORD 'test';
ALTER ROLE program WITH LOGIN;

CREATE DATABASE reservations;
GRANT ALL PRIVILEGES ON DATABASE reservations TO program;

CREATE DATABASE libraries;
GRANT ALL PRIVILEGES ON DATABASE libraries TO program;

CREATE DATABASE ratings;
GRANT ALL PRIVILEGES ON DATABASE ratings TO program;


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

-- UPDATE library_books SET available_count = 0 WHERE book_id = 3;
-- SELECT books.*, library_books.available_count from library_books, books, library 
-- 	where books.book_uid = 'c6cdb5f4-40c2-4658-b71d-66385e8707ee' and library.id = library_books.library_id 
-- 	and books.id = library_books.book_id;

INSERT INTO library VALUES (1, '83575e12-7ce0-48ee-9931-51919ff3c9ee', 'Библиотека имени 7 Непьющих', 'Москва', '2-я Бауманская ул., д.5, стр.1');
INSERT INTO books VALUES (1, 'f7cdc58f-2caf-4b15-9727-f89dcc629b27', 'Краткий курс C++ в 7 томах', 'Бьерн Страуструп', 'Научная фантастика', 'EXCELLENT');
INSERT INTO library_books VALUES (1, 1, 1);

-- INSERT INTO books VALUES (2, 'b0a67f71-c27b-4c1b-8360-6b9033157c3e', 'Облачный GO', 'Мэтью Титмус', 'Научная фантастика', 'EXCELLENT');
-- INSERT INTO library_books VALUES (2, 1, 2);

-- INSERT INTO library VALUES (2, 'd31f6751-9421-48af-9667-e5ca97bd6295', 'Библиотека имени Шоколада', 'Москва', 'Центральная ул., д.2, стр.1');
-- INSERT INTO books VALUES (3, 'c6cdb5f4-40c2-4658-b71d-66385e8707ee', 'Совершенный код', 'Стив Макконнелл', 'Научная фантастика', 'EXCELLENT');
-- INSERT INTO library_books VALUES (3, 2, 1);

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

-- INSERT INTO rating VALUES (1, 'godrain', 20);
INSERT INTO rating VALUES (1, 'Test Max', 20);
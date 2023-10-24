CREATE TABLE people (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255),
    age INTEGER NOT NULL,
    gender VARCHAR(6) NOT NULL,
    nationality VARCHAR(255) NOT NULL
)
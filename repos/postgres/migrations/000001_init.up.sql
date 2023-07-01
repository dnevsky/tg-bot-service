CREATE TABLE services
(
    id SERIAL PRIMARY KEY,
    userid INT NOT NULL,
    nameofservice VARCHAR(255) NOT NULL,
    login VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);
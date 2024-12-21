CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    last_notified TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    phone_number  VARCHAR(20) UNIQUE NOT NULL
);

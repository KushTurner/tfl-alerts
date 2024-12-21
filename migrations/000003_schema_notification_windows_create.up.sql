CREATE TABLE notification_windows
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER REFERENCES users (id),
    train_id   INTEGER REFERENCES trains (id),
    start_time TIMESTAMP,
    end_time   TIMESTAMP,
    UNIQUE (user_id, train_id)
);

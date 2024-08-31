CREATE TABLE reservations (
    id SERIAL PRIMARY KEY,
    room_id VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL
);

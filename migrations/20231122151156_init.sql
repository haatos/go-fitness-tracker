-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user (
	id TEXT PRIMARY KEY,
	email TEXT,
	password_hash TEXT,
	created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS exercise (
    id TEXT PRIMARY KEY,
    name TEXT,
    user_id TEXT,
    FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE,
    UNIQUE(name, user_id)
);

CREATE TABLE IF NOT EXISTS workout (
    id TEXT PRIMARY KEY,
    name TEXT,
    user_id TEXT,
    FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS junction (
    id TEXT PRIMARY KEY,
    exercise_id TEXT,
    workout_id TEXT,
    user_id TEXT,
    set_count INTEGER,
    FOREIGN KEY(exercise_id) REFERENCES exercise(id),
    FOREIGN KEY(workout_id) REFERENCES workout(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS entry (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    junction_id TEXT,
    set_number INTEGER,
    weight INTEGER,
    reps INTEGER,
    time TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY(junction_id) REFERENCES junction(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE entry;
DROP TABLE junction;
DROP TABLE workout;
DROP TABLE exercise;
DROP TABLE user;
-- +goose StatementEnd

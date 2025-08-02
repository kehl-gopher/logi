ALTER TABLE users
DROP COLUMN first_name;

ALTER TABLE users
DROP COLUMN last_name;


CREATE TABLE profiles (
    id TEXT PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    full_name VARCHAR(100),
    profile_pic TEXT,
    updated_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
)
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    email text NOT NULL,
    password text NOT NULL,
    first_name text NOT NULL,
	last_name text NOT NULL,
	phone text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email);
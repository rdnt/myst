CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    date timestamp NOT NULL DEFAULT NOW()::timestamp
);

CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
    post_id INTEGER REFERENCES posts (id),
    author TEXT NOT NULL,
    content TEXT NOT NULL,
    date timestamp NOT NULL DEFAULT NOW()::timestamp
);

-- INSERT INTO posts VALUES
--     (default, default, 'title1', 'content1', default),
--     (default, default, 'title2', 'content2', default),
--     (default, default, 'title3', 'content3', default),
--     (default, default, 'title4', 'content4', default);

-- INSERT INTO comments VALUES
--     (default, default, 25, 'author1', 'content1', default)

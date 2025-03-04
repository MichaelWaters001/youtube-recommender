-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    google_id TEXT UNIQUE NOT NULL
);

-- Create creators table
CREATE TABLE creators (
    id SERIAL PRIMARY KEY,
    youtube_id TEXT UNIQUE NOT NULL,
    youtube_handle TEXT UNIQUE,
    name TEXT NOT NULL,
    description TEXT
);

-- Create tags table
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Create creator_tags table (many-to-many relationship between creators and tags)
CREATE TABLE creator_tags (
    id SERIAL PRIMARY KEY,
    creator_id INT NOT NULL REFERENCES creators(id) ON DELETE CASCADE,
    tag_id INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (creator_id, tag_id, user_id)
);

-- Create votes table (tracks upvotes/downvotes for creator_tags)
CREATE TABLE votes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    creator_tag_id INT NOT NULL REFERENCES creator_tags(id) ON DELETE CASCADE,
    vote_type INT NOT NULL CHECK (vote_type IN (-1, 1)),
    UNIQUE (user_id, creator_tag_id)
);

-- Indexes for performance optimization
CREATE INDEX idx_creators_youtube_id ON creators(youtube_id);
CREATE INDEX idx_tags_name ON tags(name);
CREATE INDEX idx_creator_tags_creator_id ON creator_tags(creator_id);
CREATE INDEX idx_creator_tags_tag_id ON creator_tags(tag_id);
CREATE INDEX idx_votes_user_creator_tag ON votes(user_id, creator_tag_id);
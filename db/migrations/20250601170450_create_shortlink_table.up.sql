CREATE TABLE shortlinks (
   id  UUID PRIMARY KEY NOT NULL,
   owner_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
   short_name VARCHAR(255) UNIQUE NOT NULL,
   original_url VARCHAR(255) NOT NULL,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_shortlinks_id ON shortlinks (id);
CREATE INDEX idx_shortlinks_owner_id ON shortlinks (owner_id);
CREATE INDEX idx_shortlinks_short_name ON shortlinks (short_name);
CREATE INDEX idx_shortlinks_created_at ON shortlinks (created_at);
CREATE INDEX idx_shortlinks_deleted_at ON shortlinks (deleted_at);
CREATE TABLE users (
   user_id UUID PRIMARY KEY NOT NULL,
   email VARCHAR(255) UNIQUE NOT NULL,
   password VARCHAR(255) NOT NULL,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_users_user_id ON users (user_id);
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_created_at ON users (created_at);
CREATE INDEX idx_users_deleted_at ON users (deleted_at);
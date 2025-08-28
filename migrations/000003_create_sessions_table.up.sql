CREATE TABLE
    sessions (
        token VARCHAR(255),
        user_id UUID NOT NULL,
        type VARCHAR(255),
        os VARCHAR(255),
        expired_at TIMESTAMP,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (id),
        PRIMARY KEY (token, user_id)
    );
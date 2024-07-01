CREATE TABLE IF NOT EXISTS identities (
                                          id SERIAL PRIMARY KEY,
                                          email VARCHAR(100) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );


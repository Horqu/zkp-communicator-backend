-- User Schema: ZKP registration, public key
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    public_key TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Session Schema: stores active user sessions
CREATE TABLE sessions (
    session_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Contacts Schema: manage user contacts
CREATE TABLE contacts (
    contact_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id),
    contact_user_id INT NOT NULL REFERENCES users(user_id),
    status VARCHAR(20) NOT NULL, -- e.g.: 'pending', 'accepted', 'blocked'
    created_at TIMESTAMP DEFAULT NOW()
);

-- Private Messages Schema: end-to-end encrypted communication
CREATE TABLE messages (
    message_id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL REFERENCES users(user_id),
    recipient_id INT NOT NULL REFERENCES users(user_id),
    content TEXT NOT NULL,  -- encrypted message content
    created_at TIMESTAMP DEFAULT NOW()
);

-- Groups Schema: group chat data, age verification/ZKP
CREATE TABLE groups (
    group_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_by INT NOT NULL REFERENCES users(user_id),
    created_at TIMESTAMP DEFAULT NOW(),
    requires_age_verification BOOLEAN DEFAULT FALSE -- e.g.: 18+ group
);

-- Group Members Schema: users assigned to a group
CREATE TABLE group_members (
    group_id INT NOT NULL REFERENCES groups(group_id),
    user_id INT NOT NULL REFERENCES users(user_id),
    joined_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (group_id, user_id)
);

-- (Optional) Schema for storing age verification proofs
CREATE TABLE age_proofs (
    proof_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id),
    proof_data TEXT NOT NULL, -- encrypted proof data
    created_at TIMESTAMP DEFAULT NOW()
);
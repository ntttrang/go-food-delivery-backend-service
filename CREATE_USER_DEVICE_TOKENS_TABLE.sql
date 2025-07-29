-- Create user_device_tokens table
-- This table replaces the old refresh_tokens table with enhanced device tracking

CREATE TABLE user_device_tokens (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    is_production BOOLEAN DEFAULT TRUE,
    os VARCHAR(50) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_user_device_tokens_user_id (user_id),
    INDEX idx_user_device_tokens_token (token),
    INDEX idx_user_device_tokens_expires_at (expires_at),
    INDEX idx_user_device_tokens_is_revoked (is_revoked),
    INDEX idx_user_device_tokens_is_production (is_production),
    INDEX idx_user_device_tokens_os (os),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Sample data for testing (optional)
-- INSERT INTO user_device_tokens (
--     id, 
--     user_id, 
--     token, 
--     expires_at, 
--     is_revoked,
--     is_production,
--     os
-- ) VALUES (
--     UUID(),
--     'your-test-user-uuid-here',
--     'Kv7QzJj2fJ8vP2R5xN9mL3wY6tE1sA4hG8nB7cX0zM9',
--     DATE_ADD(NOW(), INTERVAL 30 DAY),
--     FALSE,
--     TRUE,
--     'iOS'
-- );

-- Cleanup query for expired/revoked tokens
-- DELETE FROM user_device_tokens 
-- WHERE expires_at < NOW() OR is_revoked = TRUE;

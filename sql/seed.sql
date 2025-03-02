-- Sample data for testing

-- Insert sample users
INSERT INTO users (id, name, email, created_at, updated_at)
VALUES 
    ('00000000-0000-0000-0000-000000000001', 'John Doe', 'john@example.com', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000002', 'Jane Smith', 'jane@example.com', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000003', 'Bob Johnson', 'bob@example.com', NOW(), NOW());

-- Insert sample posts
INSERT INTO posts (id, user_id, title, content, created_at, updated_at)
VALUES
    ('00000000-0000-0000-0000-000000000101', '00000000-0000-0000-0000-000000000001', 'First Post', 'This is the content of the first post', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000102', '00000000-0000-0000-0000-000000000001', 'Second Post', 'This is the content of the second post', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000103', '00000000-0000-0000-0000-000000000002', 'Hello World', 'Introduction to the world of blogging', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000104', '00000000-0000-0000-0000-000000000003', 'My Experience', 'Sharing my personal experience with transactions', NOW(), NOW()); 

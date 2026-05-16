INSERT INTO customers (first_name, last_name, email)
VALUES
    ('John', 'Doe', 'john.doe@example.com'),
    ('Jane', 'Smith', 'jane.smith@example.com'),
    ('Alice', 'Johnson', 'alice.johnson@example.com'),
    ('Bob', 'Williams', 'bob.williams@example.com'),
    ('Carol', 'Brown', 'carol.brown@example.com')
ON CONFLICT (email) DO NOTHING;

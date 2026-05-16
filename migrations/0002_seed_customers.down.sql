DELETE FROM customers
WHERE email IN (
    'john.doe@example.com',
    'jane.smith@example.com',
    'alice.johnson@example.com',
    'bob.williams@example.com',
    'carol.brown@example.com'
);

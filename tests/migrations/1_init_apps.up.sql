INSERT INTO apps (id, name, secret)
VALUES (1, 'test-app', 'test-secret-key')
ON CONFLICT DO NOTHING;

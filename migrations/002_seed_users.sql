-- Seed users table with 10 sample rows
INSERT INTO
  users (name, email, created_at, updated_at)
VALUES
  ('John Doe', 'john.doe@example.com', NOW(), NOW()),
  (
    'Jane Smith',
    'jane.smith@example.com',
    NOW(),
    NOW()
  ),
  (
    'Bob Johnson',
    'bob.johnson@example.com',
    NOW(),
    NOW()
  ),
  (
    'Alice Williams',
    'alice.williams@example.com',
    NOW(),
    NOW()
  ),
  (
    'Charlie Brown',
    'charlie.brown@example.com',
    NOW(),
    NOW()
  ),
  (
    'Diana Ross',
    'diana.ross@example.com',
    NOW(),
    NOW()
  ),
  (
    'Edward Norton',
    'edward.norton@example.com',
    NOW(),
    NOW()
  ),
  (
    'Fiona Apple',
    'fiona.apple@example.com',
    NOW(),
    NOW()
  ),
  (
    'George Lucas',
    'george.lucas@example.com',
    NOW(),
    NOW()
  ),
  (
    'Hannah Montana',
    'hannah.montana@example.com',
    NOW(),
    NOW()
  ) ON CONFLICT (email) DO NOTHING;
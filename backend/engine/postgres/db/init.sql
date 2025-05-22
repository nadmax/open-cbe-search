SELECT 'CREATE DATABASE open_cbe_database'
WHERE NOT EXISTS (
  SELECT FROM pg_database WHERE datname = 'open_cbe_database'
)\gexec
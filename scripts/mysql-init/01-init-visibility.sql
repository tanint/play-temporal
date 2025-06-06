-- Create temporal_visibility database
CREATE DATABASE IF NOT EXISTS temporal_visibility;

-- Grant all privileges to temporal user for both databases
GRANT ALL PRIVILEGES ON temporal.* TO 'temporal'@'%';
GRANT ALL PRIVILEGES ON temporal_visibility.* TO 'temporal'@'%';

-- Apply changes
FLUSH PRIVILEGES;

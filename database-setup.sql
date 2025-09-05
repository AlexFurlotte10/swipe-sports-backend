-- Production Database Setup for Swipe Sports
-- Run these commands on your production MySQL server

-- 1. Create the database
CREATE DATABASE IF NOT EXISTS swipe_sports
  CHARACTER SET utf8mb4 
  COLLATE utf8mb4_unicode_ci;

-- 2. Create a dedicated user for the application
CREATE USER IF NOT EXISTS 'swipe_user'@'%' IDENTIFIED BY 'CHANGE_THIS_SECURE_PASSWORD';

-- 3. Grant necessary permissions
GRANT SELECT, INSERT, UPDATE, DELETE ON swipe_sports.* TO 'swipe_user'@'%';
GRANT CREATE, ALTER, INDEX ON swipe_sports.* TO 'swipe_user'@'%';

-- 4. Apply changes
FLUSH PRIVILEGES;

-- 5. Verify setup
USE swipe_sports;
SHOW GRANTS FOR 'swipe_user'@'%';

-- Note: Tables will be created automatically by your Go application
-- The users table will be created with this schema:
/*
CREATE TABLE users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  oauth_id VARCHAR(255) UNIQUE,
  oauth_provider VARCHAR(50),
  name VARCHAR(255) NOT NULL,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  age INT,
  email VARCHAR(255) UNIQUE,
  gender ENUM('male', 'female', 'other'),
  location VARCHAR(255),
  latitude DECIMAL(10, 8),
  longitude DECIMAL(11, 8),
  `rank` INT DEFAULT 1000,
  profile_pic_url VARCHAR(500),
  bio TEXT,
  sport_preferences JSON,
  skill_level VARCHAR(50),
  ntrp_rating DECIMAL(3, 1),
  play_style VARCHAR(100),
  preferred_timeslots VARCHAR(100),
  availability JSON,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_location (location),
  INDEX idx_gender (gender),
  INDEX idx_rank (`rank`),
  INDEX idx_oauth (oauth_id, oauth_provider),
  INDEX idx_age (age),
  INDEX idx_skill_level (skill_level)
);
*/

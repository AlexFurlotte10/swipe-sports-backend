-- Initialize Swipe Sports Database
-- This script creates the database schema and adds some sample data

USE swipe_sports;

-- Create tables (these will be created by the Go application, but we can add indexes here)
-- The actual table creation is handled in the Go code

-- Add some sample users for testing
INSERT INTO users (oauth_id, oauth_provider, name, email, gender, location, latitude, longitude, rank, bio, sport_preferences, skill_level, play_style, availability) VALUES
('google_123456789', 'google', 'John Doe', 'john@example.com', 'male', 'New York, NY', 40.7128, -74.0060, 1200, 'Basketball enthusiast looking for pickup games', '{"basketball": true, "soccer": false, "tennis": true}', 'intermediate', 'competitive', '{"monday": ["18:00", "20:00"], "wednesday": ["18:00", "20:00"], "saturday": ["10:00", "14:00"]}'),
('google_987654321', 'google', 'Jane Smith', 'jane@example.com', 'female', 'Los Angeles, CA', 34.0522, -118.2437, 1100, 'Soccer player seeking team for weekend matches', '{"soccer": true, "basketball": false, "volleyball": true}', 'advanced', 'team_player', '{"saturday": ["09:00", "12:00"], "sunday": ["14:00", "17:00"]}'),
('google_555666777', 'google', 'Mike Johnson', 'mike@example.com', 'male', 'Chicago, IL', 41.8781, -87.6298, 1300, 'Tennis player looking for doubles partners', '{"tennis": true, "basketball": true, "soccer": false}', 'expert', 'strategic', '{"tuesday": ["19:00", "21:00"], "thursday": ["19:00", "21:00"], "sunday": ["08:00", "11:00"]}'),
('google_111222333', 'google', 'Sarah Wilson', 'sarah@example.com', 'female', 'Miami, FL', 25.7617, -80.1918, 1050, 'Volleyball enthusiast seeking beach volleyball games', '{"volleyball": true, "tennis": false, "soccer": true}', 'intermediate', 'casual', '{"friday": ["17:00", "19:00"], "saturday": ["16:00", "18:00"]}'),
('google_444555666', 'google', 'David Brown', 'david@example.com', 'male', 'Seattle, WA', 47.6062, -122.3321, 1150, 'Basketball player looking for indoor court games', '{"basketball": true, "soccer": true, "tennis": false}', 'advanced', 'aggressive', '{"monday": ["20:00", "22:00"], "wednesday": ["20:00", "22:00"], "friday": ["19:00", "21:00"]}');

-- Add some sample swipes (for testing match creation)
INSERT INTO swipes (swiper_id, swipee_id, direction) VALUES
(1, 2, 'right'),
(2, 1, 'right'),
(1, 3, 'left'),
(3, 1, 'right'),
(2, 3, 'right'),
(3, 2, 'right'),
(1, 4, 'right'),
(4, 1, 'left'),
(2, 4, 'right'),
(4, 2, 'right');

-- Create matches based on mutual right swipes
INSERT INTO matches (user1_id, user2_id) VALUES
(1, 2),  -- John and Jane matched
(2, 3),  -- Jane and Mike matched
(3, 2);  -- Mike and Jane matched (duplicate, will be handled by unique constraint)

-- Add some sample messages
INSERT INTO messages (match_id, sender_id, content, message_type) VALUES
(1, 1, 'Hey Jane! Would you be interested in playing soccer this weekend?', 'text'),
(1, 2, 'Absolutely! I love soccer. What time works for you?', 'text'),
(1, 1, 'How about Saturday at 10 AM?', 'text'),
(1, 2, 'Perfect! I\'ll bring my gear.', 'text'),
(2, 2, 'Hi Mike! I saw you play tennis. Interested in a doubles match?', 'text'),
(2, 3, 'Definitely! I\'m always looking for good partners.', 'text'); 
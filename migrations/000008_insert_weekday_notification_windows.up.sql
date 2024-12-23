ALTER TABLE IF EXISTS notification_windows ADD COLUMN weekday INTEGER CHECK ( weekday between 0 AND 6);

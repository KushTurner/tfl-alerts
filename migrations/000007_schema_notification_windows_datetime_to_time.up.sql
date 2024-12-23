ALTER TABLE IF EXISTS notification_windows
    ALTER COLUMN start_time TYPE TIME,
    ALTER COLUMN end_time TYPE TIME;
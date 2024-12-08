### I don't want to get a notification if I've already got it before

1. table with users - id, notification_windows, last_notified, number, lines
2. table with train_lines - name, status, description
3. get users which have times within notification window
4. if status from lines table is not equal to status from api or last_notified was not within window or null, send notification 
5. update last_notified for user
6. update status for train
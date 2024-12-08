### I don't want to get a notification if I've already got it before

1. table with users - id, notification_windows, last_notified, number, lines
2. table with train_lines - name, status, description
3. get users which have times within notification window and last_notified was not within window or null
4. if status from lines table is not equal to status from api, send notification 
5. update last_notified for user
6. update status for train
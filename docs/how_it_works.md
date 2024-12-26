### How it works

1. Poll TFL (every x) for current status of trains and update database
2. Get any trains that any user has subscribed to during the current time
3. For each train, look for users that have subscribed to trains during the current time
4. Check if current severity is different to previous severity
5. If different, notify each user
6. Update user last notified
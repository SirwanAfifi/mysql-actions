# MySQL Actions

This is a lot like GitHub Actions; Let's say you have a MySQL database with tables and you want something to happen when a new row is added, changed, or deleted, this is for you. For example, if you have a table of users and want a message to show up when a new user is inserted or deleted, you can set that up with a YAML file.

```yaml
name: Send Notification Example

on:
  insert:
    tables:
      - users
  delete:
    tables:
      - users

jobs:
  - name: Send Notification
    steps:
      - name: Run send_notification.sh script
        shell: bash
        # run: ./send_notification.sh
        run: |
          echo "Hello World"
```

This then prints "Hello World" when a new user is created or deleted.

## Running

Make sure you have a MySQL server running and you have set the these environment variables in a `.env` file:

```bash
MYSQL_USER=
MYSQL_PASSWORD=
MYSQL_HOST=
MYSQL_DB=
```

Then run the following command to start the application:

```bash
> go run ./cmd/mysql-actions/main.go
```

## How it works

The application basically uses MySQL triggers. Upon startup, it creates a table called `event_log` along with triggers for each table specified in the YAML file. When a row is either inserted or deleted from the tables, the corresponding trigger inserts a new record to the `event_log`. Subsequently, the application scans the `event_log` at regular intervals for new entries and executes the script delineated in the YAML file.

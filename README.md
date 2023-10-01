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

```bash
> go run ./cmd/mysql-actions/main.go
```

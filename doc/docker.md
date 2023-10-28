# Docker and PostgreSQL

@author: Deming Chen (陈德铭)

## Connect to PostgreSQL

While PostgreSQL (the docker container) is running, you are able to connect it. Run

```bash
docker exec -it [container name] psql -U [postgres_user]
```

to connect it. In our case, `[container name]` is `flagcamp-db` and `[postgres_user]` is just `postgres`. Then, you will see

```bash
psql (16.0 (Debian 16.0-1.pgdg120+1))
Type "help" for help.

postgres=#
```

To change database, use `\c`

```bash
postgres=# \c temp_db
You are now connected to database "temp_db" as user "postgres".
```

To list all tables, use `\d`

```bash
temp_db=# \d
              List of relations
 Schema |     Name     |   Type   |  Owner
--------+--------------+----------+----------
 public | items        | table    | postgres
 public | items_id_seq | sequence | postgres
 public | users        | table    | postgres
 public | users_id_seq | sequence | postgres
```


# Go Book-Rent-App

This is a REST API project. Main functionality is to rent books. This project is developed in Go.

## Project Features

- CRUD books
- Rent a book
- Return the book you rented
- Login & Register

## Technical Details

- Go
- MySQL
- Routing
- Middleware
- API Error Handling
- Database transactions
- JWT Authentication & Authorization

## Project Structure

I didn't use a complex folder structure because I didn't need it. I used very simple structure.

```bash
├── api
│   ├── book_handler.go
│   ├── middleware.go
│   ├── rent_handler.go
│   └── user_handler.go
├── database
│   └── db.go
├── go.mod
├── go.sum
├── helpers
│   └── helpers.go
├── main.go
├── store
│   ├── book_store.go
│   ├── rent_store.go
│   └── user_store.go
└── types
    └── types.go
```

## MySQL Tables

books:

```bash
+------------+----------+------+-----+---------+----------------+
| Field      | Type     | Null | Key | Default | Extra          |
+------------+----------+------+-----+---------+----------------+
| id         | int      | NO   | PRI | NULL    | auto_increment |
| name       | text     | NO   |     | NULL    |                |
| created_at | datetime | YES  |     | NULL    |                |
| quantity   | int      | YES  |     | 0       |                |
+------------+----------+------+-----+---------+----------------+
```

<br>
users:

```bash
+------------+-------------+------+-----+---------+-------+
| Field      | Type        | Null | Key | Default | Extra |
+------------+-------------+------+-----+---------+-------+
| id         | varchar(40) | NO   | PRI | NULL    |       |
| username   | text        | NO   |     | NULL    |       |
| password   | text        | NO   |     | NULL    |       |
| first_name | text        | NO   |     | NULL    |       |
| last_name  | text        | NO   |     | NULL    |       |
| created_at | datetime    | YES  |     | NULL    |       |
| role       | varchar(32) | YES  |     | user    |       |
+------------+-------------+------+-----+---------+-------+
```

<br>
book_rent_history:

```bash
+-----------------------+-------------+------+-----+---------+-------+
| Field                 | Type        | Null | Key | Default | Extra |
+-----------------------+-------------+------+-----+---------+-------+
| id                    | varchar(40) | NO   | PRI | NULL    |       |
| book_id               | int         | NO   | MUL | NULL    |       |
| user_id               | varchar(40) | NO   | MUL | NULL    |       |
| rent_start_time       | datetime    | YES  |     | NULL    |       |
| rent_return_time      | datetime    | YES  |     | NULL    |       |
| rent_duration_in_days | int         | NO   |     | NULL    |       |
+-----------------------+-------------+------+-----+---------+-------+
```

## How rent works?

1. Insert new record to the "book_rent_history" table
2. Decrease the quantity variable by one in the "books" table

## How return works?

1. Update the rent_end_time variable in "book_rent_history" table
2. Increase the quantity variable by one in the "books" table

## Future improvements

1. Pagination
2. Filtered lists (delayed returns, old returns, etc.)
3. Better validation solution and more meaningful error messages

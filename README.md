# Gator

Gator is a command-line RSS feed aggregator written in Go. It allows multiple local users to add and follow RSS feeds, continuously collect posts from those feeds, and browse collected posts directly from the terminal.

Gator uses PostgreSQL for persistent storage, Goose for database migrations, and SQLC for generating type-safe Go database access code.

## Requirements

Before installing Gator, make sure you have:

- Go 1.27 or later
- PostgreSQL 15 or later

PostgreSQL must be running and accessible from your local environment.

## Installation

Install Gator directly from GitHub using Go:

```bash
go install github.com/JB4779/gator@latest
```

Make sure your Go binary directory is included in your `PATH`. After installation, you should be able to run:

```bash
gator
```

You do not need to use `go run .` after installing Gator. `go run .` is primarily used when developing the application from its source code.

## Database Setup

Create a PostgreSQL database for Gator. For example:

```sql
CREATE DATABASE gator;
```

This project uses Goose to manage the database schema.

Install Goose if necessary:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Clone the repository if you need the migration files:

```bash
git clone https://github.com/JB4779/gator.git
cd gator
```

From the `sql/schema` directory, run the migrations using your PostgreSQL connection string:

```bash
goose postgres "postgres://username:password@localhost:5432/gator" up
```

## Configuration

Gator reads its local configuration from:

```text
~/.gatorconfig.json
```

Create the file with a PostgreSQL connection string:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

For example, a local development configuration might look like:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

The `current_user_name` value is updated by Gator when you register or log in as a user.

Do not commit a configuration file containing real database credentials to source control.

## Usage

After installation, commands are run using:

```bash
gator <command> [arguments]
```

### Register a user

```bash
gator register tom
```

This creates the user and makes that user current.

### Log in as an existing user

```bash
gator login tom
```

### List users

```bash
gator users
```

The currently selected user is marked with `(current)`.

### Add an RSS feed

```bash
gator addfeed "Hacker News" "https://news.ycombinator.com/rss"
```

The user who adds a feed automatically follows it.

### List feeds

```bash
gator feeds
```

### Follow a feed

```bash
gator follow "https://news.ycombinator.com/rss"
```

### List followed feeds

```bash
gator following
```

### Unfollow a feed

```bash
gator unfollow "https://news.ycombinator.com/rss"
```

### Aggregate posts

Run the RSS collector with a duration between requests:

```bash
gator agg 1m
```

Gator continuously rotates through the stored feeds, downloads their RSS data, and stores new posts in PostgreSQL.

The aggregator is intended to keep running. Stop it with `Ctrl+C`.

Avoid using unnecessarily short intervals when fetching third-party RSS feeds.

### Browse posts

Show the two most recent posts from feeds followed by the current user:

```bash
gator browse
```

Specify a different limit if desired:

```bash
gator browse 10
```

### Reset the database data

For development and testing:

```bash
gator reset
```

This deletes the users and, through the database relationships, their associated feed data.

**Warning:** `reset` is destructive and is intended for development use.

## Development

Clone the repository:

```bash
git clone https://github.com/JB4779/gator.git
cd gator
```

Run Gator directly from source:

```bash
go run . <command>
```

Check the Go code:

```bash
go vet ./...
```

After changing SQL queries, regenerate the database code:

```bash
sqlc generate
```

Database schema changes should be added as new Goose migrations in:

```text
sql/schema
```

SQLC queries are stored in:

```text
sql/queries
```

Generated database code is stored in:

```text
internal/database
```

## Technology

Gator is built with:

- Go
- PostgreSQL
- SQLC
- Goose
- RSS/XML over HTTP
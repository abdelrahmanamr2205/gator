# Gator - RSS Feed Aggregator CLI

Gator is a command-line interface (CLI) application written in Go that acts as an RSS feed aggregator. It allows you to manage users, add RSS feeds, follow/unfollow them, fetch posts continuously in the background, and browse the aggregated posts.

This is a guided project from [Boot.dev](https://boot.dev).

## Features

* **User Management:** Register and log in as different users.
* **Feed Management:** Add new RSS feeds to the database.
* **Subscriptions:** Follow and unfollow feeds.
* **Continuous Aggregation:** Run a long-running worker to continuously fetch and parse RSS feeds at a configurable interval.
* **Browse Posts:** View the latest posts from the feeds you follow.

## Prerequisites

* [Go](https://golang.org/dl/) (version 1.26.1 or higher)
* [PostgreSQL](https://www.postgresql.org/)
* *(Optional)* [Goose](https://github.com/pressly/goose) for database migrations.

## Configuration

Before running the application, you need to configure your database connection. Gator expects a configuration file named `.gatorconfig.json` in your user's home directory (e.g., `~/.gatorconfig.json`).

Create the file and add your PostgreSQL connection string:


```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```
*(Replace the `username`, `password`, and database name with your actual PostgreSQL credentials).*

## Installation & Setup

1. **Clone the repository** (or navigate to your project directory).
2. **Run database migrations**: Navigate to `sql/schema` and run Goose to create the tables.
   ```bash
   cd sql/schema
   goose postgres "postgres://username:password@localhost:5432/gator?sslmode=disable" up
   ```
3. **Build the CLI tool**:
   ```bash
   go build -o gator .
   ```

## Usage

Run the compiled `gator` executable followed by a command.

### User Commands
* `gator register <username>`: Register a new user and automatically log in.
* `gator login <username>`: Log in as an existing user.
* `gator users`: List all registered users. The current user will be marked with `(current)`.
* `gator reset`: Delete all users from the database (use with caution!).

### Feed Commands
* `gator addfeed <feed_name> <feed_url>`: Add a new RSS feed and automatically follow it.
* `gator feeds`: List all available feeds in the database.
* `gator follow <feed_url>`: Follow an existing feed.
* `gator unfollow <feed_url>`: Unfollow a feed.
* `gator following`: List all feeds the current logged-in user is following.

### Aggregation & Reading Commands
* `gator agg <time_between_requests>`: Start the continuous feed aggregator worker. It will fetch the oldest unfetched feed at the specified interval (e.g., `gator agg 1m` or `gator agg 30s`). Keep this running in a separate terminal.
* `gator browse [limit]`: View the latest posts from your followed feeds. The default limit is 2 posts, but you can specify a custom number (e.g., `gator browse 5`).

## Technologies Used
* **Go**: Core language.
* **PostgreSQL**: Relational database for storing users, feeds, feed follows, and posts.
* **sqlc**: Type-safe Go SQL query generation.
* **Goose**: Database schema migrations.

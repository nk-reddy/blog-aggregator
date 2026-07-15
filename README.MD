# Overview
Gator is a RSS feed aggregator built with Go. 

You can register a user, follow RSS feeds by URL, periodically collect posts, and browse posts from feeds.

# Requirements
- Go 
- PostgreSQL

# Setup
1. Database
Create a PostgreSQL database named gator.

`psql postgres`

`CREATE DATABASE gator;`

`\q`

2. Configuration

Gator uses a JSON config file in your home directory:

`~/.gatorconfig.json`

Create the file with the following:

`{ "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable"}`

where username and password are your PostgreSQL credentials.


# Installing CLI
Clone the repo:
`git clone https://github.com/nk-reddy/blog-aggregator.git`
`cd blog-aggregator`
Then, apply the DB migrations with Goose:
`cd sql/schema`
`goose postgres "postgres://username:password@localhost:5432/gator" up` 
Now, install the CLI:
`go install .`
If the binary is named blog-aggregator, use that or build explicitly:
`go build -o gator .`
Now, you can run:
`./gator`

# Usage
## Account Commands
Register a user: `./gator register <username>`
Log into a user: `./gator login <username>`
List users: `./gator users`

## Feed Commands
Add a RSS feed: `./gator addfeed "<feed-name>" "<feed-url>"`
Example: `./gator add "TechCrunch" "https://techcrunch.com/feed/"`

List available feeds: `./gator feeds`
Follow existing feed: `./gator follow "<feed-url>"`
List followed feeds: `./gator following`

## Feed Aggregation 
Run the feed aggregator: `./gator agg <time-between-requests>`
The time interval must be a Go duration string, i.e. `10s`, `1m`, or `1h`.
Example: `./gator agg 1m`

The aggregator runs continuously, saving new posts from the feeds to the database.
Least recently fetched feeds are downloaded first.
Press `Ctrl+C` to stop the aggregator. Run commands in a new terminal.

## Browse Posts
`./gator browse`
Displays the 2 most recent posts from feeds followed by the user. 
Provide an optional argument to display more posts, i.e. `./gator browse 10`

## Reset Database
`./gator reset`
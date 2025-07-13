# Gossip Bot

A bot that collects top news titles from various sources and creates a daily digest. The project consists of two main components:

1. **GitHub Bot (ghbot)**: Creates a GitHub issue with the daily news titles
2. **Markdown Bot (mdbot)**: Creates a markdown file in the `news` directory with the daily news titles

## Features

- Scheduled pipeline that runs daily at 00:00 UTC
- Collects top news titles from configurable sources (e.g., Hacker News, Reddit)
- Creates a GitHub issue with a digest of the collected news titles
- Creates markdown files with news titles in the `news` directory
- Configurable via environment variables

## Supported Sources

- **Hacker News**: Fetches top stories from Hacker News
- **Reddit**: Fetches top posts from specified subreddits

## Configuration

The bot is configured via environment variables:

### GitHub Configuration

- `GITHUB_TOKEN`: GitHub token with permission to create issues
- `GITHUB_OWNER`: Owner of the GitHub repository
- `GITHUB_REPO`: Name of the GitHub repository

### Sources Configuration

- `SOURCES`: Comma-separated list of source names (e.g., `HackerNews,RedditGo,RedditPython`)

For each source, the following environment variables are required:

- `SOURCE_{NAME}_TYPE`: Type of the source (e.g., `hackernews`, `reddit`)
- `SOURCE_{NAME}_URL`: Base URL of the source
- `SOURCE_{NAME}_LIMIT`: Maximum number of news items to fetch (default: 10)
- `SOURCE_{NAME}_SUBSOURCE`: Sub-source for sources like Reddit (e.g., subreddit name)

### Reddit App Configuration (Optional)

For better Reddit API access and higher rate limits, you can configure Reddit OAuth2 app credentials:

- `REDDIT_APP_ID`: Reddit app ID (client ID) from your Reddit app
- `REDDIT_APP_SECRET`: Reddit app secret (client secret) from your Reddit app

To create a Reddit app:
1. Go to https://www.reddit.com/prefs/apps
2. Click "Create App" or "Create Another App"
3. Choose "script" as the app type
4. Use any URL for redirect URI (not used for script apps)
5. Copy the app ID (under the app name) and secret

If these credentials are not provided, the bot will fall back to unauthenticated requests, which have lower rate limits.

## Example Configuration

```env
GITHUB_TOKEN=your_github_token
GITHUB_OWNER=your_github_username
GITHUB_REPO=your_github_repo
SOURCES=HackerNews,RedditGo,RedditPython
SOURCE_HackerNews_TYPE=hackernews
SOURCE_HackerNews_URL=https://hacker-news.firebaseio.com/v0
SOURCE_HackerNews_LIMIT=10
SOURCE_RedditGo_TYPE=reddit
SOURCE_RedditGo_URL=https://www.reddit.com
SOURCE_RedditGo_SUBSOURCE=golang
SOURCE_RedditGo_LIMIT=5
SOURCE_RedditPython_TYPE=reddit
SOURCE_RedditPython_URL=https://www.reddit.com
SOURCE_RedditPython_SUBSOURCE=python
SOURCE_RedditPython_LIMIT=5
REDDIT_APP_ID=your_reddit_app_id
REDDIT_APP_SECRET=your_reddit_app_secret
```

## Running Locally

1. Clone the repository
2. Create a `.env` file with the required environment variables
3. Build and run the bots:

```bash
# Build and run the GitHub Bot
go build -o ghbot ./cmd/ghbot
./ghbot

# Build and run the Markdown Bot
go build -o mdbot ./cmd/mdbot
./mdbot
```

## Visual Studio Code Integration

The project includes VSCode configuration files for easy development:

1. **Launch Configurations**:
   - `Run Current Package`: Runs the package of the currently open file
   - `Run GitHub Bot`: Runs the GitHub bot
   - `Run Markdown Bot`: Runs the Markdown bot
   - `Run Tests`: Runs tests for the current package

2. **Tasks**:
   - `Build GitHub Bot`: Builds the GitHub bot
   - `Build Markdown Bot`: Builds the Markdown bot
   - `Build All`: Builds both bots
   - `Run GitHub Bot`: Runs the GitHub bot
   - `Run Markdown Bot`: Runs the Markdown bot
   - `Run Tests`: Runs all tests

## GitHub Actions

The bot is configured to run daily at 00:00 UTC via GitHub Actions. You need to set up the following secrets in your GitHub repository:

- `GOSSIP_GITHUB_TOKEN`: GitHub token with permission to create issues

The workflow can also be triggered manually via the GitHub Actions UI.

## Requirements

- Go 1.22 or higher
- GitHub token with permission to create issues

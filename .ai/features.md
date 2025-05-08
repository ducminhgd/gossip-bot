# Features

## Basics

1. This project is hosted on Github.
2. A scheduled pipeline with run daily at 00:00 UTC to get the top latest and most read or most views or top rated news by sources. They will be configurable in environment variables. For example: get top 10 news from HackerNews, get 10 posts in a forum of Reddit.
3. The news should be: top rated, or top liked, or top discussed, or top viewed, or top commented.
4. The news titles will be collected and a Github issue is created each day, using markdown, with the titles linked to the news. For example:

    ```markdown
    # <YYYY-MM-DD>

    ## Hacker News

    1. [title](https://news.ycombinator.com/item?id=37137555)
    2. [title](https://news.ycombinator.com/item?id=37137555)
    3. [title](https://news.ycombinator.com/item?id=37137555)

    ## Reddit Go

    1. [title](https://www.reddit.com/r/Go/comments/1310111/golang_120_is_now_available/)
    2. [title](https://www.reddit.com/r/Go/comments/1310111/golang_120_is_now_available/)

    ## Reddit Python

    1. [title](https://www.reddit.com/r/Python/comments/1310111/python_312_is_now_available/)
    2. [title](https://www.reddit.com/r/Python/comments/1310111/python_312_is_now_available/)
    ```

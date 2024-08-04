# HackerJobs
A command line tool index job postings from [Hacker News](https://news.ycombinator.com/) using the [Hacker News API from Firebase](https://github.com/HackerNews/API) and [Bleve](http://blevesearch.com/).


[![asciicast](https://asciinema.org/a/KxpasyadU7dwdAkaYSLrhc0l3.svg)](https://asciinema.org/a/KxpasyadU7dwdAkaYSLrhc0l3)

## Building and Running

### With Go
* Build with `go build`
* Then run with `./HackerJobs`

### With Docker
* Build with `docker build -t hacker_jobs . `
* Then run with `docker run hacker_jobs`

### With pre-build Docker image from [ghcr.io](https://github.com/amscotti/HackerJobs/pkgs/container/hackerjobs)
* To download and run, use `docker run ghcr.io/amscotti/hackerjobs:main`

### Command Line Arguments
```
Usage of ./HackerJobs:
  -c int
        Count of posting to be return (default 100)
  -j int
        Job posting ID from HackerNews (default 32677265)
  -q string
        Text to search for in postings (default "+text:golang +text:remote")
```

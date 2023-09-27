File to be completed later

- `update.py`: cleans go modcache and go cache and gets updated packages (and, opt., updated test stuff too), then opt. runs: \
`go mod tidy`
  - `usage: python3 update.py [t] [n]`
  - `python3 update.py` runs: \
    ```
    go clean -cache
    go clean -modcache
    go get -u ./...
    go mod tidy
    ```
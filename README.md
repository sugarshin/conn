# conn

The wrapper for [confluence-go-api](https://github.com/Virtomize/confluence-go-api)

Currently, It's available is `subpage` utilities only.

## Installation

```sh
go get github.com/sugarshin/conn
```

## Usage

```go
package main

import (
	"log"
	conn "github.com/sugarshin/conn"
)

func main() {
	client, err := conn.New("<confluence_endpoint>", "<confluence_username>", "<confluence_token_or_password>")
	if err != nil {
		log.Fatal("failed")
	}
	c, err := client.CreateSubPageContent(parentPageID, content)
}
```

### `CreateSubPageContentWithLatest`

### `CreateSubPageContentWith`

### `CreateSubPageContent`

### `GetLatestChildPageContent`

### `GetChildPageContentByID`

### `GetChildPageContentWith`

Besides these, all confluence-go-api functions are available. ref: https://pkg.go.dev/github.com/virtomize/confluence-go-api#pkg-functions



# follow

[![Build Status](https://travis-ci.org/zemirco/follow.svg)](https://travis-ci.org/zemirco/follow)
[![GoDoc](https://godoc.org/github.com/zemirco/follow?status.svg)](https://godoc.org/github.com/zemirco/follow)

Go client CouchDB [_changes](http://docs.couchdb.org/en/latest/api/database/changes.html) API.

## Example

```go
package main

import "fmt"
import "github.com/zemirco/follow"

func main() {

  // set CouchDB url and database name
  follow.Url = "http://127.0.0.1:5984/"
  follow.Database = "_users"

  // set query parameters
  params := follow.QueryParameters{
    Limit: 10,
  }

  // get all changes at once
  changes, err := follow.Changes(params)
  if err != nil {
    panic(err)
  }
  fmt.Println(changes)

  // listen continuously for changes
  changes, errors := follow.Continuous(params)
  for {
    select {
      // something changed
      case change, ok := <-changes:
      fmt.Println(ok, change)
      // an error happenend
      case err := <-errors:
      panic(err)
      // stop after 5 seconds
      case <-time.After(5 * time.Second):
      fmt.Println("done")
    }
  }

}
```

## Test

`go test`

## License

MIT

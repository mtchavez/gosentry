# Go Sentry

Go package for [sentry](http://app.getsentry.com)

## Install

Run `go get -u github.com/mtchavez/gosentry`

## Usage

Set up a `RavenConfig` with your Project DSN

```go
package main

import (
    sentry "githubcom/mtchavez/gosentry"
)

var Raven *sentry.RavenConfig

func main() {
    var err
    myDSN := ""http://username:my-pw@app.getsentry.com/{project_id}"
    Raven, err = sentry.Setup(myDSN)
    if err != nil {
        println(err.Error())
    }
}
```

Capture a panic and sent error to sentry

```go
func myFunc() {
    defer func() {
        if r := recover(); r != nil {
            Raven.Message(r, "in package.myFunc()")
        }
    }()

}
```

## TODO

* Beef up tests
* Need to get backtrace to show up
* Add go examples
* Clean up code

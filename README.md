# Controller Go SDK
[![Build Status](https://ci.deis.io/buildStatus/icon?job=Deis/controller-sdk-go/master)](https://ci.deis.io/job/Deis/job/controller-sdk-go/job/master/)
[![Go Report Card](https://goreportcard.com/badge/github.com/deis/controller-sdk-go)](https://goreportcard.com/report/github.com/deis/controller-sdk-go)
[![codebeat badge](https://codebeat.co/badges/2fdee091-714d-4860-ab19-dba7587a3158)](https://codebeat.co/projects/github-com-deis-controller-sdk-go)
[![GoDoc](https://godoc.org/github.com/deis/controller-sdk-go?status.svg)](https://godoc.org/github.com/deis/controller-sdk-go)

This is the Go SDK for interacting with the [Deis Controller](https://github.com/deis/controller).

### Usage

```go
import deis "github.com/deis/controller-sdk-go"
import "github.com/deis/controller-sdk-go/apps"
```

Construct a deis client to interact with the controller API. Then, get the first 100 apps the user has access to.

```go
//                    Verify SSL, Controller URL, API Token
client, err := deis.New(true, "deis.test.io", "abc123")
if err != nil {
    log.Fatal(err)
}
apps, _, err := apps.List(client, 100)
if err != nil {
    log.Fatal(err)
}
```

### Authentication

```go
import deis "github.com/deis/controller-sdk-go"
import "github.com/deis/controller-sdk-go/auth"
```

If you don't already have a token for a user, you can retrieve one with a username and password.

```go
// Create a client with a blank token to pass to login.
client, err := deis.New(true, "deis.test.io", "")
if err != nil {
    log.Fatal(err)
}
token, err := auth.Login(client, "user", "password")
if err != nil {
    log.Fatal(err)
}
// Set the client to use the retrieved token
client.Token = token
```

For a complete usage guide to the SDK, see [full package documentation](https://godoc.org/github.com/deis/controller-sdk-go).

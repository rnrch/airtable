# airtable

[![PkgGoDev](https://pkg.go.dev/badge/github.com/rnrch/airtable)](https://pkg.go.dev/github.com/rnrch/airtable)
[![Go Report Card](https://goreportcard.com/badge/github.com/rnrch/airtable)](https://goreportcard.com/report/github.com/rnrch/airtable)

A simple airtable client.

## Usage

Go to your [account page](https://airtable.com/account) to get an api token and go to [API page](https://airtable.com/api) and select the database to get the base ID.

Then create the client.

```go
c := airtable.NewClient("your_api_token","your_base_ID")
```

Call the client's `ListRecords`, `GetRecord`, `CreateRecords`, `DeleteRecords` and `PatchRecords` methods to access the API.

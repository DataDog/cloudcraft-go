# cloudcraft-go

[![Go Documentation](https://godocs.io/github.com/DataDog/cloudcraft-go?status.svg)](https://godocs.io/github.com/DataDog/cloudcraft-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/DataDog/cloudcraft-go)](https://goreportcard.com/report/github.com/DataDog/cloudcraft-go)

Package `cloudcraft-go` is a simple and easy-to-use package for interacting with [Cloudcraft's developer API](https://developers.cloudcraft.co/).

## Installation

To install `cloudcraft-go`, run:

```console
go get github.com/DataDog/cloudcraft-go
```

## Usage

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/DataDog/cloudcraft-go"
)

func main() {
	key, ok := os.LookupEnv("CLOUDCRAFT_API_KEY")
	if !ok {
		log.Fatal("missing env var: CLOUDCRAFT_API_KEY")
	}

	// Create new Config to be initialize a Client.
	cfg := cloudcraft.NewConfig(key)

	// Create a new Client instance with the given Config.
	client, err := cloudcraft.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// List all blueprints in an account.
	blueprints, _, err := client.Blueprint.List(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Print the name of each blueprint.
	for _, blueprint := range blueprints {
		log.Println(blueprint.Name)
	}
}
```

For more examples and usage details, please [check the Go reference documentation](https://godocs.io/github.com/DataDog/cloudcraft-go).

## Contributing

Anyone can help make `cloudcraft-go` better. Check out [the contribution guidelines](CONTRIBUTING.md) for more information.

---

Released under the [Apache-2.0 License](LICENSE.md).

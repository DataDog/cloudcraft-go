# cloudcraft-go

[![Go Documentation](https://godocs.io/github.com/DataDog/cloudcraft-go?status.svg)](https://godocs.io/github.com/DataDog/cloudcraft-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/DataDog/cloudcraft-go)](https://goreportcard.com/report/github.com/DataDog/cloudcraft-go)

![Cloudcraft diagram](https://static.cloudcraft.co/sdk/cloudcraft-sdk-example-1.svg)

Visualize your cloud architecture with Cloudcraft by Datadog, [the best way to create smart AWS and Azure diagrams](https://www.cloudcraft.co/).

Cloudcraft supports both manual and programmatic diagramming, as well as automatic reverse engineering of existing cloud environments into
beautiful system architecture diagrams.

This `cloudcraft-go` package provides an easy-to-use native Go SDK for interacting with [the Cloudcraft API](https://developers.cloudcraft.co/).

Use case examples:
- Snapshot and visually compare your live AWS or Azure environment before and after a deployment, in your app or as part of your automated CI pipeline
- Download an inventory of all your cloud resources from a linked account as JSON
- Write a converter from a third party data format to Cloudcraft diagrams
- Backup, export & import your Cloudcraft data
- Programmatically create Cloudcraft diagrams

This SDK requires a [Cloudcraft API key](https://developers.cloudcraft.co/#authentication) to use. [A free trial of Cloudcraft Pro](https://www.cloudcraft.co/pricing) with API access is available.

## Installation

To install `cloudcraft-go`, run:

```console
go get github.com/DataDog/cloudcraft-go
```

## Go SDK Documentation

Usage details and more examples, please [see the Go reference documentation](https://godocs.io/github.com/DataDog/cloudcraft-go).

## Example Usage

In the below example the Cloudcraft API key is read from the `CLOUDCRAFT_API_KEY` environment variable. Alternatively, pass in the key to the configuration directly.

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

## Contributing

Anyone can help make `cloudcraft-go` better. Check out [the contribution guidelines](CONTRIBUTING.md) for more information.

---

Released under the [Apache-2.0 License](LICENSE.md).

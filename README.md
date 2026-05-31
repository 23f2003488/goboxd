# goboxd (Team One8Go)

A Go-based HTTP service that securely executes untrusted code inside an nsjail sandbox. Built for the SEEK Hackathon.

We chose the standard library `net/http` for our web framework. The routing requirements for this service are limited to a few static endpoints, meaning third-party frameworks like Gin or Echo would introduce unnecessary dependency bloat without providing tangible architectural benefits.

## Running the service

The environment is fully containerized. Do not use bare `go run` commands.

To build the runtime container:
make build

To run the server locally on port 8080:
make run

To run isolated unit tests and linters:
make test
make lint

## Documentation
* System design and lifecycle: docs/architecture.md
* Security closure logs: docs/ai/adrs.md
* Generative AI logs and Postmortem: docs/ai/
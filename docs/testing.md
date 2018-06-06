# Testing

```sh
# Run basic tests
go test ./... --cover

# Build a coverage profile
go test ./pkg/chat -coverprofile=coverage.out -covermode=count

# View the coverage profile
go tool cover -html=coverage.out
```

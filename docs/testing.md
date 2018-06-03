# Testing

```sh
# Run basic tests
go test ./... --cover
# ?   	github.com/rocheio/chatbot/cmd/chatbot	[no test files]
# ok  	github.com/rocheio/chatbot/pkg/chat	0.006s	coverage: 89.0% of statements
# ok  	github.com/rocheio/chatbot/pkg/tally	0.006s	coverage: 33.3% of statements

# Build a coverage profile
go test ./pkg/chat -coverprofile=coverage.out -covermode=count

# View the coverage profile
go tool cover -html=coverage.out
```

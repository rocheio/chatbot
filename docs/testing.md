# Testing

```sh
# Run basic tests
go test ./... --cover
# ?   	github.com/rocheio/chatbot/cmd/chatbot	[no test files]
# ?   	github.com/rocheio/chatbot/cmd/chatbot-example	[no test files]
# ok  	github.com/rocheio/chatbot/pkg/chat	0.006s	coverage: 91.4% of statements
# ok  	github.com/rocheio/chatbot/pkg/tally	0.006s	coverage: 31.2% of statements

# Build a coverage profile
go test ./pkg/chat -coverprofile=coverage.out -covermode=count

# View the coverage profile
go tool cover -html=coverage.out
```

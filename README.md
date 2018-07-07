# chatbot

A (work-in-progress, hopefully someday) interactive chat bot written in Go.

## Getting Started

```sh
# Download Hitchhiker's Guide to the Galaxy for source material
wget http://www.clearwhitelight.org/hitch/hhgttg.txt -P ./data

# Run unit tests
go test ./... --cover

# Build and run the interactive chatbot program
go build ./cmd/chatbot
./chatbot --input="Hello, how are you?"
./chatbot --interactive
```

## Downloading Corpuses

Corpuses give the chatbot content to learn speech from.

- [Collection of corpuses](http://freeconnection.blogspot.hu/2016/04/conversational-datasets-for-train.html)
- [Librarian Inquiries](https://academiccommons.columbia.edu/catalog/ac:176612)
- [Hitchhiker's Guide to the Galaxy](http://www.clearwhitelight.org/hitch/hhgttg.txt)

# chatbot

A (work-in-progress, hopefully someday) interactive chat bot written in Go.

## Getting Started

```sh
# Run unit tests
go test ./pkg/chat --cover

# Build the chatbot program
go build -o ./chatbot ./cmd/chatbot

# Download Hitchhiker's Guide to the Galaxy for source material
wget http://www.clearwhitelight.org/hitch/hhgttg.txt -P ./data

# Run the chatbot program
./chatbot
```

## Downloading Corpuses

Corpuses give the chatbot content to learn speech from.

- [Collection of corpuses](http://freeconnection.blogspot.hu/2016/04/conversational-datasets-for-train.html)
- [Librarian Inquiries](https://academiccommons.columbia.edu/catalog/ac:176612)
- [Hitchhiker's Guide to the Galaxy](http://www.clearwhitelight.org/hitch/hhgttg.txt)

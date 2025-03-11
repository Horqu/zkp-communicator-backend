# Client for ZKP Communicator

Simple chat client built using the Fyne framework in Go. This application allows users to connect to a chat server and exchange messages in real-time.

## Project Structure

```
client
├── cmd
│   └── main.go          # Entry point of the application
├── internal
│   ├── ui
│   │   └── ui.go       # User interface components
│   └── logic
│       └── logic.go    # Core application logic
├── go.mod               # Module definition and dependencies
└── README.md            # Project documentation
```

## Setup Instructions

1. **Install dependencies:**
   Ensure you have Go installed on your machine. Run the following command to download the necessary dependencies:
   ```
   go mod tidy
   ```

2. **Run the application:**
   To start the application, execute:
   ```
   go run cmd/main.go
   ```
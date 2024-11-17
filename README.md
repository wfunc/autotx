
# AUTOTX

`AUTOTX` is an automated task execution framework built with Go. It is designed to perform tasks like logging into web pages, fetching data, and executing browser interactions using `chromedp`. The project is modular and extensible, allowing you to define and manage tasks efficiently.

> 🌐 **Available in other languages:**  
> [🇨🇳 简体中文](./README_CN.md)

## Directory Structure

```
AUTOTX
├── build/server/         # Compiled executables
│   └── autotx
├── runner/               # Task runner module
│   ├── runner_test.go    # Unit tests for runner
│   └── runner.go         # Runner implementation
├── task/                 # Task definitions and logic
│   ├── base.go           # BaseTask for shared functionality
│   ├── example.go        # Example task implementation
│   ├── items.go          # Additional task logic
│   ├── login.go          # Login task logic
│   ├── sign_in.go        # Sign-in task logic
│   ├── task.go           # Task interface and utilities
├── .gitignore            # Git ignore file
├── autotx                # Build output
├── build-go.sh           # Go build script
├── DockerfileGo          # Docker build configuration
├── go.mod                # Go module configuration
├── go.sum                # Go module dependencies
├── README_CN.md          # Chinese documentation
├── README.md             # English documentation (default)
├── run.sh                # Script to run the project
└── service.go            # Service entry point
```

## Features

- **Modular Design**: 
  - Centralized `BaseTask` for shared task properties and methods.
  - Extensible task definitions (e.g., `LoginTask`, `ExampleTask`).
- **Browser Automation**: 
  - Uses `chromedp` for headless Chrome interactions.
- **Task Runner**:
  - Manages task execution lifecycle (start, stop, loop).

## Usage

### 1. Clone the Repository
```bash
git clone https://github.com/your-repo/autotx.git
cd autotx
```

### 2. Install Dependencies
Ensure you have Go 1.18+ installed.
```bash
go mod tidy
```

### 3. Run the Project

#### Locally
```bash
go run service.go
```

#### Using the Run Script
```bash
./run.sh
```

#### Using Docker Build
```bash
docker build -t autotx -f DockerfileGo .
```

#### Using Docker run
```bash
docker run -it \
  -e HEADLESS=1 \
  -e Verbose=1 \
  -e CodeURL=https://example.com \
  autotx
```

## Testing

Run unit tests:
```bash
go test ./runner/ ./task/
```

## Adding New Tasks

1. Create a new file in the `task/` directory, e.g., `my_task.go`, and implement the `Task` interface.
2. Use the `BaseTask` to inherit shared logic.
3. Register the new task in your runner or service logic.

---

## Contributing

Contributions are welcome! Please submit a pull request with your changes.

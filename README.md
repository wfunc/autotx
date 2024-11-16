
# AUTOTX

`AUTOTX` is an automated task execution framework built with Go, designed to perform webpage operations using `chromedp`. The project is modular in design, featuring a task scheduler (`runner`) and task execution logic (`task`), making it easy to extend and manage.

> 🌐 **Available in other languages:**  
> [🇨🇳 简体中文](./README_CN.md)

## Introduction

This is the English version of the documentation. If you prefer reading in Chinese, please click the link above to access the Simplified Chinese version.

## Directory Structure

```
AUTOTX
├── build/server/         # Compiled executable files
│   └── autotx
├── runner/               # Task scheduling module
│   ├── runner_test.go    # Unit tests for the scheduler
│   └── runner.go         # Scheduler implementation
├── task/                 # Task module
│   ├── task_test.go      # Unit tests for tasks
│   └── task.go           # Task logic implementation
├── .gitignore            # Git ignore file
├── autotx                # Executable file (build output)
├── build-go.sh           # Go build script
├── DockerfileGo          # Docker build file
├── go.mod                # Go module configuration file
├── go.sum                # Go module dependency file
├── run.sh                # Run script
└── service.go            # Service entry point
```

## Features

- **Task Module**:
  - Each task implements the `Task` interface with `Run` and `Stop` methods.
  - The example task `ExampleTask` demonstrates webpage operations using `chromedp`.

- **Scheduler Module**:
  - Manages the addition, starting, stopping, and looping of tasks.
  - Supports high-concurrency task execution, ensuring graceful shutdown with `sync.WaitGroup`.

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

#### Run Locally
```bash
go run service.go
```

#### Using the Run Script
```bash
./run.sh
```

#### Using the Build Script
```bash
./build-go.sh
./build/server/autotx
```

### 4. Use Docker
Build the Docker image:
```bash
docker build -t autotx -f DockerfileGo .
```

Run the container:
```bash
docker run -it autotx
```

## Testing

Run unit tests:
```bash
go test ./runner/ ./task/
```

## Extending Features

### Add a New Task
Create a new file in the `task/` directory, e.g., `my_task.go`, and implement the `Task` interface:
```go
type MyTask struct{}

func (t *MyTask) Run(ctx context.Context) error {
    // Implement task logic
    return nil
}

func (t *MyTask) Stop() error {
    // Implement stop logic
    return nil
}
```

Add the task in `service.go`:
```go
runner.AddTask(&MyTask{})
```

### Modify Execution Logic
You can customize the scheduling logic in `runner/runner.go`, such as adding priorities or dynamically loading tasks.

## Contributing
Contributions are welcome! Please submit a Pull Request with your changes.

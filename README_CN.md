
# AUTOTX

`AUTOTX` 是一个基于 Go 语言的自动化任务执行框架，支持通过 `chromedp` 操作网页任务。项目具有模块化设计，包含任务调度器 (`runner`) 和任务执行逻辑 (`task`)，方便扩展和管理。

## 目录结构

```
AUTOTX
├── build/server/         # 构建后的可执行文件
│   └── autotx
├── runner/               # 任务调度模块
│   ├── runner_test.go    # 调度模块单元测试
│   └── runner.go         # 调度器实现
├── task/                 # 任务模块
│   ├── task_test.go      # 任务模块单元测试
│   └── task.go           # 任务逻辑实现
├── .gitignore            # Git 忽略文件
├── autotx                # 可执行文件（构建产物）
├── build-go.sh           # Go 构建脚本
├── DockerfileGo          # Docker 构建文件
├── go.mod                # Go 模块配置文件
├── go.sum                # Go 模块依赖文件
├── run.sh                # 运行脚本
└── service.go            # 服务入口
```

## 功能说明

- **任务模块 (task)**:
  - 每个任务实现 `Task` 接口，包含 `Run` 和 `Stop` 方法。
  - 示例任务 `ExampleTask` 使用 `chromedp` 执行网页操作。

- **调度模块 (runner)**:
  - 负责管理任务的添加、启动、停止以及循环执行。
  - 支持高并发任务运行，使用 `sync.WaitGroup` 确保任务优雅退出。

## 使用方法

### 1. 克隆项目
```bash
git clone https://github.com/your-repo/autotx.git
cd autotx
```

### 2. 安装依赖
确保您已安装 Go 1.18+。

```bash
go mod tidy
```

### 3. 运行项目

#### 本地运行
```bash
go run service.go
```

#### 使用运行脚本
```bash
./run.sh
```

#### 使用构建脚本
```bash
./build-go.sh
./build/server/autotx
```

### 4. 使用 Docker
构建 Docker 镜像：
```bash
docker build -t autotx -f DockerfileGo .
```

运行容器：
```bash
docker run -it autotx
```

## 测试

运行单元测试：
```bash
go test ./runner/ ./task/
```

## 扩展功能

### 添加新任务
在 `task/` 目录中创建一个新文件，例如 `my_task.go`，实现 `Task` 接口：
```go
type MyTask struct{}

func (t *MyTask) Run(ctx context.Context) error {
    // 实现任务逻辑
    return nil
}

func (t *MyTask) Stop() error {
    // 实现停止逻辑
    return nil
}
```

在 `service.go` 中添加任务：
```go
runner.AddTask(&MyTask{})
```

### 修改运行逻辑
可在 `runner/runner.go` 中调整任务的调度逻辑，例如增加优先级或动态加载任务。

## 贡献
欢迎贡献代码！请通过提交 Pull Request 提交您的更改。

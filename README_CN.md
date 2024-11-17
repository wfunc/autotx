
# AUTOTX

`AUTOTX` 是一个基于 Go 的自动化任务执行框架。它旨在使用 `chromedp` 实现任务，例如网页登录、数据获取以及浏览器交互。项目采用模块化设计，便于扩展和高效管理任务。

> 🌐 **其他语言版本:**  
> [🇬🇧 English](./README.md)

## 目录结构

```
AUTOTX
├── build/server/         # 编译后的可执行文件
│   └── autotx
├── runner/               # 任务运行模块
│   ├── runner_test.go    # 运行模块的单元测试
│   └── runner.go         # Runner 实现
├── task/                 # 任务定义和逻辑
│   ├── base.go           # BaseTask 定义公共逻辑
│   ├── example.go        # 示例任务实现
│   ├── items.go          # 其他任务逻辑
│   ├── login.go          # 登录任务逻辑
│   ├── sign_in.go        # 注册任务逻辑
│   ├── task.go           # Task 接口与工具
├── .gitignore            # Git 忽略文件
├── autotx                # 构建输出
├── build-go.sh           # Go 构建脚本
├── DockerfileGo          # Docker 构建配置
├── go.mod                # Go 模块配置
├── go.sum                # Go 模块依赖
├── README_CN.md          # 中文文档
├── README.md             # 英文文档（默认）
├── run.sh                # 项目运行脚本
└── service.go            # 服务入口
```

## 功能

- **模块化设计**:
  - 集中管理 `BaseTask`，定义共享属性和方法。
  - 可扩展的任务定义（例如 `LoginTask`, `ExampleTask`）。
- **浏览器自动化**:
  - 使用 `chromedp` 控制无头浏览器。
- **任务管理器**:
  - 管理任务的生命周期（启动、停止、循环执行）。

## 使用方法

### 1. 克隆项目
```bash
git clone https://github.com/your-repo/autotx.git
cd autotx
```

### 2. 安装依赖
确保已安装 Go 1.18+。
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

#### 使用 Docker
```bash
docker build -t autotx -f DockerfileGo .
docker run -it autotx
```

## 测试

运行单元测试:
```bash
go test ./runner/ ./task/
```

## 添加新任务

1. 在 `task/` 目录下创建一个新文件，例如 `my_task.go`，并实现 `Task` 接口。
2. 使用 `BaseTask` 继承共享逻辑。
3. 在运行器或服务逻辑中注册新任务。

---

## 贡献

欢迎贡献代码！请提交 Pull Request。

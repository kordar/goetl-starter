# goetl-starter

`goetl-starter` 用于将 **goetl 引擎（`*engine.Engine`）** 以“模块化配置”的方式批量装配与注册，方便在应用启动阶段按配置加载多个 ETL 任务实例，并通过 `Get(id)` 获取对应引擎对象。

## 安装

```bash
go get github.com/kordar/goetl-starter
```

## 设计

- `GoetlModule`：负责解析配置并逐条调用 `load` 回调构建引擎
- `EngineLoader`：用户提供的回调，用于将单个 item 配置装配为 `*engine.Engine`
- 注册表：`Provide(id, eng)` 保存；`Get(id)` 获取

核心代码见：

- [module.go](file:///Users/mac/Documents/Projects/golang/github.com/etl/goetl-starter/module.go)
- [helper.go](file:///Users/mac/Documents/Projects/golang/github.com/etl/goetl-starter/helper.go)

## 配置形态

支持两种形态：

### 1) 单实例（包含 id）

```go
cfg := map[string]any{
  "id": "job_users",
  "source": map[string]any{ /* ... */ },
  "sink":   map[string]any{ /* ... */ },
}
```

### 2) 多实例（map key 作为 id）

```go
cfg := map[string]any{
  "job_users": map[string]any{
    "source": map[string]any{ /* ... */ },
  },
  "job_orders": map[string]any{
    "source": map[string]any{ /* ... */ },
  },
}
```

## 使用示例

```go
package main

import (
  "context"
  "log"

  "github.com/kordar/goetl/engine"
  goetlstarter "github.com/kordar/goetl-starter"
)

func main() {
  m := goetlstarter.NewGoetlModule("goetl", func(moduleName, itemID string, item map[string]any) (*engine.Engine, error) {
    _ = moduleName
    _ = itemID
    _ = item
    eng := &engine.Engine{}
    return eng, nil
  })

  m.Load(map[string]any{
    "job_a": map[string]any{"source": map[string]any{}},
    "job_b": map[string]any{"source": map[string]any{}},
  })

  eng := goetlstarter.Get("job_a")
  if err := eng.Run(context.Background()); err != nil {
    log.Fatal(err)
  }
}
```

## 测试

```bash
go test ./...
go vet ./...
```

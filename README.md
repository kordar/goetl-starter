# goetl-starter

`goetl-starter` 用于批量加载与注册 `*engine.Engine`，适合在应用启动阶段按配置装配多个 ETL 任务。

## 核心设计

- `GoetlModule`：读取配置并逐项调用 `EngineLoader`
- `EngineLoader`：你提供的回调，负责把一条配置构建成 `*engine.Engine`
- `Provide/Get`：全局注册与获取引擎实例

## 配置输入形态

支持两种：

1. 单实例（包含 `id`）

```go
cfg := map[string]any{
	"id": "job_users",
	"source": map[string]any{},
	"sink": map[string]any{},
}
```

2. 多实例（map key 作为 id）

```go
cfg := map[string]any{
	"job_users": map[string]any{"source": map[string]any{}},
	"job_orders": map[string]any{"source": map[string]any{}},
}
```

## 使用示例

```go
package main

import (
	"context"
	"log"

	"github.com/kordar/goetl"
	"github.com/kordar/goetl/engine"
	goetlstarter "github.com/kordar/goetl-starter"
)

func main() {
	m := goetlstarter.NewGoetlModule("goetl", func(moduleName, itemID string, item map[string]any) (*engine.Engine, error) {
		_ = moduleName
		_ = item

		eng := engine.NewEngine(
			/* your sink */,
			engine.WithPipeline(goetl.NewPipeline()),
		)
		return eng, nil
	})

	m.Load(map[string]any{
		"job_a": map[string]any{},
		"job_b": map[string]any{},
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

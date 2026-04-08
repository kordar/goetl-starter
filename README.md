# goetl-starter

Go 版 ETL 启动器，提供模块化的引擎装载能力和 DBWalker 规范化封装，便于在应用启动阶段按配置或代码快速组装多个 ETL 任务。

## 组件

- GoetlModule：读取配置并逐项调用 EngineLoader，装配引擎
- Provide/Get：全局注册与获取引擎实例
- DBWalkerSpec：标准化数据库扫描规范接口
- BaseDBWalkerSpec：可嵌入的基础对象，覆盖必要方法即可复用默认行为
- LoadDBWalkerWithGorm/LoadDBWalkerWithSQL：将 DBWalkerSpec 转为标准 DBWalker 并注入 Engine

## DBWalkerSpec

必选方法
- CheckpointKey() string
- PageSize() int
- MaxItems() int
- Store() checkpoint.CheckpointStore
- BuildQuery(ctx, cur, limit)
- ExtractCursor(row)
- MapRow(row)
- BuildAck(row, cur)

可选扩展（实现 DBWalkerSpecOptional）
- Name() string
- Partition() string
- BuildAttrs(row, cur) (map[string]any, error)

## 快速开始

最小实现一个按 id 递增的扫描器，每秒拉取一条记录：

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kordar/goetl/checkpoint"
	"github.com/kordar/goetl/engine"
	goetlstarter "github.com/kordar/goetl-starter"
	"gorm.io/gorm"
)

type MySpec struct {
	goetlstarter.BaseDBWalkerSpec
	Query string
}

func (s MySpec) BuildQuery(ctx context.Context, cur checkpoint.Cursor, limit int) (string, []any, error) {
	last := int64(0)
	if len(cur.Values) > 0 {
		if v, ok := cur.Values[0].(int64); ok { last = v }
	}
	return s.Query, []any{last, limit}, nil
}

func (s MySpec) ExtractCursor(row map[string]any) (checkpoint.Cursor, error) {
	return checkpoint.Cursor{Values: []any{row["id"].(int64)}}, nil
}

func main() {
	var gdb *gorm.DB

	spec := MySpec{
		BaseDBWalkerSpec: goetlstarter.BaseDBWalkerSpec{
			CheckpointKeyValue: "t:id",
			PageSizeValue: 1,
			MaxItemsValue: 1,
			NameValue: "walker",
			PartitionValue: "t",
		},
		Query: "SELECT id, name FROM t WHERE id > ? ORDER BY id LIMIT ?",
	}

	eng := engine.NewEngine()
	goetlstarter.LoadDBWalkerWithGorm(eng, gdb, spec, time.Second, time.Second, true)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eng.Run(ctx, func(m goetl.Message) {
		fmt.Println("msg", m)
	}, func(err error) {
		fmt.Println("err", err)
	})
	time.Sleep(5 * time.Second)
}
```

## 使用 GoetlModule 按配置装载

```go
m := goetlstarter.NewGoetlModule("etl", func(moduleName, itemID string, item map[string]any) (*engine.Engine, error) {
	eng := engine.NewEngine()
	// 依据 item 构造 spec（省略），然后：
	// goetlstarter.LoadDBWalkerWithGorm(eng, gdb, spec, interval, retry, true)
	return eng, nil
})

m.Load(map[string]any{"job_a": map[string]any{}})
eng := goetlstarter.Get("job_a")
```

## 测试

```bash
go test ./...
go vet ./...
```

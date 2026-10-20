// 命令 table-demo 演示表格类型 *Table 的
// 构建、选列、设索引与 JSON 序列化，不需要运行组件参与。
package main

import (
	"fmt"

	thsdk "thsdk"
)

func main() {
	// 手工构造一个表格
	t := &thsdk.Table{
		Columns: []string{"code", "close", "open"},
		Rows: []map[string]any{
			{"code": "000001", "close": 10.5, "open": 10.0},
			{"code": "600000", "close": 11.0, "open": 11.2},
		},
		Index: []any{"a", "b"},
	}
	fmt.Println("列：", t.Columns, "行数：", t.RowCount(), "空：", t.Empty())

	// 选取字段
	selected, err := t.Select([]string{"code", "close"})
	if err != nil {
		panic(err)
	}
	fmt.Println("选取后列：", selected.Columns)

	// 设索引（把列提升为 Index 并从数据列移除）
	selected.SetIndex("code")
	fmt.Println("索引名：", selected.IndexName, "索引：", selected.Index, "剩余列：", selected.Columns)

	// JSON 序列化
	data, err := t.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Printf("JSON：%s\n", data)
}

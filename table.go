package thsdk

// 本文件实现 pandas.DataFrame 的 Go 等价类型 *Table，并移植 api.py 中
// _frame / _records_frame / _identity_rows / _dataframe_result 的全部逻辑。
//
// Table 与 DataFrame 的对应关系：
//   - Columns ↔ DataFrame.columns（有序列名）
//   - Rows    ↔ 行记录（[]map[string]any，缺失列值为 nil，对应 NaN）
//   - Index   ↔ DataFrame.index（set_index 后生效）
//   - Attrs   ↔ DataFrame.attrs（thsdk 元数据等）

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Table 是 pandas.DataFrame 的等价实现，用于承载公开层表格结果。
type Table struct {
	Columns   []string
	Rows      []map[string]any
	Index     []any
	IndexName string
	Attrs     map[string]any
}

// eightDigits 对应 _frame 中的 ^\d{8}$ 全匹配判断。
var eightDigits = regexp.MustCompile(`^\d{8}$`)

// NewTable 创建空表格。
func NewTable() *Table {
	return &Table{Attrs: map[string]any{}}
}

// Empty 判断表格是否没有行。
func (t *Table) Empty() bool {
	return t == nil || len(t.Rows) == 0
}

// RowCount 返回行数。
func (t *Table) RowCount() int {
	if t == nil {
		return 0
	}
	return len(t.Rows)
}

// ColumnCount 返回列数。
func (t *Table) ColumnCount() int {
	if t == nil {
		return 0
	}
	return len(t.Columns)
}

// HasColumn 判断是否存在指定列。
func (t *Table) HasColumn(name string) bool {
	if t == nil {
		return false
	}
	for _, column := range t.Columns {
		if column == name {
			return true
		}
	}
	return false
}

// SetAttr 设置 attrs 元数据。
func (t *Table) SetAttr(key string, value any) {
	if t.Attrs == nil {
		t.Attrs = map[string]any{}
	}
	t.Attrs[key] = value
}

// Attr 读取 attrs 元数据。
func (t *Table) Attr(key string) any {
	if t.Attrs == nil {
		return nil
	}
	return t.Attrs[key]
}

// JSON 把表格序列化为 {"columns","rows","index","index_name","attrs"}。
func (t *Table) JSON() ([]byte, error) {
	return marshalTable(t)
}

// String 实现 fmt.Stringer，输出紧凑 JSON。
func (t *Table) String() string {
	data, err := marshalTable(t)
	if err != nil {
		return fmt.Sprintf("<Table error: %v>", err)
	}
	return string(data)
}

// newTableFromRecords 从记录集合构造表格，对应 pandas.DataFrame.from_records：
// 列顺序为各记录键首次出现的顺序，缺失键值为 nil。
//
// Go 的 map 迭代顺序是随机的（Python dict 保持插入序），因此单个记录内的
// 键先按字典序排序，使列顺序确定且可复现；跨记录仍保持首次出现语义。
func newTableFromRecords(records []map[string]any) *Table {
	var columns []string
	seen := make(map[string]bool)
	for _, record := range records {
		keys := sortedKeys(record)
		for _, key := range keys {
			if !seen[key] {
				seen[key] = true
				columns = append(columns, key)
			}
		}
	}
	rows := make([]map[string]any, len(records))
	for i, record := range records {
		row := make(map[string]any, len(columns))
		for _, column := range columns {
			row[column] = record[column]
		}
		rows[i] = row
	}
	return &Table{Columns: columns, Rows: rows, Attrs: map[string]any{}}
}

// sortedKeys 返回映射键的字典序排列。
func sortedKeys(mapping map[string]any) []string {
	keys := make([]string, 0, len(mapping))
	for key := range mapping {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// keyEntry 是映射排序后的条目。
type keyEntry struct {
	key   any
	value any
}

// sortedMappingEntries 对应 _records_frame 中的排序逻辑：
// 优先按整数键升序，其余按字符串键；可选 mapping_order 优先置顶。
func sortedMappingEntries(mapping map[string]any, mappingOrder any) []keyEntry {
	entries := make([]keyEntry, 0, len(mapping))
	for key, value := range mapping {
		entries = append(entries, keyEntry{key: key, value: value})
	}
	order, ok := toSequenceAny(mappingOrder)
	if !ok || len(order) == 0 {
		sortEntries(entries)
		return entries
	}
	var ordered []keyEntry
	remaining := entries
	for _, requested := range order {
		requestedString := fmt.Sprint(requested)
		for index, entry := range remaining {
			if entry.key == requested || fmt.Sprint(entry.key) == requestedString {
				ordered = append(ordered, entry)
				remaining = append(remaining[:index], remaining[index+1:]...)
				break
			}
		}
	}
	sortEntries(remaining)
	return append(ordered, remaining...)
}

func sortEntries(entries []keyEntry) {
	sort.SliceStable(entries, func(i, j int) bool {
		leftRank, leftKey := mappingSortKey(entries[i].key)
		rightRank, rightKey := mappingSortKey(entries[j].key)
		if leftRank != rightRank {
			return leftRank < rightRank
		}
		return leftKey < rightKey
	})
}

// mappingSortKey 对应 Python mapping_sort_key：整数键排前，字符串键排后。
func mappingSortKey(key any) (int, string) {
	if value, ok := asInt64(key); ok {
		return 0, fmt.Sprintf("%020d", value)
	}
	return 1, fmt.Sprint(key)
}

// recordsFrame 对应 Python _records_frame：把明确的记录集合转换为表格，
// 不展开行内嵌套值。
func recordsFrame(rows any, mappingKeys bool, mappingOrder any) (*Table, error) {
	switch value := rows.(type) {
	case nil:
		return NewTable(), nil
	case map[string]any:
		if !mappingKeys {
			return newTableFromRecords([]map[string]any{value}), nil
		}
		entries := sortedMappingEntries(value, mappingOrder)
		records := make([]map[string]any, 0, len(entries))
		indexValues := make([]any, 0, len(entries))
		for _, entry := range entries {
			if item, ok := entry.value.(map[string]any); ok {
				records = append(records, item)
			} else {
				records = append(records, map[string]any{"value": entry.value})
			}
			indexValues = append(indexValues, entry.key)
		}
		table := newTableFromRecords(records)
		table.Index = indexValues
		table.IndexName = "group_key"
		return table, nil
	case string:
		return newTableFromRecords([]map[string]any{{"value": value}}), nil
	}

	values, ok := toSequenceAny(rows)
	if !ok {
		values = []any{rows}
	}
	if len(values) == 0 {
		return NewTable(), nil
	}
	allMaps := true
	for _, item := range values {
		if _, isMap := item.(map[string]any); !isMap {
			allMaps = false
			break
		}
	}
	if allMaps {
		records := make([]map[string]any, 0, len(values))
		for _, item := range values {
			records = append(records, item.(map[string]any))
		}
		return newTableFromRecords(records), nil
	}
	return newTableFromRecords([]map[string]any{{"value": values}}), nil
}

// frame 对应 Python _frame：可选把 time 列转为日期索引。
func frame(rows any, timeIndex bool) (*Table, error) {
	var table *Table
	if existing, ok := rows.(*Table); ok {
		table = existing.clone()
	} else {
		built, err := recordsFrame(rows, false, nil)
		if err != nil {
			return nil, err
		}
		table = built
	}
	if timeIndex && table.HasColumn("time") && !table.Empty() {
		allMatch := true
		for _, row := range table.Rows {
			text, isString := row["time"].(string)
			if !isString || !eightDigits.MatchString(text) {
				allMatch = false
				break
			}
		}
		if allMatch {
			for _, row := range table.Rows {
				text := row["time"].(string)
				if parsed, err := time.Parse("20060102", text); err == nil {
					row["time"] = parsed
				}
			}
		}
	}
	if timeIndex && table.HasColumn("time") {
		table.SetIndex("time")
	}
	return table, nil
}

// clone 深拷贝表格（对应 pandas DataFrame.copy(deep=True)）。
func (t *Table) clone() *Table {
	copied := &Table{
		Columns:   append([]string(nil), t.Columns...),
		Rows:      make([]map[string]any, len(t.Rows)),
		Index:     append([]any(nil), t.Index...),
		IndexName: t.IndexName,
		Attrs:     map[string]any{},
	}
	for i, row := range t.Rows {
		copied.Rows[i] = deepCopyAny(row).(map[string]any)
	}
	for key, value := range t.Attrs {
		copied.Attrs[key] = deepCopyAny(value)
	}
	return copied
}

// SetIndex 把指定列设为索引并从数据列中移除，对应 pandas set_index。
func (t *Table) SetIndex(name string) {
	index := make([]any, len(t.Rows))
	for i, row := range t.Rows {
		index[i] = row[name]
	}
	t.Index = index
	t.IndexName = name
	var columns []string
	for _, column := range t.Columns {
		if column != name {
			columns = append(columns, column)
		}
	}
	t.Columns = columns
	for _, row := range t.Rows {
		delete(row, name)
	}
}

// Select 按字段名选取列，对应 Python get_price 中的 fields 过滤。
func (t *Table) Select(fields []string) (*Table, error) {
	selected := make([]string, len(fields))
	for i, field := range fields {
		selected[i] = fmt.Sprint(field)
	}
	var missing []string
	for _, field := range selected {
		if !t.HasColumn(field) {
			missing = append(missing, field)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("返回结果缺少字段: %s", strings.Join(missing, ", "))
	}
	table := &Table{
		Columns:   selected,
		Index:     append([]any(nil), t.Index...),
		IndexName: t.IndexName,
		Attrs:     t.Attrs,
		Rows:      make([]map[string]any, 0, len(t.Rows)),
	}
	for _, row := range t.Rows {
		newRow := make(map[string]any, len(selected))
		for _, field := range selected {
			newRow[field] = row[field]
		}
		table.Rows = append(table.Rows, newRow)
	}
	return table, nil
}

// identityRows 对应 Python _identity_rows：把行内嵌套的 security 展平为
// code / market / full_code 列（仅当这些列尚未存在时补充）。
func identityRows(rows any) []map[string]any {
	values, ok := toSequenceAny(rows)
	if !ok {
		return nil
	}
	result := make([]map[string]any, 0, len(values))
	for _, item := range values {
		row, isMap := item.(map[string]any)
		if !isMap {
			continue
		}
		copied := make(map[string]any, len(row))
		for key, value := range row {
			copied[key] = value
		}
		security, _ := copied["security"].(map[string]any)
		delete(copied, "security")
		if security != nil {
			code := mapString(security, "code")
			market := mapString(security, "market")
			if _, present := copied["code"]; !present {
				copied["code"] = code
			}
			if _, present := copied["market"]; !present {
				copied["market"] = market
			}
			if _, present := copied["full_code"]; !present {
				if market != "" && code != "" {
					copied["full_code"] = market + code
				} else {
					copied["full_code"] = ""
				}
			}
		}
		result = append(result, copied)
	}
	return result
}

// mapString 取映射中的字符串值，缺失或 nil 返回空串。
func mapString(mapping map[string]any, key string) string {
	value, ok := mapping[key]
	if !ok || value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

// dataframeResult 按 manifest 中的稳定主集合字段构造公开表格结果，
// 对应 Python _dataframe_result。
func dataframeResult(method string, result any) (*Table, error) {
	spec, ok := nativeMethods[method]
	if !ok {
		return nil, newTHSDKError("未知的 Canonical 方法: %s", method)
	}
	rowField := spec.DataFrameRowsField
	var rows any
	metadata := map[string]any{}
	if rowField == "" {
		rows = result
	} else if result == nil {
		rows = nil
	} else if mapping, isMap := result.(map[string]any); isMap {
		rows = mapping[rowField]
		for key, value := range mapping {
			if key != rowField {
				metadata[key] = deepCopyAny(value)
			}
		}
	} else {
		// 兼容旧运行组件偶尔直接返回记录数组的情况；类型仍稳定为 Table。
		rows = result
	}

	mappingKeys := false
	if _, isMap := rows.(map[string]any); rowField == "groups" && isMap {
		mappingKeys = true
	}
	var mappingOrder any
	if order, present := metadata["group_order"]; present {
		mappingOrder = order
	}
	table, err := recordsFrame(rows, mappingKeys, mappingOrder)
	if err != nil {
		return nil, err
	}
	table.SetAttr("thsdk", map[string]any{
		"method":               spec.Name,
		"declared_result_type": spec.ResultType,
		"collection_key":       rowField,
		"metadata":             metadata,
	})
	return table, nil
}

// marshalTable 序列化表格的完整状态。
func marshalTable(t *Table) ([]byte, error) {
	if t == nil {
		t = NewTable()
	}
	return json.Marshal(map[string]any{
		"columns":    t.Columns,
		"rows":       t.Rows,
		"index":      t.Index,
		"index_name": t.IndexName,
		"attrs":      t.Attrs,
	})
}

// toSequenceAny 把任意切片值转为 []any。
func toSequenceAny(value any) ([]any, bool) {
	switch typed := value.(type) {
	case []any:
		return typed, true
	case nil:
		return nil, false
	}
	return toSequence(value)
}

func deepCopyAny(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(typed))
		for key, item := range typed {
			out[key] = deepCopyAny(item)
		}
		return out
	case []any:
		out := make([]any, len(typed))
		for i, item := range typed {
			out[i] = deepCopyAny(item)
		}
		return out
	default:
		return value
	}
}

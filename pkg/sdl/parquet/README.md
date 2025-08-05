# Apache Parquet Implementation

這個目錄包含了完整的Apache Parquet實現，展示了如何在Go中使用Parquet進行高效的列式數據存儲和處理。Parquet是一種開源的列式存儲格式，特別適合大數據分析和ETL工作流。

## 📁 目錄結構

```
pkg/sdl/parquet/
├── models.go              # Parquet數據模型定義
├── simple_manager.go      # 基本Parquet文件操作管理器
├── workflows.go           # 數據處理工作流示例
├── *_test.go             # 測試文件
├── benchmark_test.go      # 性能測試
├── workflows_test.go      # 工作流測試
└── README.md             # 本文件
```

## 🚀 快速開始

### 環境要求

1. **Go 1.19+**
2. **依賴庫**: 
   ```bash
   go get github.com/segmentio/parquet-go
   ```

### 基本使用示例

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "go-transport-prac/pkg/sdl/parquet"
)

func main() {
    // 創建管理器
    manager := parquet.NewSimpleManager("data/parquet")
    
    // 創建示例數據
    users := []parquet.User{
        {
            ID:     1,
            Email:  "john@example.com",
            Name:   "John Doe",
            Status: "active",
            Profile: &parquet.Profile{
                FirstName: "John",
                LastName:  "Doe",
                Phone:     "+1-555-0123",
                Address: &parquet.Address{
                    Street:     "123 Main St",
                    City:       "New York",
                    State:      "NY",
                    PostalCode: "10001",
                    Country:    "USA",
                },
                Interests: []string{"technology", "sports"},
                Metadata: map[string]string{
                    "source": "manual",
                    "type":   "premium",
                },
            },
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
    }
    
    // 寫入Parquet文件
    err := manager.WriteUsers("users.parquet", users)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("✓ Data written to Parquet file")
    
    // 讀取Parquet文件
    readUsers, err := manager.ReadUsers("users.parquet")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("✓ Read %d users from Parquet file\n", len(readUsers))
    
    // 獲取文件信息
    info, err := manager.GetBasicFileInfo("users.parquet")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("File info: %d rows, %d bytes\n", info.NumRows, info.FileSize)
}
```

## 🧪 運行測試

### 運行所有測試

```bash
cd /path/to/go-transport-prac
go test ./pkg/sdl/parquet -v
```

### 運行特定測試類型

```bash
# 基本操作測試
go test ./pkg/sdl/parquet -v -run TestSimpleParquetOperations

# 產品操作測試
go test ./pkg/sdl/parquet -v -run TestProductOperations

# ETL工作流測試
go test ./pkg/sdl/parquet -v -run TestETLWorkflow

# 批處理測試
go test ./pkg/sdl/parquet -v -run TestBatchProcessing

# 數據質量測試
go test ./pkg/sdl/parquet -v -run TestDataQualityCalculation
```

### 查看測試覆蓋率

```bash
go test ./pkg/sdl/parquet -cover
```

詳細覆蓋率報告：

```bash
go test ./pkg/sdl/parquet -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 📊 性能測試

### 運行所有基準測試

```bash
go test -bench=. -benchmem ./pkg/sdl/parquet
```

### 運行特定基準測試

```bash
# Parquet vs JSON 序列化比較
go test -bench=BenchmarkParquet.*Serialization -benchmem ./pkg/sdl/parquet
go test -bench=BenchmarkJSON.*Serialization -benchmem ./pkg/sdl/parquet

# 反序列化性能
go test -bench=BenchmarkParquet.*Deserialization -benchmem ./pkg/sdl/parquet
go test -bench=BenchmarkJSON.*Deserialization -benchmem ./pkg/sdl/parquet

# 大小比較
go test -bench=BenchmarkParquetVsJSONSize -benchmem ./pkg/sdl/parquet

# 完整週期測試
go test -bench=BenchmarkParquet.*FullCycle -benchmem ./pkg/sdl/parquet
go test -bench=BenchmarkJSON.*FullCycle -benchmem ./pkg/sdl/parquet

# 不同數據集大小測試
go test -bench=BenchmarkParquet.*Dataset -benchmem ./pkg/sdl/parquet

# 內存使用測試
go test -bench=BenchmarkParquetMemoryUsage -benchmem ./pkg/sdl/parquet
```

### 基準測試結果解讀

```
BenchmarkParquetUserSerialization-16    1206    917806 ns/op   6985921 B/op   4798 allocs/op
```

- `1206`: 執行次數
- `917806 ns/op`: 每次操作平均耗時 (納秒)
- `6985921 B/op`: 每次操作平均內存分配 (字節)
- `4798 allocs/op`: 每次操作平均內存分配次數
- `-16`: 使用的CPU核心數

### 性能基準結果總結

基於我們的基準測試結果：

| 指標 | Parquet | JSON | Parquet優勢 |
|------|---------|------|-------------|
| 序列化速度 | 918μs | 824μs | JSON快 11% |
| 反序列化速度 | 1525μs | 2630μs | **Parquet快 72%** |
| 數據大小 | 174KB | 456KB | **Parquet小 62%** |
| 完整週期 | 604μs | 337μs | JSON快 79% |

**結論**: Parquet在數據大小和反序列化性能上具有顯著優勢，特別適合大數據場景。

## 🔄 數據處理工作流

### ETL工作流示例

```go
// 創建數據處理管道
pipeline := parquet.NewDataPipeline("data/pipeline")

// 運行完整的ETL工作流
err := pipeline.RunETLWorkflow()
if err != nil {
    log.Fatal(err)
}

// 清理工作流文件
defer pipeline.CleanupWorkflow()
```

### 批處理工作流

```go
// 運行批處理工作流
err := pipeline.RunBatchProcessing()
if err != nil {
    log.Fatal(err)
}
```

### 分析工作流

```go
// 運行分析工作流
err := pipeline.RunAnalyticsWorkflow()
if err != nil {
    log.Fatal(err)
}
```

## 📈 Parquet特點與優勢

### 核心優勢

1. **列式存儲**: 高效的數據壓縮和查詢性能
2. **Schema演進**: 支持添加、删除字段而不影響兼容性
3. **壓縮效率**: 相比JSON節省約62%的存儲空間
4. **跨平台**: 支持多種編程語言和數據處理框架
5. **元數據豐富**: 包含詳細的Schema和統計信息

### 適用場景

**優先選擇Parquet:**
- 大數據分析和ETL管道
- 數據倉庫和數據湖存儲
- 批處理和OLAP查詢
- 長期數據歸檔
- 需要高壓縮比的場景
- 列式查詢和聚合操作

**可考慮其他格式:**
- 實時數據交換 (考慮JSON)
- 簡單的鍵值存儲 (考慮JSON/MessagePack)
- 流式處理 (考慮Avro)
- 前端應用 (考慮JSON)

### 性能特徵

| 方面 | Parquet表現 | 說明 |
|------|-------------|------|
| 讀取性能 | ⭐⭐⭐⭐⭐ | 列式存儲，支持謂詞下推 |
| 寫入性能 | ⭐⭐⭐⭐ | 需要構建列式結構，稍慢於行式 |
| 壓縮率 | ⭐⭐⭐⭐⭐ | 優秀的壓縮算法支持 |
| Schema靈活性 | ⭐⭐⭐⭐ | 支持複雜嵌套結構 |
| 跨語言支持 | ⭐⭐⭐⭐⭐ | 廣泛的生態系統支持 |

## 🛠️ 數據模型

### 支持的數據類型

```go
// 基本數據類型
type User struct {
    ID        int64     `parquet:"id"`           // 64位整數
    Email     string    `parquet:"email"`        // UTF-8字符串
    Name      string    `parquet:"name"`         
    Status    string    `parquet:"status"`       
    Profile   *Profile  `parquet:"profile"`      // 嵌套結構
    CreatedAt time.Time `parquet:"created_at"`   // 時間戳
    UpdatedAt time.Time `parquet:"updated_at"`   
}

// 嵌套結構
type Profile struct {
    FirstName string            `parquet:"first_name"`
    LastName  string            `parquet:"last_name"`
    Phone     string            `parquet:"phone,optional"`     // 可選字段
    Address   *Address          `parquet:"address,optional"`   // 可選嵌套結構
    Interests []string          `parquet:"interests"`          // 字符串數組
    Metadata  map[string]string `parquet:"metadata"`           // 鍵值對
}

// 複雜數據類型示例
type Product struct {
    ID            int64                 `parquet:"id"`
    Price         *Price                `parquet:"price"`          // 嵌套結構
    Categories    []string              `parquet:"categories"`     // 數組
    Specifications map[string]string   `parquet:"specifications"` // Map
    Inventory     *Inventory            `parquet:"inventory"`      // 另一個嵌套結構
}
```

### Parquet標籤說明

- `parquet:"field_name"`: 指定字段名
- `parquet:"field_name,optional"`: 可選字段
- 支持的Go類型: `int64`, `int32`, `string`, `bool`, `float32`, `float64`, `time.Time`
- 支持複雜類型: `struct`, `slice`, `map`

## 💡 最佳實踐

### Schema設計

1. **字段命名**: 使用snake_case命名約定
2. **可選字段**: 為可能缺失的字段添加`optional`標籤
3. **嵌套結構**: 合理使用嵌套結構組織相關數據
4. **數據類型**: 選擇合適的數據類型以優化存儲和性能

### 性能優化

1. **批量操作**: 一次處理多條記錄而不是逐條處理
2. **內存管理**: 對大數據集使用流式處理
3. **壓縮算法**: 根據數據特徵選擇合適的壓縮算法
4. **文件大小**: 保持合理的文件大小 (通常128MB-1GB)

### 錯誤處理

```go
// 良好的錯誤處理示例
func processParquetFile(filename string) error {
    manager := parquet.NewSimpleManager("data")
    
    users, err := manager.ReadUsers(filename)
    if err != nil {
        return fmt.Errorf("failed to read parquet file %s: %w", filename, err)
    }
    
    if len(users) == 0 {
        return fmt.Errorf("no users found in file %s", filename)
    }
    
    // 處理數據...
    
    return nil
}
```

## 🔧 工作流集成

### ETL管道

```go
// 1. Extract (提取)
rawData, err := extractFromDatabase()
if err != nil {
    return err
}

// 2. Transform (轉換)
cleanData := transformAndValidate(rawData)

// 3. Load (加載)
err = manager.WriteUsers("processed_users.parquet", cleanData)
if err != nil {
    return err
}
```

### 批處理

```go
// 處理大數據集的批次
batchSize := 1000
for i := 0; i < len(allUsers); i += batchSize {
    end := i + batchSize
    if end > len(allUsers) {
        end = len(allUsers)
    }
    
    batch := allUsers[i:end]
    filename := fmt.Sprintf("batch_%d.parquet", i/batchSize)
    
    err := manager.WriteUsers(filename, batch)
    if err != nil {
        return fmt.Errorf("failed to write batch %d: %w", i/batchSize, err)
    }
}
```

## 🐛 故障排除

### 常見問題

1. **標籤語法錯誤**
   ```
   錯誤: `parquet:"id,int64"`
   正確: `parquet:"id"`
   ```

2. **內存不足**
   ```go
   // 對大文件使用流式處理而不是一次性加載
   // 考慮分批處理大數據集
   ```

3. **字段類型不匹配**
   ```
   確保Go結構體字段類型與Parquet schema匹配
   ```

4. **文件權限問題**
   ```bash
   # 確保有讀寫權限
   chmod 755 data/parquet/
   ```

### 調試技巧

```go
// 獲取詳細文件信息
info, err := manager.GetBasicFileInfo("test.parquet")
if err != nil {
    log.Printf("File info error: %v", err)
} else {
    log.Printf("File: %d rows, %d bytes, %d fields", 
        info.NumRows, info.FileSize, len(info.Schema.Fields()))
}

// 列出所有文件
files, err := manager.ListFiles()
if err != nil {
    log.Printf("List files error: %v", err)
} else {
    log.Printf("Available files: %v", files)
}
```

## 📚 相關資源

- [Apache Parquet官方文檔](https://parquet.apache.org/)
- [Go Parquet庫文檔](https://github.com/segmentio/parquet-go)
- [Parquet格式規範](https://github.com/apache/parquet-format)
- [列式存儲介紹](https://en.wikipedia.org/wiki/Column-oriented_DBMS)

## 🔄 與其他格式對比

| 特性 | Parquet | JSON | Avro | Protocol Buffers |
|------|---------|------|------|------------------|
| 存儲效率 | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| 查詢性能 | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| Schema演進 | ⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| 人類可讀 | ⭐ | ⭐⭐⭐⭐⭐ | ⭐ | ⭐ |
| 跨語言支持 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| 適用場景 | 分析、歸檔 | Web、API | 流處理 | 微服務通信 |

每種格式都有其最適合的使用場景，選擇時應根據具體需求考慮性能、兼容性和開發便利性等因素。
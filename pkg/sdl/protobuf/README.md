# Protocol Buffers Implementation

這個目錄包含了完整的Protocol Buffers實現，展示了如何在Go中使用protobuf進行高效的數據序列化和反序列化。

## 📁 目錄結構

```
pkg/sdl/protobuf/
├── proto/                    # Protocol Buffer 定義文件
│   ├── user.proto           # 用戶相關訊息定義
│   ├── product.proto        # 產品相關訊息定義
│   ├── order.proto          # 訂單相關訊息定義
│   ├── common.proto         # 通用訊息定義
│   └── userv2/              # 版本2用戶定義（兼容性測試）
│       └── user_v2.proto
├── gen/                     # 生成的Go代碼
│   ├── user/               # 用戶相關生成代碼
│   ├── product/            # 產品相關生成代碼
│   ├── order/              # 訂單相關生成代碼
│   ├── common/             # 通用生成代碼
│   └── userv2/             # 版本2用戶生成代碼
├── manager.go              # Protocol Buffers管理器
├── examples.go             # 使用示例
├── compatibility.go        # 兼容性演示
├── *_test.go              # 測試文件
├── benchmark_test.go       # 性能測試
├── PERFORMANCE.md          # 性能分析報告
└── README.md              # 本文件
```

## 🚀 快速開始

### 環境要求

1. **Go 1.19+**
2. **Protocol Buffers編譯器**:
   ```bash
   # macOS
   brew install protobuf
   
   # Ubuntu/Debian
   sudo apt install protobuf-compiler
   
   # 其他平台請參考: https://protobuf.dev/downloads/
   ```

3. **Go插件**:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   ```

### 編譯Proto文件

如果需要重新生成Go代碼：

```bash
# 生成所有proto文件
protoc --go_out=. --go_opt=paths=source_relative \
    pkg/sdl/protobuf/proto/*.proto \
    pkg/sdl/protobuf/proto/userv2/*.proto
```

## 💡 使用示例

### 基本序列化和反序列化

```go
package main

import (
    "fmt"
    "log"
    
    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/types/known/timestamppb"
    
    "go-transport-prac/pkg/sdl/protobuf/gen/user"
)

func main() {
    // 創建用戶對象
    user := &user.User{
        Id:    1,
        Email: "john@example.com",
        Name:  "John Doe",
        Status: user.UserStatus_USER_STATUS_ACTIVE,
        Profile: &user.Profile{
            FirstName: "John",
            LastName:  "Doe",
            Phone:     "+1-555-0123",
            Address: &user.Address{
                Street:     "123 Main St",
                City:       "New York",
                State:      "NY",
                PostalCode: "10001",
                Country:    "USA",
            },
        },
        CreatedAt: timestamppb.Now(),
        UpdatedAt: timestamppb.Now(),
    }
    
    // 序列化
    data, err := proto.Marshal(user)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("序列化後大小: %d bytes\n", len(data))
    
    // 反序列化
    var deserializedUser user.User
    err = proto.Unmarshal(data, &deserializedUser)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("用戶名稱: %s\n", deserializedUser.Name)
    fmt.Printf("用戶郵箱: %s\n", deserializedUser.Email)
}
```

### 使用管理器

```go
package main

import (
    "fmt"
    
    "go-transport-prac/pkg/sdl/protobuf"
)

func main() {
    manager := protobuf.NewManager()
    
    // 創建示例用戶
    user := manager.CreateSampleUser()
    
    // 序列化
    data, err := manager.SerializeUser(user)
    if err != nil {
        panic(err)
    }
    
    // 反序列化
    deserializedUser, err := manager.DeserializeUser(data)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("原始用戶ID: %d\n", user.Id)
    fmt.Printf("反序列化用戶ID: %d\n", deserializedUser.Id)
}
```

## 🧪 運行測試

### 運行所有測試

```bash
cd /path/to/go-transport-prac
go test ./pkg/sdl/protobuf -v
```

### 運行特定測試類型

```bash
# 兼容性測試
go test ./pkg/sdl/protobuf -v -run="Test.*Compatibility"

# 序列化測試
go test ./pkg/sdl/protobuf -v -run="TestSerialization"

# 管理器測試
go test ./pkg/sdl/protobuf -v -run="TestManager"
```

### 查看測試覆蓋率

```bash
go test ./pkg/sdl/protobuf -cover
```

詳細覆蓋率報告：

```bash
go test ./pkg/sdl/protobuf -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 📊 性能測試

### 運行所有基準測試

```bash
go test -bench=. -benchmem ./pkg/sdl/protobuf
```

### 運行特定基準測試

```bash
# Protocol Buffers vs JSON 比較
go test -bench=BenchmarkProtobufVsJSON -benchmem ./pkg/sdl/protobuf

# 序列化性能
go test -bench=BenchmarkProtobuf.*Serialization -benchmem ./pkg/sdl/protobuf

# 反序列化性能
go test -bench=BenchmarkProtobuf.*Deserialization -benchmem ./pkg/sdl/protobuf

# 完整週期測試
go test -bench=BenchmarkProtobufFullCycle -benchmem ./pkg/sdl/protobuf

# 大數據集測試
go test -bench=BenchmarkProtobufLargeDataSet -benchmem ./pkg/sdl/protobuf
```

### 性能基準輸出解讀

```
BenchmarkProtobufUserSerialization-16    1720834    686.7 ns/op
```

- `1720834`: 執行次數
- `686.7 ns/op`: 每次操作平均耗時
- `-16`: 使用的CPU核心數

### 生成性能報告

```bash
# 生成基準測試報告
go test -bench=. -benchmem ./pkg/sdl/protobuf > benchmark_results.txt

# 比較兩次測試結果（需要先安裝benchcmp）
go install golang.org/x/tools/cmd/benchcmp@latest
benchcmp old_results.txt new_results.txt
```

## 🔄 兼容性演示

### 運行兼容性演示

```bash
# 運行兼容性演示程序
go run ./pkg/sdl/protobuf/examples.go

# 運行兼容性測試
go test ./pkg/sdl/protobuf -v -run="TestSchemaEvolutionRoundTrip"
```

### 兼容性測試涵蓋範圍

1. **前向兼容性**: 舊代碼讀取新數據
2. **後向兼容性**: 新代碼讀取舊數據
3. **未知字段保存**: 確保未知字段在序列化過程中被保持
4. **枚舉演進**: 新枚舉值的處理
5. **字段添加/刪除**: Schema演進的安全性

## 📈 性能特點

基於我們的基準測試結果：

### Protocol Buffers vs JSON

| 指標 | Protocol Buffers | JSON | 優勢 |
|------|------------------|------|------|
| 序列化速度 | 686.7 ns/op | 665.3 ns/op | JSON略快 3% |
| 反序列化速度 | 903.9 ns/op | 2835 ns/op | **Protobuf快 3.1x** |
| 數據大小 | 245 bytes | 471 bytes | **Protobuf小 48%** |
| 完整週期 | 1657 ns/op | 3566 ns/op | **Protobuf快 2.2x** |

### 適用場景

**優先選擇Protocol Buffers:**
- 微服務間通信
- 移動應用（節省流量）
- 高性能API
- 數據存儲
- 實時系統

**可考慮JSON:**
- 公開REST API
- 前端開發
- 調試和開發階段
- 簡單數據交換

## 🛠️ 開發工具

### VS Code擴展

推薦安裝以下擴展以獲得更好的開發體驗：

- **vscode-proto3**: Protocol Buffers語法高亮
- **Go**: Go語言支持

### Makefile命令

如果項目中有Makefile，可以添加以下命令：

```makefile
# 生成protobuf代碼
.PHONY: proto-gen
proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative \
		pkg/sdl/protobuf/proto/*.proto \
		pkg/sdl/protobuf/proto/userv2/*.proto

# 運行protobuf測試
.PHONY: test-protobuf
test-protobuf:
	go test ./pkg/sdl/protobuf -v

# 運行protobuf基準測試
.PHONY: bench-protobuf
bench-protobuf:
	go test -bench=. -benchmem ./pkg/sdl/protobuf
```

## 🤝 最佳實踐

### Schema設計

1. **字段編號**: 1-15使用單字節編碼，用於常用字段
2. **保留字段**: 使用`reserved`關鍵字保護已刪除的字段
3. **枚舉**: 總是包含`UNSPECIFIED = 0`默認值
4. **可選字段**: 新字段添加為可選，保持兼容性

### 代碼組織

1. **分離關注點**: 不同領域使用不同的proto文件
2. **版本管理**: 為破壞性更改創建新版本
3. **文檔**: 為所有字段和消息添加註釋
4. **測試**: 為所有兼容性場景編寫測試

### 性能優化

1. **重用對象**: 避免頻繁創建新對象
2. **池化**: 對高頻操作使用對象池
3. **流式處理**: 對大數據使用流式API
4. **監控**: 監控序列化性能和數據大小

## 📚 相關資源

- [Protocol Buffers官方文檔](https://protobuf.dev/)
- [Go中的Protocol Buffers](https://protobuf.dev/getting-started/gotutorial/)
- [性能分析報告](./PERFORMANCE.md)
- [兼容性指南](https://protobuf.dev/programming-guides/proto3/#updating)

## 🐛 故障排除

### 常見問題

1. **protoc命令找不到**
   ```bash
   # 確保protoc已安裝並在PATH中
   protoc --version
   ```

2. **Go插件找不到**
   ```bash
   # 確保Go bin目錄在PATH中
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

3. **編譯錯誤**
   ```bash
   # 確保依賴已安裝
   go mod tidy
   ```

4. **測試失敗**
   ```bash
   # 清理並重新生成
   rm -rf pkg/sdl/protobuf/gen/*
   make proto-gen  # 或手動運行protoc命令
   ```

如有其他問題，請查看測試輸出或創建issue。
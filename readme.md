# TiDB测试
## 如何使用:
```go
//tidb连接配置信息
dbConfig := DbConfig{
	Host:     "127.0.0.1",
	Port:     4000,
	User:     "root",
	Password: "",
	DbName:   "test_db",
}
//从文件读取测试sql
reader, _ := NewFileReader([]string{"test1.sql", "test2.sql"， "test3.sql"})
//选择枚举的方式混淆执行顺序
strategy := NewExhaustive()
//随时取消
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
//执行测试，阻塞方法，如果sql执行出错，立刻结束
err := Run(ctx, dbConfig, reader, strategy)
```

## 其他：
* 标注TODO的地方可以继续思考和优化。
* 部分单元测试使用了本地数据库，配置在tidb_test.go/testDbConfig
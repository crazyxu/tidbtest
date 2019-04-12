package tidbtest

//Reader 测试数据
type Reader interface {
	//输出是{输入源名称：文本内容}
	Read() (map[string]string, error)
}

//Strategy 混淆策略
type Strategy interface {
	//输入{输入源名称：sql条数}，返回所有的执行顺序组合
	Shuffle(map[string]int) (turns [][]string)
}

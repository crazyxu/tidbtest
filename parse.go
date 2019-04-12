package tidbtest

import (
	"fmt"
	"strings"
)

//parseSQL 将多行字符串解析成多条sql语句
//按照分号来划分sql，可以自动忽略单行注释和多行注释
func parseSQL(s string) []string {
	lines := strings.Split(s, "\n")
	var sqls, sqlLines []string
	comment := false
	for _, line := range lines {
		s := strings.TrimSpace(line)
		if strings.HasPrefix(s, "/*") {
			comment = true
		} else if strings.HasSuffix(s, "*/") {
			comment = false
		} else if strings.HasPrefix(s, "--") || comment || len(s) == 0 {
			continue
		} else {
			sqlLines = append(sqlLines, line)
			if strings.HasSuffix(s, ";") {
				sqls = append(sqls, strings.Join(sqlLines, "\n"))
				sqlLines = nil
			}
		}
	}
	//最后一行可能没有分号，自动补上
	if len(sqlLines) != 0 {
		sqls = append(sqls, fmt.Sprintf("%s;", strings.Join(sqlLines, "\n")))
	}
	return sqls
}

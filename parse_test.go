package tidbtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSQL1(t *testing.T) {
	sqls := parseSQL(`
--使用；来区分多条
select 1;
/*
自动忽略注释
最后一条sql没有分号会自动加上
*/
select 2;
select 3
	`)
	assert.Equal(t, []string{
		"select 1;",
		"select 2;",
		"select 3;",
	}, sqls)
}

func TestParseSQL2(t *testing.T) {
	sqls := parseSQL(`
	--支持多行语句，保持原有格式
select 1
	from users;
select 2;
	`)
	assert.Equal(t, []string{
		"select 1\n	from users;",
		"select 2;",
	}, sqls)
}

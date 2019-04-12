package tidbtest

import (
	"context"
	"fmt"
)

type worker struct {
	name       string
	sqls       []string
	currentIdx int
	db         *tiDB
}

//Run 在tidb中跑测试sql
//ctx 此方法是阻塞的，ctx能随时停止测试
//dbCfg 连接数据库相关配置
//reader 测试数据源
//strategy 数据混淆策略
//err 如果全部sql执行成功，返回nil
func Run(ctx context.Context, dbCfg DbConfig, reader Reader, strategy Strategy) error {
	//1.获取测试数据，字符串
	ss, err := reader.Read()
	if err != nil {
		return fmt.Errorf("read error %s", err)
	}

	//2.解析成多条sql
	allSQL := make(map[string][]string)
	counts := make(map[string]int)
	for name, s := range ss {
		sqls := parseSQL(s)
		allSQL[name] = sqls
		counts[name] = len(sqls)
	}

	//3.使用策略混淆，产生多种组合
	turns := strategy.Shuffle(counts)

	//4.依次并发执行每种组合，TODO:考虑执行完一种组合后是否需要恢复数据

	//生成多个worker并发执行sql，每个worker一个tidb实例
	var workers []*worker
	var workerFns = make(map[string]workerFn)
	for name, sqls := range allSQL {
		db, err := newTiDB(dbCfg)
		if err != nil {
			return fmt.Errorf("tidb error %s", err)
		}
		defer db.close()

		w := &worker{
			name: name,
			db:   db,
			sqls: sqls,
		}
		workers = append(workers, w)
		workerFns[name] = w.work
	}
	for _, trun := range turns {
		//一种组合测试完后重置计数
		for _, w := range workers {
			w.reset()
		}

		err = workInTurns(ctx, trun, workerFns)
		if err != nil {
			return fmt.Errorf("run test error, order: %s, error: %s", trun, err)
		}
	}

	return nil
}

func (r *worker) work(ctx context.Context) error {
	if r == nil {
		return fmt.Errorf("worker is nil")
	}
	if r.db == nil {
		return fmt.Errorf("db is nil")
	}

	if r.currentIdx > len(r.sqls)-1 {
		return fmt.Errorf("invalid state, currentIdx is %d, sql counts is %d",
			r.currentIdx, len(r.sqls))
	}

	if err := r.db.exec(r.sqls[r.currentIdx]); err != nil {
		return fmt.Errorf(`exec sql "%s" error %s`, r.sqls[r.currentIdx], err)
	}
	r.currentIdx++
	return nil
}

func (r *worker) reset() {
	if r == nil {
		return
	}
	r.currentIdx = 0
}

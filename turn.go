package tidbtest

import (
	"context"
	"fmt"
	"sync"
)

type workerFn func(context.Context) error

//workInTurns 使用类似击鼓传花的方式在多个routine中轮来work
func workInTurns(ctx context.Context, turns []string, fns map[string]workerFn) (err error) {
	if len(turns) == 0 {
		return nil
	}
	//验证turns元素都在fns中
	for _, name := range turns {
		if _, ok := fns[name]; !ok {
			return fmt.Errorf("turn %s not specifiy func", name)
		}
	}

	//用于停止所有启动的routine
	turnCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	//等待routine结束
	wg := &sync.WaitGroup{}
	//routine之间的唤醒机制
	chans := make(map[string]chan struct{}, 0)
	for name := range fns {
		chans[name] = make(chan struct{}, 1)
	}

	for name, fn := range fns {
		wg.Add(1)
		go func(name string, c chan struct{}, fn workerFn) {
			defer wg.Done()
			for {
				//等待被唤醒work,或被ctx取消
				select {
				case <-c:
					//如果work报错，停止整个测试
					err1 := fn(turnCtx)
					if err1 != nil {
						if err == nil {
							cancel()
							err = err1
						}
						return
					}
				case <-turnCtx.Done():
					return
				}

				//如果执行完了，就停止所有等待的routine
				if len(turns) == 0 {
					cancel()
					return
				}

				//下一次该谁了
				turn := turns[0]
				turns = turns[1:len(turns)]
				if turn == name {
					//还是自己
					c <- struct{}{}
				} else {
					//不是自己，唤醒该turn
					chans[turn] <- struct{}{}
				}
			}
		}(name, chans[name], fn)
	}

	//激活第一个worker
	turn := turns[0]
	turns = turns[1:len(turns)]
	chans[turn] <- struct{}{}

	wg.Wait()
	return
}

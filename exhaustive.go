package tidbtest

type exhaustive struct{}

//NewExhaustive 构造一个穷举策略
func NewExhaustive() Strategy {
	return &exhaustive{}
}

//Shuffle 打乱多个输入源的操作顺序
//例如输入A、B分别需要执行2、1个操作
//inputs是{"A":2,"B":1}，输出是[["A","A","B"]，["A","B","A"]，["B","A","A"]]
func (e *exhaustive) Shuffle(inputs map[string]int) (turns [][]string) {
	if e == nil {
		return nil
	}

	//最终输出的序列的长度
	var totalNum int
	var names []string
	var nums []int
	for name, num := range inputs {
		names = append(names, name)
		nums = append(nums, num)
		totalNum += num
	}

	//输出序列的总index
	availableIdx := make([]int, totalNum)
	for i := 0; i < totalNum; i++ {
		availableIdx[i] = i
	}

	idxTurns := selectIdx(0, availableIdx, names, nums)

	turns = make([][]string, len(idxTurns))
	for i, idxTurn := range idxTurns {
		turn := make([]string, totalNum)
		for name, idxs := range idxTurn {
			for _, i := range idxs {
				turn[i] = name
			}
		}
		turns[i] = turn
	}
	return
}

//selectIdx 一组names依次从availableIdx中选择需要的数量的index，结果返回各种可能的不重复的组合
func selectIdx(current int, availableIdx []int, names []string, nums []int) (turns []map[string][]int) {
	if current >= len(names) || len(names) != len(nums) {
		return nil
	}

	name := names[current]
	num := nums[current]

	//最后一个，没得选
	if current == len(names)-1 {
		turns = []map[string][]int{
			{
				name: availableIdx,
			},
		}
		return turns
	}

	//从availableIdx中选择num个
	selections, remainings := cNum(0, num, availableIdx)

	//递归的选择、合并结果
	for i, selection := range selections {
		subTurns := selectIdx(current+1, remainings[i], names, nums)
		for _, t := range subTurns {
			mTurn := map[string][]int{name: selection}
			for k, v := range t {
				mTurn[k] = v
			}
			turns = append(turns, mTurn)
		}
	}
	return
}

//cNum 从nums数组中无序的选择n个元素，其中nums和selections都是单调递增且无重复
//selections是可能性组合,remainings是剩余组合
func cNum(p, n int, nums []int) (selections, remainings [][]int) {
	if n == 0 || len(nums) < n {
		return nil, [][]int{nums}
	}

	for i := 0; i < len(nums); i++ {
		if nums[i] < p {
			continue
		}
		//递归子集
		subNums := []int{}
		subNums = append(subNums, nums[0:i]...)
		subNums = append(subNums, nums[i+1:]...)
		subSelections, subRemainings := cNum(nums[i], n-1, subNums)
		if n == 1 {
			//最后一层
			selections = append(selections, []int{nums[i]})
		} else {
			for _, v := range subSelections {
				selections = append(selections, append([]int{nums[i]}, v...))
			}
		}
		remainings = append(remainings, subRemainings...)
	}
	return
}

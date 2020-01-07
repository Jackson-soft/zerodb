package backend

import (
	"math/rand"
	"time"

	"git.2dfire.net/zerodb/proxy/pkg/errcode"
)

func Gcd(ary []int) int {
	var i int
	min := ary[0]
	length := len(ary)
	for i = 0; i < length; i++ {
		if ary[i] < min {
			min = ary[i]
		}
	}

	for {
		isCommon := true
		for i = 0; i < length; i++ {
			if ary[i]%min != 0 {
				isCommon = false
				break
			}
		}
		if isCommon {
			break
		}
		min--
		if min < 1 {
			break
		}
	}
	return min
}

func (node *HostGroupNode) InitBalancer() {
	var sum int
	node.LastReadIndex = 0
	gcd := Gcd(node.ReadWeights)

	for _, weight := range node.ReadWeights {
		sum += weight / gcd
	}

	node.RoundRobinQ = make([]int, 0, sum)
	for index, weight := range node.ReadWeights {
		for j := 0; j < weight/gcd; j++ {
			node.RoundRobinQ = append(node.RoundRobinQ, index)
		}
	}

	//random order
	if 1 < len(node.ReadWeights) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := 0; i < sum; i++ {
			x := r.Intn(sum)
			temp := node.RoundRobinQ[x]
			other := sum % (x + 1)
			node.RoundRobinQ[x] = node.RoundRobinQ[other]
			node.RoundRobinQ[other] = temp
		}
	}
}

func (node *HostGroupNode) GetNextRead() (*DBPool, error) {
	var index int
	queueLen := len(node.RoundRobinQ)
	if queueLen == 0 {
		return nil, errcode.ErrNoDBPool
	}
	if queueLen == 1 {
		index = node.RoundRobinQ[0]
		return node.Read[index], nil
	}

	node.LastReadIndex = node.LastReadIndex % queueLen
	index = node.RoundRobinQ[node.LastReadIndex]
	if len(node.Read) <= index {
		return nil, errcode.ErrNoDBPool
	}
	DBPool := node.Read[index]
	node.LastReadIndex++
	node.LastReadIndex = node.LastReadIndex % queueLen
	return DBPool, nil
}

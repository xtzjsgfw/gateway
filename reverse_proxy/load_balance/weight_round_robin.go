package load_balance

import (
	"errors"
	"strconv"
)

type WeightRoundRobinLoadBalance struct {
	curIndex int
	rss      []*WeightNode
	rsw      []int
	conf     LoadBalanceConf
}

type WeightNode struct {
	addr            string
	weight          int //权重值
	currentWeight   int //节点当前权重
	effectiveWeight int //有效权重
}

func (w *WeightRoundRobinLoadBalance) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("param len need 2")
	}
	parInt, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}
	node := &WeightNode{addr: params[0], weight: int(parInt)}
	node.effectiveWeight = node.weight
	w.rss = append(w.rss, node)
	return nil
}

func (w *WeightRoundRobinLoadBalance) Next() string {
	total := 0
	var best *WeightNode
	for i := 0; i < len(w.rss); i++ {
		weightNode := w.rss[i]
		//step 1 统计所有有效权重之和
		total += weightNode.effectiveWeight

		//step 2 变更节点临时权重为的节点临时权重+节点有效权重
		weightNode.currentWeight += weightNode.effectiveWeight

		//step 3 有效权重默认与权重相同，通讯异常时-1, 通讯成功+1，直到恢复到weight大小
		if weightNode.effectiveWeight < weightNode.weight {
			weightNode.effectiveWeight++
		}

		//step 4 选择最大临时权重点节点
		if best == nil || weightNode.currentWeight > best.currentWeight {
			best = weightNode
		}
	}
	if best == nil {
		return ""
	}
	//step 5 变更临时权重为 临时权重-有效权重之和
	best.currentWeight -= total
	return best.addr
}

func (w *WeightRoundRobinLoadBalance) Get(key string) (string, error) {
	return w.Next(), nil
}

func (w *WeightRoundRobinLoadBalance) Update() {
}

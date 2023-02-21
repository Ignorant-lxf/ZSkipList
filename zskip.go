package zsetlist

import (
	"math/rand"
	"time"
)

type ZSkipList struct {
	level  int    // 列表的层次
	length uint64 // 跳表的长度

	head *ZSkipListNode
	tail *ZSkipListNode
}

type ZSkipListNode struct {
	obj      any // 存储真正的值
	score    uint64
	backward *ZSkipListNode // 回退指针
	level    []ZSkipListNodeLevel
}

type ZSkipListNodeLevel struct {
	span    int            // 当前节点到下一个节点的距离
	forward *ZSkipListNode // 前进指针
}

var ZSKIPLIST_MAXLEVEL = 32

func NewZSkipList() *ZSkipList {
	zsl := &ZSkipList{
		level: 1,
	}

	zsl.head = new(ZSkipListNode)
	zsl.head.level = make([]ZSkipListNodeLevel, ZSKIPLIST_MAXLEVEL)
	return zsl
}

func CreateNode(level int, score uint64, obj any) *ZSkipListNode {
	zsln := &ZSkipListNode{
		obj:   obj,
		score: score,
	}
	zsln.level = make([]ZSkipListNodeLevel, level)
	return zsln
}

func (zl *ZSkipList) insert(obj any, score uint64) *ZSkipListNode {
	var (
		node   = zl.head
		updata = make([]*ZSkipListNode, ZSKIPLIST_MAXLEVEL) // 记录能到达的最右边的节点
		rank   = make([]int, ZSKIPLIST_MAXLEVEL)            // 记录每层的距离
	)

	// 确定新节点所处的位置和 span
	for i := zl.level - 1; i >= 0; i-- {
		if i != zl.level-1 {
			rank[i] = rank[i+1] // 沿用层次表上一层的 span
		}
		for node.level[i].forward != nil && node.level[i].forward.score < score {
			rank[i] += node.level[i].span
			node = node.level[i].forward
		}
		updata[i] = node
	}

	level := zl.RandomLevel()
	if level > zl.level {
		for i := zl.level; i < level; i++ {
			rank[i] = 0
			updata[i] = zl.head
			updata[i].level[i].span = int(zl.length)
		}
	}

	targetNode := CreateNode(level, score, obj)
	for i := 0; i < level; i++ {
		// 每层节点插入
		targetNode.level[i].forward = updata[i].level[i].forward
		updata[i].level[i].forward = targetNode

		targetNode.level[i].span = updata[i].level[i].span - (rank[0] - rank[i])
		updata[i].level[i].span = rank[0] - rank[i] + 1
	}

	// 更新上面层次节点
	for i := level; i < zl.level; i++ {
		updata[i].level[i].span++
	}

	// 设置回退指针
	if updata[0] != zl.head {
		targetNode.backward = updata[0]
	}

	if targetNode.level[0].forward != nil {
		targetNode.level[0].forward.backward = targetNode
	} else {
		zl.tail = targetNode
	}

	zl.length++
	return targetNode
}

func (zl *ZSkipList) RandomLevel() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(ZSKIPLIST_MAXLEVEL)
}

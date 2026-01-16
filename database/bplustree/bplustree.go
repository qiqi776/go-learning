package main

import (
	"fmt"
	"sort"
	"strings"
)

// Order (阶数/扇出)：一个节点最多包含的子节点数 (M)
// 为了演示分裂效果，我们将 M 设得很小。在 InnoDB 中这个值通常是 1000+
const M = 4

// BPNode B+树节点 (对应 InnoDB 的 Page)
type BPNode struct {
	IsLeaf bool // 是否是叶子节点

	// 共有部分
	Keys []int // 索引键 (有序)

	// 非叶子节点特有
	Children []*BPNode // 指向子节点的指针

	// 叶子节点特有
	Records []string // 实际数据 (模拟 Row Data)
	Next    *BPNode  // 指向下一个叶子节点的指针 (双向链表通常还有 Prev，这里简化为单向)
}

// BPTree B+树结构
type BPTree struct {
	Root *BPNode
}

// NewBPTree 初始化
func NewBPTree() *BPTree {
	return &BPTree{
		Root: &BPNode{IsLeaf: true},
	}
}

// ==========================================
// 1. 查找 (Search)
// ==========================================

func (t *BPTree) Get(key int) (string, bool) {
	curr := t.Root
	for !curr.IsLeaf {
		// 在内部节点找合适的子节点分支
		// 找到第一个大于 key 的位置 i，子节点就是 Children[i]
		idx := 0
		for idx < len(curr.Keys) && key >= curr.Keys[idx] {
			idx++
		}
		curr = curr.Children[idx]
	}

	// 在叶子节点内二分查找
	idx := sort.SearchInts(curr.Keys, key)
	if idx < len(curr.Keys) && curr.Keys[idx] == key {
		return curr.Records[idx], true
	}
	return "", false
}

// RangeScan 范围查询 (模拟 SELECT * FROM t WHERE id > X)
// 这是 B+ 树最强大的特性
func (t *BPTree) RangeScan(startKey int) []string {
	var results []string
	
	// 1. 先像普通查找一样找到起点叶子节点
	curr := t.Root
	for !curr.IsLeaf {
		idx := 0
		for idx < len(curr.Keys) && startKey >= curr.Keys[idx] {
			idx++
		}
		curr = curr.Children[idx]
	}

	// 2. 从起点开始，顺着 Next 指针遍历链表
	for curr != nil {
		for i, key := range curr.Keys {
			if key >= startKey {
				results = append(results, fmt.Sprintf("[%d:%s]", key, curr.Records[i]))
			}
		}
		curr = curr.Next
	}
	return results
}

// ==========================================
// 2. 插入 (Insert)
// ==========================================

func (t *BPTree) Insert(key int, value string) {
	// 递归插入，splitKey 和 newChild 是分裂时向上冒泡的数据
	splitKey, newChild := t.insertRecursive(t.Root, key, value)

	// 如果根节点分裂了，需要创建新的根
	if newChild != nil {
		newRoot := &BPNode{
			IsLeaf:   false,
			Keys:     []int{splitKey},
			Children: []*BPNode{t.Root, newChild},
		}
		t.Root = newRoot
	}
}

// 递归插入辅助函数
// 返回值: (向上提拔的Key, 新分裂出来的右兄弟节点)
func (t *BPTree) insertRecursive(node *BPNode, key int, value string) (int, *BPNode) {
	// --- A. 如果是叶子节点 (真正存数据) ---
	if node.IsLeaf {
		// 1. 找到插入位置
		idx := 0
		for idx < len(node.Keys) && key > node.Keys[idx] {
			idx++
		}

		// 2. 插入 Key 和 Value
		node.Keys = insertInt(node.Keys, idx, key)
		node.Records = insertString(node.Records, idx, value)

		// 3. 检查是否需要分裂 (Keys数量 >= M)
		if len(node.Keys) < M {
			return 0, nil
		}

		// --- 叶子节点分裂逻辑 ---
		// 裂变点: M/2
		splitIdx := M / 2
		
		// 创建右兄弟
		newLeaf := &BPNode{
			IsLeaf:  true,
			Keys:    append([]int(nil), node.Keys[splitIdx:]...),       // 后半部分 Key
			Records: append([]string(nil), node.Records[splitIdx:]...), // 后半部分 Data
			Next:    node.Next, // 维护链表: Old -> New -> Next
		}
		node.Next = newLeaf

		// 截断左节点 (保留前半部分)
		node.Keys = node.Keys[:splitIdx]
		node.Records = node.Records[:splitIdx]

		// 叶子分裂时，中间 Key (newLeaf的第一个) 需要复制一份提拔给父节点
		return newLeaf.Keys[0], newLeaf
	}

	// --- B. 如果是内部节点 (只存索引) ---
	// 1. 找到该去哪个子节点
	idx := 0
	for idx < len(node.Keys) && key >= node.Keys[idx] {
		idx++
	}

	// 2. 递归向下
	childUpKey, childNewNode := t.insertRecursive(node.Children[idx], key, value)

	// 3. 如果子节点没分裂，直接返回
	if childNewNode == nil {
		return 0, nil
	}

	// 4. 子节点分裂了，把提拔上来的 Key 插到当前节点
	node.Keys = insertInt(node.Keys, idx, childUpKey)
	node.Children = insertNodePtr(node.Children, idx+1, childNewNode)

	// 5. 检查当前内部节点是否需要分裂
	if len(node.Keys) < M {
		return 0, nil
	}

	// --- 内部节点分裂逻辑 ---
	// 注意：内部节点分裂与叶子不同，中间的 Key 会被“挤”上去，而不会保留在当前层
	splitIdx := M / 2
	upKey := node.Keys[splitIdx] // 这个 Key 还要继续往上提

	newInternal := &BPNode{
		IsLeaf:   false,
		Keys:     append([]int(nil), node.Keys[splitIdx+1:]...), // 注意: splitIdx 那个 key 被拿走了
		Children: append([]*BPNode(nil), node.Children[splitIdx+1:]...),
	}

	// 截断左节点
	node.Keys = node.Keys[:splitIdx]
	node.Children = node.Children[:splitIdx+1]

	return upKey, newInternal
}

// ==========================================
// 辅助工具函数
// ==========================================

func insertInt(s []int, idx int, val int) []int {
	s = append(s, 0)
	copy(s[idx+1:], s[idx:])
	s[idx] = val
	return s
}

func insertString(s []string, idx int, val string) []string {
	s = append(s, "")
	copy(s[idx+1:], s[idx:])
	s[idx] = val
	return s
}

func insertNodePtr(s []*BPNode, idx int, val *BPNode) []*BPNode {
	s = append(s, nil)
	copy(s[idx+1:], s[idx:])
	s[idx] = val
	return s
}

// Print 打印树结构 (可视化)
func (t *BPTree) Print() {
	fmt.Println("=== B+ Tree Structure ===")
	t.printRecursive(t.Root, 0)
	fmt.Println("=========================")
}

func (t *BPTree) printRecursive(node *BPNode, level int) {
	indent := strings.Repeat("    ", level)
	if node.IsLeaf {
		fmt.Printf("%s[Leaf] Keys: %v -> Next: %v\n", indent, node.Keys, node.Next != nil)
	} else {
		fmt.Printf("%s[Internal] Keys: %v\n", indent, node.Keys)
		for _, child := range node.Children {
			t.printRecursive(child, level+1)
		}
	}
}

func main() {
	bt := NewBPTree()

	// 模拟插入数据 (为了触发分裂，M设为了4)
	// 插入 1, 2, 3, 4 -> 触发第一次叶子分裂
	// 插入 5, 6, 7 -> 触发更多分裂
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	
	fmt.Println("开始插入数据...")
	for _, v := range data {
		bt.Insert(v, fmt.Sprintf("RowData-%d", v))
	}

	// 1. 打印树形结构
	bt.Print()

	// 2. 单点查询
	fmt.Println("\n查找 Key=5:")
	if val, ok := bt.Get(5); ok {
		fmt.Printf("Found: %s\n", val)
	} else {
		fmt.Println("Not Found")
	}

	// 3. 范围查询 (核心功能)
	fmt.Println("\n范围查找 > 8 (Range Scan):")
	rows := bt.RangeScan(8)
	for _, r := range rows {
		fmt.Print(r, " -> ")
	}
	fmt.Println("END")
}
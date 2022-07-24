package contains

import (
	"errors"
	"fmt"
	"sort"
)

type HuffTreeMap map[byte]*BitList

type HuffTree struct {
	Val         byte
	Weight      int
	Left, Right *HuffTree
}

// SearchCode 搜索指定字符的哈夫曼编码，填充到 code 里，如果找到则返回 true
func (t HuffTree) SearchCode(c byte, code *BitList) bool {
	if t.Val == c {
		return true
	}
	if t.Left == nil && t.Right == nil {
		return false
	}
	code.Push(0)
	ok := t.Left.SearchCode(c, code)
	if ok {
		return true
	}
	code.Pop()
	code.Push(1)
	ok = t.Right.SearchCode(c, code)
	if ok {
		return true
	}
	code.Pop()
	return false
}

// Encode 把字符串转化成哈夫曼编码
func (t HuffTree) Encode(str []byte) ([]byte, error) {
	res := NewBitList()
	m := HuffTreeMap{}
	for i := 0; i < len(str); i++ {
		bl, ok := m[str[i]]
		if !ok {
			bl = NewBitList()
			if !t.SearchCode(str[i], bl) {
				return nil, errors.New(fmt.Sprintf("%c is not exist in huffman tree", str[i]))
			}
			m[str[i]] = bl
		}
		for j := 0; j < bl.Size; j++ {
			res.Push(bl.Index(j))
		}
	}
	return res.DumpBytes(), nil
}

// Decode 把哈夫曼编码转化成字符串
func (t HuffTree) Decode(length int, b []byte) (res []byte) {
	tmp := &t
	code := NewBitListN(length, b)
	for i := 0; i < code.Size; i++ {
		if tmp.Left == nil && tmp.Right == nil {
			res = append(res, tmp.Val)
			tmp = &t
		}
		if code.Index(i) == 0 {
			tmp = tmp.Left
		} else {
			tmp = tmp.Right
		}
	}
	return
}

type huffQue struct {
	len  int
	elem []*HuffTree
}

// Len 实现 sort.interface 用于排序
func (q *huffQue) Len() int {
	return len(q.elem)
}

func (q *huffQue) Less(i, j int) bool {
	if q.elem[i] == nil {
		if q.elem[j] == nil {
			return true
		} else {
			return false
		}
	} else if q.elem[j] == nil {
		return true
	}
	return q.elem[i].Weight < q.elem[j].Weight
}

func (q *huffQue) Swap(i, j int) {
	q.elem[i], q.elem[j] = q.elem[j], q.elem[i]
}

func newHuffQue() *huffQue {
	return &huffQue{elem: make([]*HuffTree, 'z'+1)}
}

func newHuffNode(val byte) *HuffTree {
	return &HuffTree{Val: val}
}

func buildTree(que *huffQue) *HuffTree {
	for {
		if que.len == 2 {
			return &HuffTree{
				Val:    0,
				Weight: que.elem[0].Weight + que.elem[1].Weight,
				Left:   que.elem[0],
				Right:  que.elem[1],
			}
		}
		tmp := &HuffTree{
			Val:    0,
			Weight: que.elem[0].Weight + que.elem[1].Weight,
			Left:   que.elem[0],
			Right:  que.elem[1],
		}
		que.elem = append(que.elem[2:], tmp)
		sort.Sort(que)
		que.len -= 1
	}
}

func newHufTreeEmpty(tier int) *HuffTree {
	if tier == 0 {
		return nil
	}
	res := &HuffTree{}
	res.Left = newHufTreeEmpty(tier - 1)
	res.Right = newHufTreeEmpty(tier - 1)
	return res
}

// NewHuffTree 根据字符串新建哈夫曼树
func NewHuffTree(str []byte) (res *HuffTree, size int) {
	que := newHuffQue()
	for _, v := range str {
		if que.elem[v] == nil {
			que.elem[v] = newHuffNode(v)
			que.len += 1
		}
		que.elem[v].Weight += 1
	}
	size = que.len
	sort.Sort(que)
	que.elem = que.elem[:que.len]
	return buildTree(que), size
}

type huffTreeRaw struct {
	val byte
	bl  *BitList
}

// 通过剪枝构造哈夫曼树
func cropHuffTree(tree *HuffTree, raw *huffTreeRaw, tier int) {
	if tier == raw.bl.Size {
		tree.Left = nil
		tree.Right = nil
		tree.Val = raw.val
		return
	}
	if raw.bl.Index(tier) == 0 {
		cropHuffTree(tree.Left, raw, tier+1)
	} else {
		cropHuffTree(tree.Right, raw, tier+1)
	}
}

// UnmarshalHuffTree 根据字节流构造成哈夫曼树，字节流由 MarshalHuffTree 生成
func UnmarshalHuffTree(raw []byte) *HuffTree {
	tier := 0
	nodes := make([]*huffTreeRaw, 0)
	for i := 0; i < len(raw); {
		lenn := (int(raw[i+1]) + ByteBitLens - 1) / ByteBitLens
		if tier < int(raw[i+1]) {
			tier = int(raw[i+1])
		}
		bl := NewBitListN(int(raw[i+1]), raw[i+2:])
		nodes = append(nodes, &huffTreeRaw{val: raw[i], bl: bl})
		i += lenn + 2
	}
	emptyTree := newHufTreeEmpty(tier + 1)
	for _, node := range nodes {
		cropHuffTree(emptyTree, node, 0)
	}
	return emptyTree
}

type byteArrPtr struct {
	bytes []byte
}

// MarshalHuffTree 把哈夫曼树序列化成字节流
func MarshalHuffTree(t *HuffTree) []byte {
	res := &byteArrPtr{bytes: []byte{}}
	dfsMarshalTree(t, NewBitList(), res)
	return res.bytes
}

func dfsMarshalTree(root *HuffTree, code *BitList, res *byteArrPtr) {
	if root.Left == nil && root.Right == nil {
		res.bytes = append(res.bytes, root.Val)
		res.bytes = append(res.bytes, byte(code.Size))
		res.bytes = append(res.bytes, code.DumpBytes()...)
		return
	}
	code.Push(byte(0))
	dfsMarshalTree(root.Left, code, res)
	code.Pop()
	code.Push(byte(1))
	dfsMarshalTree(root.Right, code, res)
	code.Pop()
}

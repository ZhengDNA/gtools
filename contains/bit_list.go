package contains

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"unsafe"
)

const (
	Int64BitLens = int(unsafe.Sizeof(int64(1)) * 8)
	ByteBitLens  = int(unsafe.Sizeof(byte(1)) * 8)
)

const (
	Int64Byte1 = uint64(256-1) << (iota * 8)
	Int64Byte2
	Int64Byte3
	Int64Byte4
	Int64Byte5
	Int64Byte6
	Int64Byte7
	Int64Byte8
)

var int64Bytes = []uint64{Int64Byte8, Int64Byte7, Int64Byte6, Int64Byte5, Int64Byte4, Int64Byte3, Int64Byte2, Int64Byte1}
var int64BytesLen = len(int64Bytes)

type bit64 struct {
	len byte
	val uint64
}

type BitList struct {
	Size int
	arr  []bit64
}

// DumpBytes 把数据部分转成字节流
func (a BitList) DumpBytes() (res []byte) {
	i := 0
	for ; i < len(a.arr)-1; i++ {
		for j := 0; j < int64BytesLen; j++ {
			res = append(res, byte(a.arr[i].val&int64Bytes[j]>>((int64BytesLen-j-1)*8)&Int64Byte1))
		}
	}
	b, k := byte(0), 0
	for j := 0; j < int(a.arr[i].len); j, k = j+1, k+1 {
		b = b << 1
		b += a.Index(i*Int64BitLens + j)
		if k == 7 {
			k = -1
			res = append(res, b)
			b = 0
		}
	}
	if k != 0 {
		res = append(res, b<<(8-k))
	}
	return
}

// Push 在尾部添加一位数据
func (a *BitList) Push(b byte) {
	if !(b == 0 || b == 1) {
		return
	}
	if a.Size == 0 || int(a.arr[len(a.arr)-1].len) == Int64BitLens {
		a.arr = append(a.arr, bit64{})
	}
	ind := len(a.arr) - 1
	a.arr[ind].len += 1
	a.Size += 1
	a.arr[ind].val = a.arr[ind].val<<1 + uint64(b)
}

// Pop 移除并返回尾部的数据
func (a *BitList) Pop() (res byte) {
	res = a.Index(a.Size - 1)
	a.Size -= 1
	lena := len(a.arr) - 1
	a.arr[lena].len -= 1
	a.arr[lena].val = a.arr[lena].val >> 1 & (uint64(math.Pow(2, float64(a.arr[lena].len)) - 1))
	if a.arr[lena].len == 0 {
		a.arr = a.arr[:lena]
	}
	return
}

// SubSlice 返回 [start, end) 区间的新 BitList
func (a BitList) SubSlice(start, end int) *BitList {
	if end > a.Size || start > a.Size {
		panic(fmt.Sprintf("[%d, %d) out of range of length %d", start, end, a.Size))
	}
	res := NewBitList()
	for i := start; i < end; i++ {
		res.Push(a.Index(i))
	}
	return res
}

// Index 返回下标是 i 的那一位
func (a BitList) Index(i int) byte {
	if !(i < a.Size) {
		panic(errors.New(fmt.Sprintf("index out of range [%d] with length %d", i, a.Size)))
	}
	ind := i / Int64BitLens
	offset := int(a.arr[ind].len) - i%Int64BitLens - 1
	return byte(a.arr[ind].val & (1 << offset) >> offset)
}

func (a BitList) String() (res string) {
	for i := 0; i < len(a.arr); i++ {
		str := fmt.Sprintf("%b", a.arr[i].val)
		str = strings.Repeat("0", int(a.arr[i].len)-len(str)) + str
		res += str
	}
	return
}

// NewBitList 新建一个空的 BitList
func NewBitList() *BitList {
	return &BitList{arr: make([]bit64, 0)}
}

// NewBitListN 根据给定的长度和数据流来建立新的 BitList
func NewBitListN(length int, raw []byte) (res *BitList) {
	res = NewBitList()
	for i := 0; i < len(raw); i++ {
		for j := 0; j < ByteBitLens; j++ {
			if i*ByteBitLens+j == length {
				return
			}
			res.Push(byte(raw[i] & (1 << (ByteBitLens - j - 1)) >> (ByteBitLens - j - 1)))
		}
	}
	return
}

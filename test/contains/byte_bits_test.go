package test

import (
	"fmt"
	"gtools/contains"
	"math/rand"
	"strconv"
	"testing"
)

func TestBitsLen(t *testing.T) {
	fmt.Println(contains.Int64BitLens)
}

func TestBitList(t *testing.T) {
	a := contains.BitList{}

	// a.Push
	res := ""
	for i := 0; i < 100; i++ {
		b := byte(rand.Int() & 1)
		a.Push(b)
		res += strconv.FormatInt(int64(b), 10)
	}
	astr := a.String()
	fmt.Println(a)
	fmt.Println(res == astr)

	// a.Index
	res2 := ""
	for i := 0; i < 100; i++ {
		b := a.Index(i)
		res2 += strconv.FormatInt(int64(b), 10)
	}
	fmt.Println(res2)
	fmt.Println(res2 == astr)

	// a.DumpBytes
	a.DumpBytes()
}

func TestBitOperation(t *testing.T) {
	for i := 0; i < 1000_000; i++ {
		if i&-1 != i {
			fmt.Println("i", i)
		}
	}
	for i := 0; i < 10; i++ {
		fmt.Println(-2 & i)
	}
	var a int8 = int8(-1)
	fmt.Printf("%b\n", int64(a))
	fmt.Printf("%b\n%b\n%b\n%b\n", contains.Int64Byte1, contains.Int64Byte2, contains.Int64Byte3, contains.Int64Byte4)
}

func TestDump(t *testing.T) {
	n := contains.NewBitListN(8, []byte{121, 133})
	fmt.Println(n)
	fmt.Println(n.DumpBytes())
}

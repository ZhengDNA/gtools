package prints

import (
	"fmt"
	"gtools/prints"
	"testing"
)

func TestExpPrintln(t *testing.T) {
	prints.ExpPrintln("1", "12", "123", "1234", "12345", "123456", "1234567", "12345678", "123456789", "1234567890")
	fmt.Println("----------")
	prints.ExpPrintln(
		1+1,
		add(1, 2),
		multi(3, 4),
		add(add(1, 2), multi(3, 4)),
	)
}

func add(a, b int) int {
	return a + b
}

func multi(a, b int) int {
	return a * b
}

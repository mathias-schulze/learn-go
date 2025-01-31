package calc_test

import (
	"fmt"
	"learn/cli/calc"
	"testing"
)

func ExampleAdd() {
	fmt.Println(calc.Add(1, 2))
	// Output: 3
}

func TestAdd(t *testing.T) {
	tests := []struct {
		a    int
		b    int
		want int
	}{
		{1, 2, 3},
		{5, 2, 7},
		{10, 12, 22},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d+%d", tt.a, tt.b), func(t *testing.T) {
			got := calc.Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add(%d,%d)=%d: want: %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

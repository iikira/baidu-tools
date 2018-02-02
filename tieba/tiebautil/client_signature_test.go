package tiebautil

import (
	"fmt"
	"testing"
)

func TestSignature(t *testing.T) {
	post := map[string]string{
		"111": "222",
	}
	TiebaClientSignature(post)

	fmt.Println(post)
}

func BenchmarkSignature(b *testing.B) {
	post := map[string]string{
		"111": "222",
	}
	for i := 0; i < b.N; i++ {
		TiebaClientSignature(post)
	}
}

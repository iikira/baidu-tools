package ramdominfo

import (
	"fmt"
	"testing"
)

func TestRandomNumber(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(RamdomNumber(49000, 10000))
		fmt.Println(RamdomBytes(4))
		fmt.Println(RamdomMD5UpperString(4))
	}
	fmt.Println(SumIMEI("sdfadsfasd"))
	fmt.Println(GetPhoneModel("ddsasdddsasd'adlffasatrim"))
}

func BenchmarkRandomNumber(b *testing.B) {
	var i2 uint64
	for i := 0; i < b.N; i++ {
		i2 = RamdomNumber(49000, 10000)
		if i2 < 10000 && i2 > 49000 {
			fmt.Println(i2)
		}
	}
}

func BenchmarkRandomBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RamdomBytes(4)
	}
}

func BenchmarkSumIMEI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SumIMEI("sadfjasfdjk")
	}
}

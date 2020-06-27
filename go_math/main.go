package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// RandTest 随机
func RandTest() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		n := r.Intn(100)
		fmt.Println(n)
	}
}

// MathTest 方法测试
func MathTest() {
	i := -100
	fmt.Println(math.Pi)
	fmt.Println(math.Abs(float64(i))) //绝对值
	fmt.Println(math.Ceil(5.0))       //向上取整
	fmt.Println(math.Floor(5.8))      //向下取整
	fmt.Println(math.Mod(11, 3))      //取余数，同11%3
	fmt.Println(math.Modf(5.26))      //取整数，取小数
	fmt.Println(math.Pow(3, 2))       //x的y次方
	fmt.Println(math.Pow10(4))        // 10的n次方
	fmt.Println(math.Sqrt(8))         //开平方
	fmt.Println(math.Cbrt(8))         //开立方

}

func main() {
	MathTest()
}

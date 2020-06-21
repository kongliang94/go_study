package main

import "fmt"

// 定义Animal接口
type Animal interface {
	eat()
	attack()
}

// 定义一个Zoo动物园结构体
type Zoo struct {
	name string
}

// 动物园的一个喂养功能
func (zoo *Zoo) feet(animal Animal) {
	fmt.Print(zoo.name + " ")
	animal.eat()
}

// 定义狗Dog
type Dog struct {
	name string
}

// 定义狗Dog实现Animal的eat的方法
func (dog Dog) eat() {
	fmt.Println(dog.name + " eat food")
}

// 定义狗Dog实现Animal的attack的方法
func (dog Dog) attack() {
	fmt.Println(dog.name + " attack peer")
}

// 定义Lion
type Lion struct {
	name string
}

// 定义Lion实现Animal的eat的方法
func (lion Lion) eat() {
	fmt.Println(lion.name + " eat food")
}

// 定义Lion实现Animal的attack的方法
func (lion Lion) attack() {
	fmt.Println(lion.name + " attack peer")
}

func main() {

	// 实例化dog
	d := Dog{
		"kdog",
	}

	l := Lion{
		"klion",
	}

	// 实例化动物园Zoo
	zoo := Zoo{
		"北京动物园",
	}

	zoo1 := Zoo{
		"四川动物园",
	}

	zoo.feet(d)

	zoo.feet(l)

	zoo1.feet(l)
}

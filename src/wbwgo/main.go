package main

import (
	"fmt"
	//"reflect"
)

//"time"
//"wbwgo/game"

type AA interface {
	say()
	sayno()
}
type BB interface {
	say()
}

type Person struct {
}

func (p *Person) say() {
	fmt.Printf("hello")
}

func (p *Person) sayno() {
	fmt.Printf("no")
}

var b AA

func main() {
	b = &Person{}
	b.say()
}

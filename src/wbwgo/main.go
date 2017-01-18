package main

import (
	"fmt"
	"time"
	//"os/exec"
	//"reflect"
	"sync"
	"sync/atomic"
)

var globl_num int32 = 0

var map_nums map[int32]int32

var total_guard sync.RWMutex

func DoSome() {
	for i := 0; i < 10000; i++ {
		go func() {
			total_guard.Lock()
			defer total_guard.Unlock()
			num := atomic.AddInt32(&globl_num, 1)
			if _, ok := map_nums[num]; ok {
				fmt.Printf("-----------------------------num:%d\n", num)
			} else {
				map_nums[num] = num
				fmt.Printf("num:%d\n", num)
			}
		}()
	}
}

func main() {
	map_nums = make(map[int32]int32)

	DoSome()

	time.Sleep(10 * time.Second)
}

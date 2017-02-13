package common

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
)

func ConsoleStart() {
	go ConsoleRoutine()
}

func ConsoleRoutine() {
	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		command := string(data)
		switch command {
		case "cpus":
			fmt.Println(runtime.NumCPU(), " cpus and ", runtime.GOMAXPROCS(0), " in use")

		case "routines":
			fmt.Println("Current number of goroutines: ", runtime.NumGoroutine())

		case "startgc":
			runtime.GC()
			fmt.Println("gc finished")
		default:
			fmt.Println("Command error, try again.command:%s", command)
		}
	}
}

package main

import (
	initModule "github.com/NeptuneYeh/simplebank/init"
)

func main() {
	initProcess := initModule.NewMainInitProcess("./")
	initProcess.Run()
}

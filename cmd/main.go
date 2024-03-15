package main

import (
	initModule "github.com/NeptuneYeh/simplerecommend/init"
)

func main() {
	initProcess := initModule.NewMainInitProcess("./")
	initProcess.Run()
}

package main

import (
	"czlog"
	"time"
)

func main() {

	flag := czlog.LstdFlags | czlog.Lshortfile
	tstlog := czlog.New("czlogtst.log",1, flag, true, true)
	tstlog.SetLevel(czlog.LevelDebug)

	defer tstlog.ReleasePool()

	for i :=0; i<10000; i++ {
		tstlog.Debug("Counting: %v", i)
	}

	time.Sleep(1 * time.Second)

}

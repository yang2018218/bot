package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"wechatbot/internal/wechatbot"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	wechatbot.NewApp("wechatbot").Run()
}

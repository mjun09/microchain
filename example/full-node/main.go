package main

import (
	"flag"
	"os"

	"github.com/vgxbj/microchain/core"
)

var nodeIPOpt = flag.String("addr", "localhost", "ip address that node runs on")
var nodePortOpt = flag.Int("node_port", 3000, "port that node binds to")
var webPortOpt = flag.Int("web_port", 8000, "port that node binds to")

var l *core.Logger

func init() {
	l = core.InitLogger(os.Stdout)
	flag.Parse()
}

var initString = "                                 _                   \n          (_)                   | |         (_)      \n _ __ ___  _  ___ _ __ ___   ___| |__   __ _ _ _ __  \n| '_ ` _ \\| |/ __| '__/ _ \\ / __| '_ \\ / _` | | '_ \\ \n| | | | | | | (__| | | (_) | (__| | | | (_| | | | | |\n|_| |_| |_|_|\\___|_|  \\___/ \\___|_| |_|\\__,_|_|_| |_|\n"

func main() {
	c, err := newClient(*nodeIPOpt, *nodePortOpt, *webPortOpt, l)
	if err != nil {
		l.Error.Println(err)
		return
	}

	c.terminal <- initString

	c.repl()
}

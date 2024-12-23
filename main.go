package main

import (
	_ "cat-voting-app/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

func main() {
	initStorage()
	go startBot()

	fasthttp.ListenAndServe(":8081", func(ctx *fasthttp.RequestCtx) {
		fmt.Fprint(ctx, `Телеграм бот https://telegram.me/cur_rub_bot`)
	})
}

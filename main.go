package main

import (
	"Tiktok/router"
	"Tiktok/utils"
)

func main() {
	utils.Init()
	r := router.Router()
	r.Run(":7080")
}

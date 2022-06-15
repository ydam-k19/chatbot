package main

import "botserver/api"

func main() {
	api.ConnectDB()
	api.ConnectBot()
	api.Route()
}



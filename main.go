package main

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"abdanhafidz.com/go-boilerplate/router"
)

func main() {
	appProvider := provider.NewAppProvider()
	router.RunRouter(appProvider)
}

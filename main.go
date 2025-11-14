package main

import (
	"abdanhafidz.com/go-clean-layered-architecture/provider"
	"abdanhafidz.com/go-clean-layered-architecture/router"
)

func main() {
	appProvider := provider.NewAppProvider()
	router.RunRouter(appProvider)
}

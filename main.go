package main

import (
	"github.com/cellargalaxy/smzdm-reptile/controller"
	"github.com/cellargalaxy/smzdm-reptile/service"
)

func main() {
	go service.StartSearchService()
	controller.StartWebService()
}

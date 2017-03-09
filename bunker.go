package main

import (
	datasource "bunker/datasource"
	routes "bunker/routes"
)

func main() {
	routes.SetupRoutes()
	datasource.ConnectDatabase()
	defer datasource.DisconnectDatabase()
	routes.Listen()
}

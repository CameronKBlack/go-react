package main

import (
	"context"
	"go-react/backend/routes"
	"go-react/backend/services"
)

const uri = "mongodb://localhost:27017"

func main() {
	client, err := services.ConnectMongoDB(uri)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	routes.RouterSetup(client)	
}
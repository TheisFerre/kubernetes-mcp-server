package main

import (
	"fmt"

	kubernetes_client "github.com/TheisFerre/kubernetes-mcp-server/pkg/client"
	"github.com/TheisFerre/kubernetes-mcp-server/pkg/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	kubernetesClient, err := kubernetes_client.NewKubernetesClient()
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v\n", err)
		return
	}

	//Create a new MCP server
	s := server.NewMCPServer(
		"Kubernetes MCP Server",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	tools.InitializeTools(s, kubernetesClient)

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}

}

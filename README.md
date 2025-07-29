# Kubernetes MCP Server

A Model Context Protocol (MCP) server implementation for Kubernetes resource management. This project serves as both a practical tool for Kubernetes operations and a learning exercise in understanding the MCP protocol.

## 🎯 Overview

This MCP server enables AI assistants and LLMs to interact with Kubernetes clusters through a standardized protocol. It provides a comprehensive set of tools for managing common Kubernetes resources, making it easier to perform cluster operations through natural language interfaces.

## ✨ Features

### Supported Kubernetes Resources
- **Pods** - Manage individual workload units
- **Deployments** - Handle application deployments
- **StatefulSets** - Manage stateful applications
- **DaemonSets** - Control node-level services
- **ReplicaSets** - Manage pod replicas
- **Jobs** - Run batch workloads
- **CronJobs** - Schedule recurring tasks
- **Services** - Expose applications
- **ConfigMaps** - Manage configuration data
- **Secrets** - Handle sensitive information

### Available Operations
For each resource type, the following operations are supported:
- `get` - Retrieve resource details
- `list` - List resources in a namespace
- `create` - Create new resources
- `update` - Modify existing resources
- `delete` - Remove resources

## 🚀 Installation

### Prerequisites
- Go 1.24+ installed
- Access to a Kubernetes cluster
- `kubectl` configured with cluster access
- Valid kubeconfig file (typically at `~/.kube/config`)

### Building from Source
```bash
git clone https://github.com/TheisFerre/kubernetes-mcp-server.git
cd kubernetes-mcp-server
go build -o kubernetes-mcp-server .
```

## ⚙️ Configuration

### Setting up in VS Code

1. **Install the MCP extension** for VS Code if you haven't already.

2. **Configure the MCP server** by adding it to your VS Code MCP configuration file (`mcp.json`):

```json
{
  "servers": {
    "golang-k8s-mcp": {
      "command": "/path/to/your/kubernetes-mcp-server/kubernetes-mcp-server",
      "args": [],
      "env": {}
    }
  }
}
```

3. **Restart VS Code** to load the new MCP server.

4. **Verify the setup** by checking that the Kubernetes tools are available in your AI assistant interface.

### Kubernetes Cluster Access

The server automatically detects and uses your kubectl configuration:
- Default kubeconfig location: `~/.kube/config`
- Uses the current context configured in kubectl
- Inherits cluster access permissions from your kubeconfig

## 🛠️ Usage Examples

Once configured in VS Code, you can use natural language to interact with your Kubernetes cluster:

### Basic Operations
- "List all pods in the default namespace"
- "Get details for the nginx deployment in the web namespace"
- "Delete the old-job job from the batch namespace"

### Resource Management
- "Create a new deployment with the following spec: [JSON]"
- "Update the replica count for my-app deployment to 3"
- "Show me all services in the production namespace"

### Advanced Queries
- "List all failed pods across all namespaces"
- "Get the configuration for the database secret"
- "Show me the status of all cronjobs in the automation namespace"

## 🏗️ Architecture

The server is built with a modular architecture:

```
kubernetes-mcp-server/
├── main.go                 # Entry point and server initialization
├── pkg/
│   ├── client/            # Kubernetes client configuration
│   │   └── client.go      # kubeconfig and cluster connection
│   └── tools/             # MCP tool implementations
│       ├── tools.go       # Tool registration and routing
│       ├── pod.go         # Pod operations
│       ├── deployment.go  # Deployment operations
│       └── ...           # Other resource handlers
```

### Key Components

1. **MCP Server**: Built using the `mcp-go` library for protocol compliance
2. **Kubernetes Client**: Leverages the official `client-go` library for cluster communication
3. **Tool Registry**: Dynamic registration system for Kubernetes resource tools
4. **Resource Handlers**: Specialized handlers for each Kubernetes resource type

## 🔧 Development

### Adding New Resource Types

To add support for a new Kubernetes resource:

1. Add the resource constant to `pkg/tools/tools.go`
2. Create a new handler function following the existing pattern
3. Register the resource in the `InitializeTools` function
4. Add the switch case in the `addTool` function

### Local Development

```bash
# Run the server locally
go run main.go

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o kubernetes-mcp-server-linux .
GOOS=windows GOARCH=amd64 go build -o kubernetes-mcp-server.exe .
```

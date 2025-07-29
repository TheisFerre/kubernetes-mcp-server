package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"k8s.io/client-go/kubernetes"
)

const (
	pod         = "pod"
	deployment  = "deployment"
	statefulset = "statefulset"
	daemonset   = "daemonset"
	replicaset  = "replicaset"
	job         = "job"
	cronjob     = "cronjob"
	service     = "service"
	configmap   = "configmap"
	secret      = "secret"
)

func InitializeTools(server *server.MCPServer, kubernetesClient *kubernetes.Clientset) {

	for _, tool := range []string{
		pod,
		deployment,
		statefulset,
		daemonset,
		replicaset,
		job,
		cronjob,
		service,
		configmap,
		secret,
	} {
		addTool(context.Background(), server, registerTool(tool, kubernetesClient), kubernetesClient)

	}
}

func registerTool(tool string, kubernetesClient *kubernetes.Clientset) mcp.Tool {
	// Register the tool with the system (this is a placeholder for actual registration logic)
	// In a real implementation, this could involve adding the tool to a registry or initializing it
	// For now, we just print the tool name to simulate registration
	resourceTool := mcp.NewTool(tool,
		mcp.WithDescription("Tool for managing "+tool+" resources in Kubernetes"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the "+tool+" resource"),
		),
		mcp.WithString("namespace",
			mcp.Required(),
			mcp.Description("The namespace where the "+tool+" resource is located"),
		),
		mcp.WithString("action",
			mcp.Required(),
			mcp.Description("The action to perform on the "+tool+" resource (e.g., create, delete, update, get)"),
			mcp.Enum("create", "delete", "update", "get", "list"),
		),
		mcp.WithString("resourceSpec",
			mcp.Description("The specification for the "+tool+" resource in JSON format (optional, used for create/update actions)"),
		),
	)

	return resourceTool
}

func addTool(ctx context.Context, server *server.MCPServer, tool mcp.Tool, kubernetesClient *kubernetes.Clientset) {

	server.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Implement the logic to handle the tool request
		// This is a placeholder for actual tool handling logic
		var mcpResult *mcp.CallToolResult
		var err error
		name, err := request.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		namespace, err := request.RequireString("namespace")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		action, err := request.RequireString("action")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		resourceSpec := request.GetString("resourceSpec", "")

		switch tool.GetName() {
		case pod:
			mcpResult, err = podMCPResponse(ctx, name, namespace, action, resourceSpec, kubernetesClient.CoreV1().Pods(namespace))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
		case deployment:
			mcpResult, err = deploymentMCPResponse(ctx, name, namespace, action, resourceSpec, kubernetesClient.AppsV1().Deployments(namespace))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
		case statefulset:
			mcpResult, err = statefulsetMCPResponse(ctx, name, namespace, action, resourceSpec, kubernetesClient.AppsV1().StatefulSets(namespace))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
		case daemonset:
			mcpResult, err = daemonsetMCPResponse(ctx, name, namespace, action, resourceSpec, kubernetesClient.AppsV1().DaemonSets(namespace))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
		case replicaset:
			mcpResult, err = replicasetMCPResponse(ctx, name, namespace, action, resourceSpec, kubernetesClient.AppsV1().ReplicaSets(namespace))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
		case job:
			mcpResult, err = jobMCPResponse(ctx, name, namespace, action, resourceSpec, kubernetesClient.BatchV1().Jobs(namespace))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
		case cronjob:
			mcpResult, err = cronjobMCPResponse(ctx, name, namespace, action, resourceSpec, kubernetesClient.BatchV1().CronJobs(namespace))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
		case service:
			mcpResult, err = serviceMCPResponse(ctx, name, namespace, action, resourceSpec, kubernetesClient.CoreV1().Services(namespace))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
		case configmap:
			mcpResult, err = configmapMCPResponse(ctx, name, namespace, action, resourceSpec, kubernetesClient.CoreV1().ConfigMaps(namespace))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
		case secret:
			mcpResult, err = secretMCPResponse(ctx, name, namespace, action, resourceSpec, kubernetesClient.CoreV1().Secrets(namespace))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
		default:
			return mcp.NewToolResultError("Unknown tool: " + tool.GetName()), nil
		}
		return mcpResult, err
	})
}

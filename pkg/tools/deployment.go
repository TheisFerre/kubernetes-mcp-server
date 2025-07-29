package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

func deploymentMCPResponse(ctx context.Context, name string, namespace string, action string, resourceSpec string, deploymentInterface v1.DeploymentInterface) (*mcp.CallToolResult, error) {

	switch action {
	case "create":

		if resourceSpec == "" {
			return mcp.NewToolResultError("resourceSpec is required for create action"), nil
		}

		deploymentJson := []byte(resourceSpec)
		var deploymentSpec appsv1.Deployment
		if err := json.Unmarshal(deploymentJson, &deploymentSpec); err != nil {
			return mcp.NewToolResultError("Invalid resourceSpec JSON: " + err.Error()), nil
		}

		deployment, err := deploymentInterface.Create(
			ctx,
			&appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: deploymentSpec.Spec,
			},
			metav1.CreateOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		createdDeploymentSpec, err := json.Marshal(deployment)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal created deployment: " + err.Error()), nil
		}

		return mcp.NewToolResultStructuredOnly(createdDeploymentSpec), nil
	case "delete":
		err := deploymentInterface.Delete(
			ctx,
			name,
			metav1.DeleteOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Deleted deployment " + name + " in namespace " + namespace), nil
	case "update":
		deployment, err := deploymentInterface.Update(
			ctx,
			&appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
			},
			metav1.UpdateOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		updatedDeploymentSpec, err := json.Marshal(deployment)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal updated deployment: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(updatedDeploymentSpec), nil
	case "get":
		deployment, err := deploymentInterface.Get(
			ctx,
			name,
			metav1.GetOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		deploymentSpec, err := json.Marshal(deployment)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal deployment: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(deploymentSpec), nil
	case "list":
		deployments, err := deploymentInterface.List(
			ctx,
			metav1.ListOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var deploymentNames []string
		for _, deployment := range deployments.Items {
			deploymentNames = append(deploymentNames, deployment.Name)
		}
		return mcp.NewToolResultStructuredOnly(deploymentNames), nil
	}
	return mcp.NewToolResultError("Unknown action: " + action), nil
}

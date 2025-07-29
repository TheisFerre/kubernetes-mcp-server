package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func serviceMCPResponse(ctx context.Context, name string, namespace string, action string, resourceSpec string, serviceInterface v1.ServiceInterface) (*mcp.CallToolResult, error) {

	switch action {
	case "create":

		if resourceSpec == "" {
			return mcp.NewToolResultError("resourceSpec is required for create action"), nil
		}

		serviceJson := []byte(resourceSpec)
		var serviceSpec corev1.Service
		if err := json.Unmarshal(serviceJson, &serviceSpec); err != nil {
			return mcp.NewToolResultError("Invalid resourceSpec JSON: " + err.Error()), nil
		}

		service, err := serviceInterface.Create(
			ctx,
			&corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: serviceSpec.Spec,
			},
			metav1.CreateOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		createdServiceSpec, err := json.Marshal(service)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal created service: " + err.Error()), nil
		}

		return mcp.NewToolResultStructuredOnly(createdServiceSpec), nil
	case "delete":
		err := serviceInterface.Delete(
			ctx,
			name,
			metav1.DeleteOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Deleted service " + name + " in namespace " + namespace), nil
	case "update":
		service, err := serviceInterface.Update(
			ctx,
			&corev1.Service{
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
		updatedServiceSpec, err := json.Marshal(service)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal updated service: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(updatedServiceSpec), nil
	case "get":
		service, err := serviceInterface.Get(
			ctx,
			name,
			metav1.GetOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		serviceSpec, err := json.Marshal(service)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal service: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(serviceSpec), nil
	case "list":
		services, err := serviceInterface.List(
			ctx,
			metav1.ListOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var serviceNames []string
		for _, service := range services.Items {
			serviceNames = append(serviceNames, service.Name)
		}
		return mcp.NewToolResultStructuredOnly(serviceNames), nil
	}
	return mcp.NewToolResultError("Unknown action: " + action), nil
}

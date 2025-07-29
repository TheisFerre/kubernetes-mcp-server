package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func configmapMCPResponse(ctx context.Context, name string, namespace string, action string, resourceSpec string, configmapInterface v1.ConfigMapInterface) (*mcp.CallToolResult, error) {

	switch action {
	case "create":

		if resourceSpec == "" {
			return mcp.NewToolResultError("resourceSpec is required for create action"), nil
		}

		configmapJson := []byte(resourceSpec)
		var configmapSpec corev1.ConfigMap
		if err := json.Unmarshal(configmapJson, &configmapSpec); err != nil {
			return mcp.NewToolResultError("Invalid resourceSpec JSON: " + err.Error()), nil
		}

		configmap, err := configmapInterface.Create(
			ctx,
			&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Data:       configmapSpec.Data,
				BinaryData: configmapSpec.BinaryData,
			},
			metav1.CreateOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		createdConfigMapSpec, err := json.Marshal(configmap)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal created configmap: " + err.Error()), nil
		}

		return mcp.NewToolResultStructuredOnly(createdConfigMapSpec), nil
	case "delete":
		err := configmapInterface.Delete(
			ctx,
			name,
			metav1.DeleteOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Deleted configmap " + name + " in namespace " + namespace), nil
	case "update":
		configmap, err := configmapInterface.Update(
			ctx,
			&corev1.ConfigMap{
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
		updatedConfigMapSpec, err := json.Marshal(configmap)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal updated configmap: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(updatedConfigMapSpec), nil
	case "get":
		configmap, err := configmapInterface.Get(
			ctx,
			name,
			metav1.GetOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		configmapSpec, err := json.Marshal(configmap)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal configmap: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(configmapSpec), nil
	case "list":
		configmaps, err := configmapInterface.List(
			ctx,
			metav1.ListOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var configmapNames []string
		for _, configmap := range configmaps.Items {
			configmapNames = append(configmapNames, configmap.Name)
		}
		return mcp.NewToolResultStructuredOnly(configmapNames), nil
	}
	return mcp.NewToolResultError("Unknown action: " + action), nil
}

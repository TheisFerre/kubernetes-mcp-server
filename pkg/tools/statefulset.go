package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

func statefulsetMCPResponse(ctx context.Context, name string, namespace string, action string, resourceSpec string, statefulsetInterface v1.StatefulSetInterface) (*mcp.CallToolResult, error) {

	switch action {
	case "create":

		if resourceSpec == "" {
			return mcp.NewToolResultError("resourceSpec is required for create action"), nil
		}

		statefulsetJson := []byte(resourceSpec)
		var statefulsetSpec appsv1.StatefulSet
		if err := json.Unmarshal(statefulsetJson, &statefulsetSpec); err != nil {
			return mcp.NewToolResultError("Invalid resourceSpec JSON: " + err.Error()), nil
		}

		statefulset, err := statefulsetInterface.Create(
			ctx,
			&appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: statefulsetSpec.Spec,
			},
			metav1.CreateOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		createdStatefulSetSpec, err := json.Marshal(statefulset)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal created statefulset: " + err.Error()), nil
		}

		return mcp.NewToolResultStructuredOnly(createdStatefulSetSpec), nil
	case "delete":
		err := statefulsetInterface.Delete(
			ctx,
			name,
			metav1.DeleteOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Deleted statefulset " + name + " in namespace " + namespace), nil
	case "update":
		statefulset, err := statefulsetInterface.Update(
			ctx,
			&appsv1.StatefulSet{
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
		updatedStatefulSetSpec, err := json.Marshal(statefulset)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal updated statefulset: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(updatedStatefulSetSpec), nil
	case "get":
		statefulset, err := statefulsetInterface.Get(
			ctx,
			name,
			metav1.GetOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		statefulsetSpec, err := json.Marshal(statefulset)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal statefulset: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(statefulsetSpec), nil
	case "list":
		statefulsets, err := statefulsetInterface.List(
			ctx,
			metav1.ListOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var statefulsetNames []string
		for _, statefulset := range statefulsets.Items {
			statefulsetNames = append(statefulsetNames, statefulset.Name)
		}
		return mcp.NewToolResultStructuredOnly(statefulsetNames), nil
	}
	return mcp.NewToolResultError("Unknown action: " + action), nil
}

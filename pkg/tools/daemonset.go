package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

func daemonsetMCPResponse(ctx context.Context, name string, namespace string, action string, resourceSpec string, daemonsetInterface v1.DaemonSetInterface) (*mcp.CallToolResult, error) {

	switch action {
	case "create":

		if resourceSpec == "" {
			return mcp.NewToolResultError("resourceSpec is required for create action"), nil
		}

		daemonsetJson := []byte(resourceSpec)
		var daemonsetSpec appsv1.DaemonSet
		if err := json.Unmarshal(daemonsetJson, &daemonsetSpec); err != nil {
			return mcp.NewToolResultError("Invalid resourceSpec JSON: " + err.Error()), nil
		}

		daemonset, err := daemonsetInterface.Create(
			ctx,
			&appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: daemonsetSpec.Spec,
			},
			metav1.CreateOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		createdDaemonSetSpec, err := json.Marshal(daemonset)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal created daemonset: " + err.Error()), nil
		}

		return mcp.NewToolResultStructuredOnly(createdDaemonSetSpec), nil
	case "delete":
		err := daemonsetInterface.Delete(
			ctx,
			name,
			metav1.DeleteOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Deleted daemonset " + name + " in namespace " + namespace), nil
	case "update":
		daemonset, err := daemonsetInterface.Update(
			ctx,
			&appsv1.DaemonSet{
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
		updatedDaemonSetSpec, err := json.Marshal(daemonset)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal updated daemonset: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(updatedDaemonSetSpec), nil
	case "get":
		daemonset, err := daemonsetInterface.Get(
			ctx,
			name,
			metav1.GetOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		daemonsetSpec, err := json.Marshal(daemonset)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal daemonset: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(daemonsetSpec), nil
	case "list":
		daemonsets, err := daemonsetInterface.List(
			ctx,
			metav1.ListOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var daemonsetNames []string
		for _, daemonset := range daemonsets.Items {
			daemonsetNames = append(daemonsetNames, daemonset.Name)
		}
		return mcp.NewToolResultStructuredOnly(daemonsetNames), nil
	}
	return mcp.NewToolResultError("Unknown action: " + action), nil
}

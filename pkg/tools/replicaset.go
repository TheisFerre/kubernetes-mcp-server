package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

func replicasetMCPResponse(ctx context.Context, name string, namespace string, action string, resourceSpec string, replicasetInterface v1.ReplicaSetInterface) (*mcp.CallToolResult, error) {

	switch action {
	case "create":

		if resourceSpec == "" {
			return mcp.NewToolResultError("resourceSpec is required for create action"), nil
		}

		replicasetJson := []byte(resourceSpec)
		var replicasetSpec appsv1.ReplicaSet
		if err := json.Unmarshal(replicasetJson, &replicasetSpec); err != nil {
			return mcp.NewToolResultError("Invalid resourceSpec JSON: " + err.Error()), nil
		}

		replicaset, err := replicasetInterface.Create(
			ctx,
			&appsv1.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: replicasetSpec.Spec,
			},
			metav1.CreateOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		createdReplicaSetSpec, err := json.Marshal(replicaset)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal created replicaset: " + err.Error()), nil
		}

		return mcp.NewToolResultStructuredOnly(createdReplicaSetSpec), nil
	case "delete":
		err := replicasetInterface.Delete(
			ctx,
			name,
			metav1.DeleteOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Deleted replicaset " + name + " in namespace " + namespace), nil
	case "update":
		replicaset, err := replicasetInterface.Update(
			ctx,
			&appsv1.ReplicaSet{
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
		updatedReplicaSetSpec, err := json.Marshal(replicaset)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal updated replicaset: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(updatedReplicaSetSpec), nil
	case "get":
		replicaset, err := replicasetInterface.Get(
			ctx,
			name,
			metav1.GetOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		replicasetSpec, err := json.Marshal(replicaset)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal replicaset: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(replicasetSpec), nil
	case "list":
		replicasets, err := replicasetInterface.List(
			ctx,
			metav1.ListOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var replicasetNames []string
		for _, replicaset := range replicasets.Items {
			replicasetNames = append(replicasetNames, replicaset.Name)
		}
		return mcp.NewToolResultStructuredOnly(replicasetNames), nil
	}
	return mcp.NewToolResultError("Unknown action: " + action), nil
}

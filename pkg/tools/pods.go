package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func podMCPResponse(ctx context.Context, name string, namespace string, action string, resourceSpec string, podInterface v1.PodInterface) (*mcp.CallToolResult, error) {

	switch action {
	case "create":

		if resourceSpec == "" {
			return mcp.NewToolResultError("resourceSpec is required for create action"), nil
		}

		podJson := []byte(resourceSpec)
		var podSpec corev1.Pod
		if err := json.Unmarshal(podJson, &podSpec); err != nil {
			return mcp.NewToolResultError("Invalid resourceSpec JSON: " + err.Error()), nil
		}

		pod, err := podInterface.Create(
			ctx,
			&corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: podSpec.Spec,
			},
			metav1.CreateOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		createdPodSpec, err := json.Marshal(pod)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal created pod: " + err.Error()), nil
		}

		return mcp.NewToolResultStructuredOnly(createdPodSpec), nil
		//return mcp.NewToolResultText("Created pod " + pod.Name + " in namespace " + namespace), nil
	case "delete":
		err := podInterface.Delete(
			ctx,
			name,
			metav1.DeleteOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Deleted pod " + name + " in namespace " + namespace), nil
	case "update":
		pod, err := podInterface.Update(
			ctx,
			&corev1.Pod{
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
		updatedPodSpec, err := json.Marshal(pod)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal updated pod: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(updatedPodSpec), nil
	case "get":
		pod, err := podInterface.Get(
			ctx,
			name,
			metav1.GetOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		podSpec, err := json.Marshal(pod)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal pod: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(podSpec), nil
	case "list":
		pods, err := podInterface.List(
			ctx,
			metav1.ListOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var podNames []string
		for _, pod := range pods.Items {
			podNames = append(podNames, pod.Name)
		}
		return mcp.NewToolResultStructuredOnly(podNames), nil
	}
	return mcp.NewToolResultError("Unknown action: " + action), nil
}

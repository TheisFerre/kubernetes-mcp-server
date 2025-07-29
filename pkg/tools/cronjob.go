package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/batch/v1"
)

func cronjobMCPResponse(ctx context.Context, name string, namespace string, action string, resourceSpec string, cronjobInterface v1.CronJobInterface) (*mcp.CallToolResult, error) {

	switch action {
	case "create":

		if resourceSpec == "" {
			return mcp.NewToolResultError("resourceSpec is required for create action"), nil
		}

		cronjobJson := []byte(resourceSpec)
		var cronjobSpec batchv1.CronJob
		if err := json.Unmarshal(cronjobJson, &cronjobSpec); err != nil {
			return mcp.NewToolResultError("Invalid resourceSpec JSON: " + err.Error()), nil
		}

		cronjob, err := cronjobInterface.Create(
			ctx,
			&batchv1.CronJob{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: cronjobSpec.Spec,
			},
			metav1.CreateOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		createdCronJobSpec, err := json.Marshal(cronjob)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal created cronjob: " + err.Error()), nil
		}

		return mcp.NewToolResultStructuredOnly(createdCronJobSpec), nil
	case "delete":
		err := cronjobInterface.Delete(
			ctx,
			name,
			metav1.DeleteOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Deleted cronjob " + name + " in namespace " + namespace), nil
	case "update":
		cronjob, err := cronjobInterface.Update(
			ctx,
			&batchv1.CronJob{
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
		updatedCronJobSpec, err := json.Marshal(cronjob)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal updated cronjob: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(updatedCronJobSpec), nil
	case "get":
		cronjob, err := cronjobInterface.Get(
			ctx,
			name,
			metav1.GetOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		cronjobSpec, err := json.Marshal(cronjob)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal cronjob: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(cronjobSpec), nil
	case "list":
		cronjobs, err := cronjobInterface.List(
			ctx,
			metav1.ListOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var cronjobNames []string
		for _, cronjob := range cronjobs.Items {
			cronjobNames = append(cronjobNames, cronjob.Name)
		}
		return mcp.NewToolResultStructuredOnly(cronjobNames), nil
	}
	return mcp.NewToolResultError("Unknown action: " + action), nil
}

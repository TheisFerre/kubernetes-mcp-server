package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/batch/v1"
)

func jobMCPResponse(ctx context.Context, name string, namespace string, action string, resourceSpec string, jobInterface v1.JobInterface) (*mcp.CallToolResult, error) {

	switch action {
	case "create":

		if resourceSpec == "" {
			return mcp.NewToolResultError("resourceSpec is required for create action"), nil
		}

		jobJson := []byte(resourceSpec)
		var jobSpec batchv1.Job
		if err := json.Unmarshal(jobJson, &jobSpec); err != nil {
			return mcp.NewToolResultError("Invalid resourceSpec JSON: " + err.Error()), nil
		}

		job, err := jobInterface.Create(
			ctx,
			&batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: jobSpec.Spec,
			},
			metav1.CreateOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		createdJobSpec, err := json.Marshal(job)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal created job: " + err.Error()), nil
		}

		return mcp.NewToolResultStructuredOnly(createdJobSpec), nil
	case "delete":
		err := jobInterface.Delete(
			ctx,
			name,
			metav1.DeleteOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Deleted job " + name + " in namespace " + namespace), nil
	case "update":
		job, err := jobInterface.Update(
			ctx,
			&batchv1.Job{
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
		updatedJobSpec, err := json.Marshal(job)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal updated job: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(updatedJobSpec), nil
	case "get":
		job, err := jobInterface.Get(
			ctx,
			name,
			metav1.GetOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		jobSpec, err := json.Marshal(job)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal job: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(jobSpec), nil
	case "list":
		jobs, err := jobInterface.List(
			ctx,
			metav1.ListOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var jobNames []string
		for _, job := range jobs.Items {
			jobNames = append(jobNames, job.Name)
		}
		return mcp.NewToolResultStructuredOnly(jobNames), nil
	}
	return mcp.NewToolResultError("Unknown action: " + action), nil
}

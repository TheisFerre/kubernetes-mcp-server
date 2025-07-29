package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func secretMCPResponse(ctx context.Context, name string, namespace string, action string, resourceSpec string, secretInterface v1.SecretInterface) (*mcp.CallToolResult, error) {

	switch action {
	case "create":

		if resourceSpec == "" {
			return mcp.NewToolResultError("resourceSpec is required for create action"), nil
		}

		secretJson := []byte(resourceSpec)
		var secretSpec corev1.Secret
		if err := json.Unmarshal(secretJson, &secretSpec); err != nil {
			return mcp.NewToolResultError("Invalid resourceSpec JSON: " + err.Error()), nil
		}

		secret, err := secretInterface.Create(
			ctx,
			&corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Type:       secretSpec.Type,
				Data:       secretSpec.Data,
				StringData: secretSpec.StringData,
			},
			metav1.CreateOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		createdSecretSpec, err := json.Marshal(secret)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal created secret: " + err.Error()), nil
		}

		return mcp.NewToolResultStructuredOnly(createdSecretSpec), nil
	case "delete":
		err := secretInterface.Delete(
			ctx,
			name,
			metav1.DeleteOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Deleted secret " + name + " in namespace " + namespace), nil
	case "update":
		secret, err := secretInterface.Update(
			ctx,
			&corev1.Secret{
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
		updatedSecretSpec, err := json.Marshal(secret)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal updated secret: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(updatedSecretSpec), nil
	case "get":
		secret, err := secretInterface.Get(
			ctx,
			name,
			metav1.GetOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		secretSpec, err := json.Marshal(secret)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal secret: " + err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(secretSpec), nil
	case "list":
		secrets, err := secretInterface.List(
			ctx,
			metav1.ListOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var secretNames []string
		for _, secret := range secrets.Items {
			secretNames = append(secretNames, secret.Name)
		}
		return mcp.NewToolResultStructuredOnly(secretNames), nil
	}
	return mcp.NewToolResultError("Unknown action: " + action), nil
}

package main

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"fmt"
	"log"
)

func getSecret(projectNum string, secretName string, versionNum string) *secretmanagerpb.SecretPayload {
	// Client
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}
	defer client.Close()

	// build request to access secret
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectNum, secretName, versionNum),
	}

	// request secret
	secret, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("error getting secret: %v", err)
	}

	return secret.Payload
}

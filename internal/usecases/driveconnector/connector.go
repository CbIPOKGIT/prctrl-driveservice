package driveconnector

import (
	"context"
	"os"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func NewService() (*drive.Service, error) {
	return drive.NewService(context.Background(), option.WithCredentialsJSON([]byte(os.Getenv("SERVICE_TOKEN"))))
}

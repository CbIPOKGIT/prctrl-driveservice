package driveconnector

import (
	"context"
	"os"
	"strings"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func NewService() (*drive.Service, error) {
	return drive.NewService(context.Background(), option.WithCredentialsJSON([]byte(os.Getenv("SERVICE_TOKEN"))))
}

func getMimeTypeFromExtension(filename string) string {
	switch {
	case strings.Contains(filename, ".xlsx"):
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	default:
		return ""
	}
}

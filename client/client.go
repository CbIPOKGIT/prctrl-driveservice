package client

import (
	"context"

	"github.com/CbIPOKGIT/prctrl-driveservice/driveservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type DriveClient struct {
	host  string
	token string
}

type FileInfo struct {
	Name       string `json:"filename"`
	Id         string `json:"fileid"`
	IsDir      bool   `json:"is_dir"`
	ParentName string `json:"parent_name"`
	ParentId   string `json:"parent_id"`
	Size       uint64 `json:"size"`
}

// NewClient - створення клієнта для взаємодії з сервісом.
//
// host - адреса сервісу.
// token - статичний токен для авторизації на сервісі.
func NewClient(host, token string) *DriveClient {
	return &DriveClient{
		host:  host,
		token: token,
	}
}

// Upload - завантаження файлу на g-диск.
func (dc *DriveClient) Upload(name string, content []byte, parent ...*FileInfo) (*FileInfo, error) {
	client, err := dc.getClient()
	if err != nil {
		return nil, err
	}

	request := &driveservice.UploadRequest{
		Name:    name,
		Content: content,
		Share:   true,
	}

	if len(parent) > 0 {
		request.Parent = &driveservice.FileInfo{
			Name:       parent[0].Name,
			Id:         parent[0].Id,
			IsDir:      true,
			ParentId:   parent[0].ParentId,
			ParentName: parent[0].ParentName,
			Size:       parent[0].Size,
		}
	}

	res, err := client.Upload(dc.getContext(), request)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Name:       res.GetName(),
		Id:         res.GetId(),
		IsDir:      res.GetIsDir(),
		ParentName: res.GetParentName(),
		ParentId:   res.GetParentId(),
		Size:       res.GetSize(),
	}, nil
}

// FilesList - отримання списку файлів по заданим параметрам.
func (dc *DriveClient) FilesList(filter *FileInfo) ([]*FileInfo, error) {
	client, err := dc.getClient()
	if err != nil {
		return nil, err
	}

	request := &driveservice.FileInfo{
		Name:       filter.Name,
		Id:         filter.Id,
		IsDir:      filter.IsDir,
		ParentId:   filter.ParentId,
		ParentName: filter.ParentName,
		Size:       filter.Size,
	}

	res, err := client.FilesList(dc.getContext(), request)
	if err != nil {
		return nil, err
	}

	files := make([]*FileInfo, 0, len(res.GetFiles()))
	for _, file := range res.GetFiles() {
		files = append(files, &FileInfo{
			Name:       file.GetName(),
			Id:         file.GetId(),
			IsDir:      file.GetIsDir(),
			ParentName: file.GetParentName(),
			ParentId:   file.GetParentId(),
			Size:       file.GetSize(),
		})
	}

	return files, nil
}

func (dc *DriveClient) Download(filter *FileInfo) ([]byte, error) {
	client, err := dc.getClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Download(dc.getContext(), &driveservice.FileInfo{
		Name:     filter.Name,
		Id:       filter.Id,
		IsDir:    filter.IsDir,
		ParentId: filter.ParentId,
		Size:     filter.Size,
	})

	if err != nil {
		return nil, err
	}

	return response.Content, nil
}

func (dc *DriveClient) getClient() (driveservice.DriveServiceClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(dc.host, opts...)
	if err != nil {
		return nil, err
	}

	return driveservice.NewDriveServiceClient(conn), nil
}

func (dc *DriveClient) getContext() context.Context {
	metaToken := map[string]string{"token": dc.token}
	return metadata.NewOutgoingContext(context.Background(), metadata.New(metaToken))
}

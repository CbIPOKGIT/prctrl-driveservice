package server

import (
	"context"
	"errors"

	"github.com/CbIPOKGIT/prctrl-driveservice/driveservice"
	"github.com/CbIPOKGIT/prctrl-driveservice/internal/entity"
	"github.com/CbIPOKGIT/prctrl-driveservice/internal/usecases/dataconverter"
	"github.com/CbIPOKGIT/prctrl-driveservice/internal/usecases/driveconnector"
)

func (ds *DriveService) FilesList(ctx context.Context, req *driveservice.FileInfo) (*driveservice.FilesInfo, error) {
	filter := dataconverter.InfoFromRequest(req)
	if files, err := driveconnector.FindEntities(filter); err == nil {
		return dataconverter.DriveFilesToList(files), nil
	} else {
		return nil, err
	}
}

// Upload - завантаження файлу на g-диск.
//
// Обов'язкові поля - ім'я файлу (name) та його вміст (content).
// Опціонально можна вказати батьківську папку (parent) та чи потрібно відкрити доступ до файлу (share).
// Батьківська папка буде автоматично створена, якщо її не існує.
func (ds *DriveService) Upload(ctx context.Context, req *driveservice.UploadRequest) (*driveservice.FileInfo, error) {
	options := &driveconnector.UploadOptions{
		Share: req.Share,
	}

	if req.Parent != nil {
		options.Parent = dataconverter.InfoFromRequest(req.Parent)
	}

	file, err := driveconnector.UploadFile(req.Name, req.Content, options)
	return dataconverter.DriveFileToInfo(file), err
}

// Download - завантаження файлу з g-диска.
//
// FileInfo - набір фільтрів для пошуку файлу. Повертаємо перший файл, який відповідає умовам.
func (ds *DriveService) Download(ctx context.Context, req *driveservice.FileInfo) (*driveservice.FileContent, error) {
	filter := dataconverter.InfoFromRequest(req)

	files, err := driveconnector.FindEntities(filter)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, errors.New(entity.ERROR_FILE_NOT_FOUND)
	}

	content, err := driveconnector.GetFileContent(files[0].Id)
	if err != nil {
		return nil, err
	}

	return &driveservice.FileContent{Content: content}, nil
}

package dataconverter

import (
	"github.com/CbIPOKGIT/prctrl-driveservice/driveservice"
	"github.com/CbIPOKGIT/prctrl-driveservice/internal/entity"
	"google.golang.org/api/drive/v3"
)

func DriveFileToInfo(file *drive.File) *driveservice.FileInfo {
	var parentID string

	if len(file.Parents) > 0 {
		parentID = file.Parents[0]
	}
	return &driveservice.FileInfo{
		Id:       file.Id,
		Name:     file.Name,
		Size:     uint64(file.Size),
		IsDir:    file.MimeType == entity.MIME_TYPE_FOLDER,
		ParentId: parentID,
		Created:  file.CreatedTime,
		Modified: file.ModifiedTime,
	}
}

func DriveFilesToList(files []*drive.File) *driveservice.FilesInfo {
	list := make([]*driveservice.FileInfo, len(files))

	for i, file := range files {
		list[i] = DriveFileToInfo(file)
	}

	return &driveservice.FilesInfo{Files: list}
}

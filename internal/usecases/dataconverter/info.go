package dataconverter

import (
	"github.com/CbIPOKGIT/prctrl-driveservice/driveservice"
	"github.com/CbIPOKGIT/prctrl-driveservice/internal/entity"
)

func InfoFromRequest(info *driveservice.FileInfo) *entity.DriveEntityFileInfo {
	return &entity.DriveEntityFileInfo{
		Name:       info.GetName(),
		Fileid:     info.GetId(),
		IsDir:      info.GetIsDir(),
		ParentName: info.GetParentName(),
		ParentId:   info.GetParentId(),
		Size:       info.GetSize(),
	}
}

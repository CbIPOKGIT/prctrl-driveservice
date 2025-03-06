package driveconnector

import (
	"errors"

	"github.com/CbIPOKGIT/prctrl-driveservice/internal/entity"
)

func DeleteFileFromDrive(filter *entity.DriveEntityFileInfo) (bool, error) {
	if filter.IsEmptyIdent() {
		return false, errors.New(entity.ERROR_FILTERS_ARE_EMPTY)
	}

	files, err := FindEntities(filter)
	if err != nil {
		return false, err
	}

	if len(files) == 0 {
		return false, nil
	}

	service, err := NewService()
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if err := service.Files.Delete(file.Id).Do(); err != nil {
			return false, err
		}
	}

	return true, nil
}

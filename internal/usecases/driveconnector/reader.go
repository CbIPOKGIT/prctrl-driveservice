package driveconnector

import (
	"errors"
	"fmt"
	"io"

	"github.com/CbIPOKGIT/prctrl-driveservice/internal/entity"
	"google.golang.org/api/drive/v3"
)

func FindEntities(filter *entity.DriveEntityFileInfo) ([]*drive.File, error) {
	if filter.IsEmpty() {
		return nil, errors.New(entity.ERROR_FILTERS_ARE_EMPTY)
	}

	service, err := NewService()
	if err != nil {
		return nil, err
	}

	if filter.NeedQueryParent() {
		if parent, err := GetFolderByName(filter.ParentName); err == nil {
			filter.ParentId = parent.Id
		} else {
			return nil, err
		}
	}

	if filter.Fileid != "" {
		if file, err := service.Files.Get(filter.Fileid).Do(); err == nil {
			return []*drive.File{file}, nil
		} else {
			return nil, err
		}
	}

	if files, err := service.Files.List().Q(filter.ToQuery()).Do(); err == nil {
		return files.Files, nil
	} else {
		return nil, err
	}
}

func GetFileContent(id string) ([]byte, error) {
	service, err := NewService()
	if err != nil {
		return nil, err
	}

	res, err := service.Files.Get(id).Download()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

func GetFolderByName(name string) (*drive.File, error) {
	service, err := NewService()
	if err != nil {
		return nil, err
	}

	files, err := service.Files.List().
		Q(fmt.Sprintf("mimeType = '%s' and name = '%s'", entity.MIME_TYPE_FOLDER, name)).
		Do()

	if err != nil {
		return nil, err
	}

	if len(files.Files) == 0 {
		return nil, errors.New(entity.ERROR_FILE_NOT_FOUND)
	}

	return files.Files[0], nil
}

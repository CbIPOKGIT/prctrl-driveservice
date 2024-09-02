package driveconnector

import (
	"bytes"
	"errors"

	"github.com/CbIPOKGIT/prctrl-driveservice/internal/entity"
	"google.golang.org/api/drive/v3"
)

// UploadOptions - додаткові параметри завантаження файлу.
//
// Parent - батьківська папка. Якщо пусто - файл завантажується в корінь.
// Share - чи потрібно відкрити доступ до файлу.
type UploadOptions struct {
	Parent *entity.DriveEntityFileInfo
	Share  bool
}

// UploadFile - завантажуємо файл на Google Drive.
//
// filename - назва файлу,
// content - вміст файлу.
// options - параметри завантаження.
func UploadFile(filename string, content []byte, options *UploadOptions) (*drive.File, error) {
	if filename == "" {
		return nil, errors.New(entity.ERROR_FILENAME_IS_EMPTY)
	}
	parents := make([]string, 0)

	if options.Parent != nil {

		// На всякий випадок вказуємо в parent що це папка.
		options.Parent.IsDir = true

		// Шукаємо або створюємо папку.
		entity, err := FindOrCreateFolder(options.Parent)
		if err != nil {
			return nil, err
		}
		parents = []string{entity.Id}
	}

	service, err := NewService()
	if err != nil {
		return nil, err
	}

	file := &drive.File{
		Name:    filename,
		Parents: parents,
	}

	return service.Files.Create(file).Media(bytes.NewReader(content)).Do()
}

// CreateFolder - створюємо папку на Google Drive.
func CreateFolder(data *entity.DriveEntityFileInfo) (*drive.File, error) {
	service, err := NewService()
	if err != nil {
		return nil, err
	}

	file := &drive.File{
		Name:     data.Name,
		MimeType: entity.MIME_TYPE_FOLDER,
	}

	return service.Files.Create(file).Do()
}

// FindOrCreateFolder - шукаємо або створюємо папку на Google Drive.
func FindOrCreateFolder(data *entity.DriveEntityFileInfo) (*drive.File, error) {
	files, err := FindEntities(data)
	if err != nil {
		return nil, err
	}

	if len(files) > 0 {
		return files[0], nil
	}

	return CreateFolder(data)
}

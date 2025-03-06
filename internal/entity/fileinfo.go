package entity

import (
	"fmt"
	"strings"
)

// DriveEntityFileInfo - інформація про файл.
// Використовується для фільтрації файлів та папок при запитах до Google Drive,
// або для збереження інформації про файл.
type DriveEntityFileInfo struct {
	Name       string `json:"name" form:"name"`
	Fileid     string `json:"file-id" form:"file-id"`
	IsDir      bool   `json:"is-dir" form:"is-dir"`
	ParentName string `json:"parent-name" form:"parent-name"`
	ParentId   string `json:"parent-id" form:"parent-id"`
	Size       uint64 `json:"size" form:"size"`
	Created    string `json:"created" form:"created"`
	Modified   string `json:"modified" form:"modified"`
}

// IsEmpty - перевіряє чи є фільтр порожнім.
// Повинно бути заповнене хоча б одне поле з ID чи іменем папки/файлу.
func (f *DriveEntityFileInfo) IsEmpty() bool {
	return f.IsEmptyIdent() && f.ParentName == "" && f.ParentId == ""
}

// IsEmptyIdent - перевіряє чи порожні дані файлу у фільтрі
func (f *DriveEntityFileInfo) IsEmptyIdent() bool {
	return f.Name == "" && f.Fileid == ""
}

// ToQuery - конвертує дані в рядок запиту для пошуку файлів та папок.
func (f *DriveEntityFileInfo) ToQuery() string {
	conditions := make([]string, 0, 6)

	if f.Name != "" {
		conditions = append(conditions, fmt.Sprintf(`name = %q`, f.Name))
	}

	if f.Fileid != "" {
		// conditions = append(conditions, fmt.Sprintf(`appProperties has { key='Id' and value='%s' }`, f.Fileid))
		conditions = append(conditions, fmt.Sprintf(`fileId = '%s'`, f.Fileid))
	}

	if f.IsDir {
		conditions = append(conditions, fmt.Sprintf("mimeType = '%s'", MIME_TYPE_FOLDER))
	}

	if f.ParentId != "" {
		conditions = append(conditions, fmt.Sprintf("'%s' in parents", f.ParentId))
	}

	return strings.Join(conditions, " and ")
}

func (f *DriveEntityFileInfo) NeedQueryParent() bool {
	return f.ParentName != "" && f.ParentId == ""
}

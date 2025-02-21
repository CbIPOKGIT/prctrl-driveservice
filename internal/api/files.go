package api

import (
	"io"
	"mime/multipart"
	"strings"

	"github.com/CbIPOKGIT/prctrl-driveservice/internal/entity"
	"github.com/CbIPOKGIT/prctrl-driveservice/internal/usecases/driveconnector"
	"github.com/gin-gonic/gin"
)

func LoadFile(c *gin.Context) {
	filters := &entity.DriveEntityFileInfo{}
	if err := c.BindQuery(filters); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if filters.IsEmpty() {
		c.JSON(400, gin.H{"error": "at least one filter is required"})
		return
	}

	files, err := driveconnector.FindEntities(filters)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(files) == 0 {
		c.JSON(404, gin.H{"error": entity.ERROR_FILE_NOT_FOUND})
		return
	}

	var fileid, filename, compareValue string
	if filters.Fileid != "" {
		compareValue = filters.Fileid
	} else {
		compareValue = filters.Name
	}
	for _, f := range files {
		var filevalue string
		if filters.Fileid != "" {
			filevalue = f.Id
		} else {
			filevalue = f.Name
		}
		if strings.EqualFold(compareValue, filevalue) {
			fileid = f.Id
			filename = f.Name
			break
		}

	}

	if fileid == "" {
		c.JSON(404, gin.H{"error": entity.ERROR_FILE_NOT_FOUND})
		return
	}

	if content, err := driveconnector.GetFileContent(fileid); err == nil {
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(200, "application/octet-stream", content)
	} else {
		c.JSON(500, gin.H{"error": err.Error()})
	}
}

func UploadFile(c *gin.Context) {
	request := new(struct {
		entity.DriveEntityFileInfo
		Content []*multipart.FileHeader `form:"content" binding:"required"`
	})
	if err := c.Bind(request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if request.Name == "" || len(request.Content) == 0 {
		c.JSON(400, gin.H{"error": "invalid request data"})
		return
	}

	file, err := request.Content[0].Open()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to read file content"})
		return
	}

	options := &driveconnector.UploadOptions{}
	if request.ParentId != "" || request.ParentName != "" {
		options.Parent = &entity.DriveEntityFileInfo{
			Fileid: request.ParentId,
			Name:   request.ParentName,
		}
	}

	if _, err = driveconnector.UploadFile(request.Name, content, options); err == nil {
		c.JSON(200, gin.H{"message": "File uploaded"})
	} else {
		c.JSON(500, gin.H{"error": err.Error()})
	}
}

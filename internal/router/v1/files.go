package v1

import (
	"fmt"
	"grubzo/internal/models/entity"
	"grubzo/internal/router/utils"
	"grubzo/internal/utils/ce"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

func (h Handlers) FileUpload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		ce.BadRequestBody(c)
	}
	files := form.File["files"]
	filesMeta := []map[string]any{}
	for idx, f := range files {
		args, err := utils.BuildFileSaveArgs(f, 2, nil, entity.O_TYPE_ITEM, idx+1)
		if err != nil {
			ce.RespondWithError(c, err)
			return
		}
		fileMeta, err := h.SS.FileManager.Save(args)
		if err != nil {
			ce.RespondWithError(c, err)
			return
		}
		filesMeta = append(filesMeta, fileMeta.JSON())
	}

	c.JSON(http.StatusOK, gin.H{
		"filesMeta": filesMeta,
	})

}

func (h Handlers) GetFileByID(c *gin.Context) {
	var params struct {
		ID string `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		ce.BadRequestParams(c)
		return
	}
	fileID, err := uuid.FromString(params.ID)
	if err != nil {
		ce.BadRequestParams(c)
		return
	}
	fileMeta, err := h.SS.FileManager.Get(fileID, 2)
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	f, err := fileMeta.Open()
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	defer f.Close()

	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fileMeta.GetFileName()))
	c.Header("Content-Type", fileMeta.GetMIMEType())
	c.Header("Content-Length", fmt.Sprintf("%d", fileMeta.GetFileSize()))

	if _, err := io.Copy(c.Writer, f); err != nil {
		ce.RespondWithError(c, err)
		return
	}
}

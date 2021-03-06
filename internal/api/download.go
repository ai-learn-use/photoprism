package api

import (
	"fmt"
	"path"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/repo"
	"github.com/photoprism/photoprism/internal/util"

	"github.com/gin-gonic/gin"
)

// TODO: GET /api/v1/dl/file/:hash
// TODO: GET /api/v1/dl/photo/:uuid
// TODO: GET /api/v1/dl/album/:uuid

// GET /api/v1/download/:hash
//
// Parameters:
//   hash: string The file hash as returned by the search API
func GetDownload(router *gin.RouterGroup, conf *config.Config) {
	router.GET("/download/:hash", func(c *gin.Context) {
		fileHash := c.Param("hash")

		r := repo.New(conf.OriginalsPath(), conf.Db())
		file, err := r.FindFileByHash(fileHash)

		if err != nil {
			c.AbortWithStatusJSON(404, gin.H{"error": err.Error()})
			return
		}

		fileName := path.Join(conf.OriginalsPath(), file.FileName)

		if !util.Exists(fileName) {
			log.Errorf("could not find original: %s", fileHash)
			c.Data(404, "image/svg+xml", photoIconSvg)

			// Set missing flag so that the file doesn't show up in search results anymore
			file.FileMissing = true
			conf.Db().Save(&file)
			return
		}

		downloadFileName := file.DownloadFileName()

		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", downloadFileName))

		c.File(fileName)
	})
}

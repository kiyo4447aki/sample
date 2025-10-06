package handlers

import (
	"api/config"
	"api/records"
	"errors"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ListRecordsHandler(cfg *config.Config) gin.HandlerFunc{
	return func (c *gin.Context){
		var err error

		//クエリパラメータでカメラ名としてディレクトリを取得
		subDir := c.Query("camera")
		cleanSub := filepath.Clean(subDir)
		if subDir == "" || cleanSub == "." {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error" : "invalid camera name",
			})
			return 
		}

		dirPath := filepath.Join(cfg.BasePath, cleanSub)

		/*
		ァイルパスに".."を含まない前提のもとディレクトリ・トラバーサルを防止
		仕様が変わる場合は変更が必要
		*/
		rel, err := filepath.Rel(cfg.BasePath, dirPath)
		if err != nil || strings.HasPrefix(rel, "..") {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error" : "invalid camera name",
			})
			return
		}


		//クエリパラメータからフィルタする日付を取得
		dateParam := c.Query("date")
		useDateFilter := false
		var filterDate time.Time
		if dateParam != "" {
			filterDate, err = time.Parse("2006-01-02", dateParam)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "failed",
					"error": "invalid date format",
				})
				return
			}
			useDateFilter = true
		}

		//指定されたカメラの全録画ファイル名の配列を取得
		names, err := records.GetRecordNames(dirPath)
		if err != nil {
			if errors.Is(err, records.ErrCameraNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"status": "failed",
					"error" : "camera not found",
				})
				return 
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "failed",
					"error" : "failed to get records",
				})
				return 
			} 
		}

		//詳細情報を取得
		recs := make([]records.Record, 0, len(names))
		for _, name := range names {
			recordInfo, err := records.GetRecordInfoWithDateFilter(
				name,
				dirPath,
				useDateFilter,
				filterDate,
				cleanSub,
				cfg.BaseUrl,
			)
			if err != nil {
				continue
			}

			recs = append(recs, recordInfo)

		}

		c.JSON(http.StatusOK, records.RecordsResponse{
			Status: "success",
			Recordings: recs,
			Count: len(recs),
		})
	}
}
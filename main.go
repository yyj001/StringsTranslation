package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	db = ConnectDB()
	var r = gin.Default()
	// 增
	r.POST("strings/add", AddString)
	r.POST("strings/modify", ModifyString)
	r.GET("strings/delete", DeleteString)
	r.GET("strings/query", QueryString)
	r.GET("strings/download", DownloadStringsFile)

	PORT := "9090"
	r.Run(":" + PORT)
}

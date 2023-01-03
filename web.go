package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func QueryString(context *gin.Context) {
	var datas []TableStrings
	word := context.Query("word")
	pageSize, _ := strconv.Atoi(context.Query("pageSize"))
	pageCursor, _ := strconv.Atoi(context.Query("pageCursor"))
	if pageSize <= 0 {
		OnError("pageSize 必须大于0", context)
		return
	}
	if pageCursor <= 0 {
		OnError("pageCursor 必须大于0", context)
		return
	}
	offset := (pageCursor - 1) * pageSize
	if len(word) == 0 {
		db.Limit(pageSize).Offset(offset).Find(&datas)
	} else {
		db.Where("translate_strs LIKE ?", "%"+word+"%").Limit(pageSize).Offset(offset).Find(&datas)
	}
	context.JSON(200, gin.H{
		"message": "success",
		"code":    "200",
		"data": gin.H{
			"total":      len(datas),
			"pageCursor": pageCursor,
			"list":       datas,
		},
	})
}

func AddString(context *gin.Context) {
	var json TableStrings
	err := context.ShouldBind(&json)
	var existDatas []TableStrings
	db.Where("name = ?", json.Name).Find(&existDatas)
	if len(existDatas) != 0 {
		OnError("已存在name相同文案", context)
		return
	}
	if err != nil {
		OnError(err.Error(), context)
		return
	}
	db.Create(&json)
	OnSuccess(context)
	isNeedUpdateFile = true
}

func DeleteString(context *gin.Context) {
	var datas []TableStrings
	name := context.Query("name")
	fmt.Println("delete name is " + name)
	db.Where("name = ?", name).Find(&datas)
	if len(datas) == 0 {
		OnError("没有此记录", context)
		return
	}
	db.Where("name = ?", name).Delete(&datas)
	OnSuccess(context)
	isNeedUpdateFile = true
}

func ModifyString(context *gin.Context) {
	var datas []TableStrings
	var data TableStrings
	modifyId := context.Query("id")
	err := context.ShouldBind(&data)
	db.Where("id = ?", modifyId).Find(&datas)
	if len(datas) == 0 {
		OnError("没有该文案记录", context)
		return
	}
	if err != nil {
		OnError(err.Error(), context)
		return
	}
	//修改
	db.Where("id = ?", modifyId).Updates(&data)
	OnSuccess(context)
	isNeedUpdateFile = true
}

func OnError(errStr string, context *gin.Context) {
	context.JSON(200, gin.H{
		"message": errStr,
		"code":    "400",
	})
	fmt.Println("error:", errStr)
}

func OnSuccess(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "success",
		"code":    "200",
	})
}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func QueryString(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "get String success")
	if r.Method == "POST" {
		err := r.ParseForm() // 解析 url 传递的参数，对于 POST 则解析响应包的主体（request body）
		if err != nil {
			log.Fatal("ParseForm: ", err)
		}
		// 请求的是登录数据，那么执行登录的逻辑判断
		fmt.Println("pCursor:", r.Form["pCursor"])
		fmt.Println("word:", r.Form["word"])
	}
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

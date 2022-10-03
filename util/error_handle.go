package util

import "github.com/gin-gonic/gin"

// CheckError 错误处理
func CheckError(err error){
	if err!=nil{
		panic(err.Error())
	}
}

// Sunccess 成功处理
func Sunccess(msg string,c *gin.Context )  {
	c.JSON(200,gin.H{"message":msg})
}


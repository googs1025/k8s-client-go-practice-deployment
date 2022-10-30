package main

import (
	"github.com/gin-gonic/gin"
	"k8s-client-go-api-practice/core"
	"k8s-client-go-api-practice/debug"
	"k8s-client-go-api-practice/deployment"
	"k8s-client-go-api-practice/util"
	"net/http"
)

func main() {


	r := gin.New()

	r.Use(func(c *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				c.AbortWithStatusJSON(400,gin.H{"error":e})
			}
		}()
		c.Next()
	})

	deployment.RegHandlers(r)
	debug.RegHandlers(r)

	r.Static("/static", "./static")
	r.LoadHTMLGlob("html/**/*")
	r.GET("/deployments", func(c *gin.Context) {
		//ns := c.Params("namespace")
		c.HTML(http.StatusOK, "deployment_list.html",
			util.DataBuilder().
			SetTitle("deployment列表").
			SetData("DepList",deployment.ListAll("default")))
		//c.JSON(200, deployment.ListAll("default"))
	})
	r.GET("/deployments/:name", func(c *gin.Context) {
		c.HTML(http.StatusOK, "deployment_detail.html",
			util.DataBuilder().
				SetTitle("deployment详细-"+c.Param("name")).
				SetData("DepDetail",deployment.GetDeployment("default", c.Param("name"))))
	})

	r.GET("/create/deployments", func(c *gin.Context) {
		c.HTML(http.StatusOK, "deployment_create.html",
			util.DataBuilder().
				SetTitle("deployment创建"))
	})

	core.InitDeployment() //初始化第一批deployment list
	_ = r.Run(":8080")

}

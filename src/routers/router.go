package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rtsp2hls/src/help"
	"rtsp2hls/src/models"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	r.GET("/cameras", func(c *gin.Context) {
		GetCameras(c)
	})
	r.PATCH("/camera/:id", func(c *gin.Context) {
		UpdateCamera(c)
	})
	r.GET("/camera/start/:id", func(c *gin.Context) {
		StartCamera(c)
	})
	r.GET("/camera/stop/:id", func(c *gin.Context) {
		StopCamera(c)
	})
	return r
}

func GetCameras(c *gin.Context) {
	data := make(map[string]interface{})

	data["lists"], data["total"] = models.GetCameras()
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func UpdateCamera(c *gin.Context) {
	ID := c.Param("id")

	data := make(map[string]interface{})
	if camera, err := models.GetCamera(ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err,
		})
	} else {
		if error := c.ShouldBindJSON(&data); error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err,
			})
		} else {
			if updateCamera, error := camera.UpdateCamera(data); error != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": err,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"update": updateCamera,
				})
			}
		}
	}
}

func StartCamera(c *gin.Context) {
	ID := c.Param("id")
	if camera, err := models.GetCamera(ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err,
		})
	} else {
		if camera.State != 1 {
			if _, err := camera.OpenCamera(); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": err,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"msg": "start successful",
				})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "this camera is already started, please do not repeat",
			})
		}

	}
}

func StopCamera(c *gin.Context) {
	ID := c.Param("id")
	if camera, err := models.GetCamera(ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err,
		})
	} else {
		if _, err := camera.CloseCamera(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"msg": "close successful",
			})
		}
	}
}

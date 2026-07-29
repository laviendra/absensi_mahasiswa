package routes

import (
	"github.com/gin-gonic/gin"
	"absensi-mahasiswa/controllers"
)


func AuthRoute(router *gin.Engine){

	router.POST("/login", controllers.Login)

}
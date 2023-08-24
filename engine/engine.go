package engine

import (
	"github.com/gin-gonic/gin"
	"io"
	"net"
	"net/http"
	"os"
)

func InitEngine() {
	// init IPCity Data
	InitIPCity()
	// disable color
	gin.DisableConsoleColor()
	// init log file
	_ = os.MkdirAll("logs", 0755)
	file, _ := os.OpenFile("logs/gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	gin.DefaultWriter = io.MultiWriter(file)
	// init Engine
	engine := gin.Default()
	// init route
	engine.GET("search/", searchIPAddress)
	// Start Engine
	err := engine.Run(":8000")
	if err != nil {
		panic(err.Error())
	}
}

func searchIPAddress(context *gin.Context) {
	// load params
	ip := context.Query("ip")
	if _ip := net.ParseIP(ip); _ip == nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"ip":          "",
			"country":     "",
			"province":    "",
			"city":        "",
			"district":    "",
			"isp":         "",
			"backboneISP": "",
			"countryCode": 0,
			"areaCode":    0,
		})
		return
	}
	// search ip
	meta := IPCityClient.Search(ip)
	context.JSON(http.StatusOK, gin.H{
		"ip":          ip,
		"country":     meta.Country(),
		"province":    meta.Province(),
		"city":        meta.City(),
		"district":    meta.District(),
		"isp":         meta.ISP(),
		"backboneISP": meta.BackboneISP(),
		"countryCode": meta.CountryCode(),
		"areaCode":    meta.AreaCode(),
	})
}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ApiInfoAdmin struct {
	ApiInfo
	Secret string `json:"secret"`
}

var apiInfoValue *ApiInfo
var apiInfoAdminValue *ApiInfoAdmin

func main() {
	apiInfoValue = &ApiInfo{
		Name:    "keycloak_demo_01__server",
		Version: "0.1.0",
	}
	apiInfoAdminValue = &ApiInfoAdmin{
		ApiInfo: *apiInfoValue,
		Secret:  "42",
	}

	ginEngine := gin.Default()

	apiv1 := ginEngine.Group("/api/v1")
	apiv1.GET("info", info)
	apiv1.GET("info_admin", infoAdmin)

	err := ginEngine.Run("0.0.0.0:8001")
	if err != nil {
		fmt.Println(err)
	}
}

func info(context *gin.Context) {
	context.JSONP(http.StatusOK, apiInfoValue)
}

func infoAdmin(context *gin.Context) {
	context.JSONP(http.StatusOK, apiInfoAdminValue)
}
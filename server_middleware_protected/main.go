package main

import (
	"fmt"
	"github.com/Nerzal/gocloak/v5"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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

var gocloakClient gocloak.GoCloak

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
	gocloakClient = gocloak.NewClient("http://keykloak:8080")

	ginEngine.Use(parseAndValidateJwt())

	apiv1 := ginEngine.Group("/api/v1")

	infoGroup := apiv1.Group("info")
	infoGroup.GET("", info)

	infoAdminGroup := apiv1.Group("info_admin")
	infoAdminGroup.Use(hasRole("gallery_admin"))
	infoAdminGroup.GET("", infoAdmin)

	err := ginEngine.Run("0.0.0.0:8000")
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

func parseAndValidateJwt() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		authHeader := ginContext.GetHeader("Authorization")
		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")

		jwtToken, claims, err := gocloakClient.DecodeAccessToken(bearerToken, "testrealm")

		if err != nil || !jwtToken.Valid {
			_ = ginContext.AbortWithError(http.StatusForbidden, err)
			return
		}

		if resourceAccess, ok := (*claims)["resource_access"].(map[string]interface{}); ok {
			if demoGallery, ok := resourceAccess["demo-gallery"].(map[string]interface{}); ok {
				if roles, ok := demoGallery["roles"].([]interface{}); ok {
					rolesStringSlice := make([]string, 0)
					for _, roleInterface := range roles {
						if role, ok := roleInterface.(string); ok {
							rolesStringSlice = append(rolesStringSlice, role)
						}
					}
					ginContext.Set("roles", rolesStringSlice)
				}
			}
		}
	}
}

func hasRole(requiredRole string) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		rolesInterace, _ := ginContext.Get("roles")
		roles := rolesInterace.([]string)

		hasRole := false
		for _, role := range roles {
			if role == requiredRole {
				hasRole = true
				break
			}
		}
		if !hasRole {
			ginContext.AbortWithStatus(http.StatusForbidden)
		}
	}
}

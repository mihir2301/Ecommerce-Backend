package routes

import (
	"fmt"
	"golang_project/auth"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc func(*gin.Context)
}

type Router []Route // problem

type Routes struct {
	router *gin.Engine
}

func CORSMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Writer.Header().Set("content-type", "application/json")
		c.Writer.Header().Set("access-control-allow-origin", "*")
		c.Writer.Header().Set("access-control-allow-credentials", "true")
		c.Writer.Header().Set("access-control-allow-methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Writer.Header().Set("access-control-allow-headers", "Authorization,content-type,X-requested-width")
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusNoContent)
			c.Abort()
			return
		}
	}
}

func (r *Routes) EcommerceHealthCheck(rg *gin.RouterGroup) {
	orderRoute := rg.Group("/ecommerce")
	orderRoute.Use(CORSMiddleware())
	for _, routes := range healthCheckRoutes {
		switch routes.Method {
		case "GET":
			orderRoute.GET(routes.Pattern, routes.HandleFunc)
		case "POST":
			orderRoute.POST(routes.Pattern, routes.HandleFunc)
		case "OPTIONS":
			orderRoute.OPTIONS(routes.Pattern, routes.HandleFunc)
		case "PUT":
			orderRoute.PUT(routes.Pattern, routes.HandleFunc)
		case "DELETE":
			orderRoute.DELETE(routes.Pattern, routes.HandleFunc)
		default:
			orderRoute.GET(routes.Pattern, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Choose a correct method",
				})
			})
		}
	}
}

func (r *Routes) EcommerceUser(rg *gin.RouterGroup) {
	userRoute := rg.Group("/ecommerce")
	userRoute.Use(CORSMiddleware())
	for _, route := range userRoutes {
		switch route.Method {
		case "GET":
			userRoute.GET(route.Pattern, route.HandleFunc)
		case "POST":
			userRoute.POST(route.Pattern, route.HandleFunc)
		case "OPTIONS":
			userRoute.OPTIONS(route.Pattern, route.HandleFunc)
		case "PUT":
			userRoute.PUT(route.Pattern, route.HandleFunc)
		case "DELETE":
			userRoute.DELETE(route.Pattern, route.HandleFunc)
		default:
			userRoute.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Choose a correct method",
				})
			})
		}
	}
}
func (r *Routes) EcommerceProduct(rg *gin.RouterGroup) {
	productroute := rg.Group("/ecommerce")
	productroute.Use(CORSMiddleware())
	for _, route := range productroutes {
		switch route.Method {
		case "GET":
			productroute.GET(route.Pattern, route.HandleFunc)
		case "POST":
			productroute.POST(route.Pattern, route.HandleFunc)
		case "OPTIONS":
			productroute.OPTIONS(route.Pattern, route.HandleFunc)
		case "PUT":
			productroute.PUT(route.Pattern, route.HandleFunc)
		case "DELETE":
			productroute.DELETE(route.Pattern, route.HandleFunc)
		default:
			productroute.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(http.StatusAccepted, gin.H{
					"message": "Choose a correct method",
				})
			})
		}
	}
}

func (r *Routes) EcommerceGlobalProductRoutes(rg *gin.RouterGroup) {
	productroute := rg.Group("/ecommerce")
	productroute.Use(CORSMiddleware())
	for _, route := range globalproducts {
		switch route.Method {
		case "GET":
			productroute.GET(route.Pattern, route.HandleFunc)
		case "POST":
			productroute.POST(route.Pattern, route.HandleFunc)
		case "OPTIONS":
			productroute.OPTIONS(route.Pattern, route.HandleFunc)
		case "PUT":
			productroute.PUT(route.Pattern, route.HandleFunc)
		case "DELETE":
			productroute.DELETE(route.Pattern, route.HandleFunc)
		default:
			productroute.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(http.StatusAccepted, gin.H{
					"message": "Choose a correct method",
				})
			})
		}
	}
}

func (r *Routes) EcommerceAuthUser(rg *gin.RouterGroup) {
	authuser := rg.Group("/ecommerce")
	authuser.Use(CORSMiddleware())
	for _, route := range userauthroutes {
		switch route.Method {
		case "GET":
			authuser.GET(route.Pattern, route.HandleFunc)
		case "POST":
			authuser.POST(route.Pattern, route.HandleFunc)
		case "OPTIONS":
			authuser.OPTIONS(route.Pattern, route.HandleFunc)
		case "PUT":
			authuser.PUT(route.Pattern, route.HandleFunc)
		case "DELETE":
			authuser.DELETE(route.Pattern, route.HandleFunc)
		default:
			authuser.GET(route.Pattern, func(g *gin.Context) {
				g.JSON(http.StatusAccepted, gin.H{
					"message": "Choose a correct method",
				})
			})
		}
	}
}

func ClientRoutes() {
	r := Routes{
		router: gin.Default(),
	}
	apiversion := os.Getenv("API_VERSION")
	if apiversion == "" {
		fmt.Println("Please give apiversion")
	}
	v1 := r.router.Group(apiversion)
	r.EcommerceHealthCheck(v1)
	r.EcommerceUser(v1)
	r.EcommerceGlobalProductRoutes(v1)
	//auth
	v1.Use(auth.Auth())
	r.EcommerceProduct(v1)
	r.EcommerceAuthUser(v1)
	err := r.router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Println("Failed to run the server")
	}
}

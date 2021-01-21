package main

import (
	"log"
	"lyrical-app/database"
	"lyrical-app/graph"
	"lyrical-app/graph/generated"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db := database.NewDb()

	// router := chi.NewRouter()
	router := gin.Default()
	router.Use(cors.Default())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	// To be able to send tokens to the server.
	corsConfig.AllowCredentials = true

	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("OPTIONS")
	router.Use(cors.New(corsConfig))

	router.GET("/", playgroundHandler())
	router.POST("/query", graphqlHandler(db))
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	router.Run(":" + port)
}

// Defining the Graphql handler
func graphqlHandler(db *database.DB) gin.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(db)}))
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

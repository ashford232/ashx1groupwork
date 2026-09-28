package main

import (
	"kolu/internal/config"
	"kolu/internal/database"
	"kolu/internal/handlers"
	"kolu/internal/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	port := cfg.Port
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.SetTrustedProxies(nil)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "kolu api",
		})
	})
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg))

	// Profile
	protected.POST("/profile/upload-url", handlers.GetPresignedAvatarUploadHandler(cfg))
	protected.POST("/profile/setup", handlers.SetupProfileHandler(pool))
	// Todos
	protected.POST("/todos", handlers.CreateTodoHandler(pool))
	protected.GET("/todos", handlers.GetAllTodosHandler(pool))
	protected.GET("/todos/trash", handlers.GetAllTrashedTodosHandler(pool))
	protected.GET("/todos/:id", handlers.GetTodoById(pool))
	protected.PUT("/todos/:id", handlers.UpdateTodo(pool))
	protected.DELETE("/todos/:id", handlers.DeleteTodoById(pool))
	protected.DELETE("/todos", handlers.DeleteAllTodos(pool))

	// Notes
	protected.POST("/notes", handlers.CreateNoteHandler(pool))
	protected.GET("/notes", handlers.GetAllNotesHandler(pool))
	protected.GET("/notes/trash", handlers.GetAllTrashedNotesHandler(pool))
	protected.GET("/notes/:id", handlers.GetNoteById(pool))
	protected.PUT("/notes/:id", handlers.UpdateNote(pool))
	protected.DELETE("/notes/:id", handlers.DeleteNoteById(pool))
	protected.DELETE("/notes", handlers.DeleteAllNotes(pool))

	// Docs
	protected.POST("/documents", handlers.CreateDocHandler(pool))
	protected.GET("/documents", handlers.GetAllDocsHandler(pool))
	protected.GET("/documents/trash", handlers.GetAllTrashedDocsHandler(pool))
	protected.GET("/documents/:id", handlers.GetDocById(pool))
	protected.PUT("/documents/:id", handlers.UpdateDoc(pool))
	protected.DELETE("/documents/:id", handlers.DeleteDocById(pool))
	protected.DELETE("/documents", handlers.DeleteAllDocs(pool))

	// Labels
	protected.POST("/labels", handlers.CreateLabelHandler(pool))
	protected.GET("/labels", handlers.GetAllUserLabels(pool))
	protected.GET("/labels/:id", handlers.GetLabelHandler(pool))
	protected.PUT("/labels/:id", handlers.UpdateLabelHandler(pool))
	protected.DELETE("/labels/:id", handlers.DeleteLabelHandler(pool))

	protected.POST("/todos/:todo_id/labels", handlers.AddTodoLabelHandler(pool))
	protected.POST("/notes/:note_id/labels", handlers.AddNoteLabelHandler(pool))
	protected.POST("/documents/:document_id/labels", handlers.AddDocLabelHandler(pool))

	// User
	protected.GET("/user/:id", handlers.GetUserByIDHandler(pool))
	protected.GET("/user", handlers.GetUser(pool))

	auth := router.Group("/auth")
	auth.POST("/signup", handlers.SignUpUserHandler(pool))
	auth.POST("/signin", handlers.SignInUserHandler(pool, cfg))
	router.Run(":" + port)
}

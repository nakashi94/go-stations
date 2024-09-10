package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	mux.HandleFunc(handler.NewHealthzHandler().Path, handler.NewHealthzHandler().ServeHTTP)
	todoService := service.NewTODOService(todoDB)
	mux.HandleFunc(handler.NewTODOHandler(todoService).Path, handler.NewTODOHandler(todoService).ServeHTTP)
	return mux
}

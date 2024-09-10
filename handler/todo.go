package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc  *service.TODOService
	Path string
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc:  svc,
		Path: "/todos",
	}
}

// ServeHTTP implements http.Handler interface
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var data model.CreateTODORequest
		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
			return
		}

		// subjectが空であればBad Requestを返す
		if len(data.Subject) == 0 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			println("subject length is 0")
			return
		}

		CreateTodoResponse, err := h.Create(r.Context(), &data)
		if err != nil {
			log.Println(err)
			fmt.Println("createTodoResponse Error")
			http.Error(w, "Failed to create TODO", http.StatusBadRequest)
			return
		}

		encoder := json.NewEncoder(w)
		err = encoder.Encode(CreateTodoResponse)
		if err != nil {
			log.Println(err)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{TODO: *todo}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

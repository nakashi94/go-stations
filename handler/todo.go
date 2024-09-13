package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

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
	case http.MethodGet:
		// GET method process
		prevIDStr := r.URL.Query().Get("prev_id")
		sizeStr := r.URL.Query().Get("size")
		var prevID int64 = 0
		var size int64 = 10
		if len(prevIDStr) != 0 {
			var err error
			prevID, err = strconv.ParseInt(prevIDStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid prev_id parameter", http.StatusBadRequest)
				return
			}
		}
		if len(sizeStr) != 0 {
			var err error
			size, err = strconv.ParseInt(sizeStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid Size", http.StatusBadRequest)
				return
			}
		}

		readTodoResponse, err := h.svc.ReadTODO(r.Context(), prevID, size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode(map[string]interface{}{"todos": readTodoResponse})
		if err != nil {
			http.Error(w, "encode response data to json format", http.StatusBadRequest)
			return
		}

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

		if len(data.Subject) == 0 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		createTodoResponse, err := h.Create(r.Context(), &data)
		if err != nil {
			http.Error(w, "Failed to create TODO", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode(createTodoResponse)
		if err != nil {
			log.Println(err)
			return
		}

	case http.MethodPut:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Faild to read request body.", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var data model.UpdateTODORequest
		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
			return
		}

		if data.ID == 0 {
			http.Error(w, "Invalid ID: ID is more than 1.", http.StatusBadRequest)
			return
		}

		if len(data.Subject) == 0 {
			http.Error(w, "Invaild Subject: The length of Subject is more than 1.", http.StatusBadRequest)
			return
		}

		updateTodoResponse, err := h.Update(r.Context(), &data)
		if err != nil {
			http.Error(w, "Failed to update TODO.", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode(updateTodoResponse)
		if err != nil {
			log.Println(err)
			return
		}

	case http.MethodDelete:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var data model.DeleteTODORequest
		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
			return
		}

		if len(data.IDs) == 0 {
			http.Error(w, "IDs is empty array.", http.StatusBadRequest)
			return
		}

		deleteTodoResponse, err := h.Delete(r.Context(), &data)
		if err != nil {
			http.Error(w, "TODO not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode(deleteTodoResponse)
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
// func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
// 	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &model.ReadTODOResponse{TODOs: todos}, nil
// }

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{TODO: *todo}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return &model.DeleteTODOResponse{}, nil
}

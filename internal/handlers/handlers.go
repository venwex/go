package handlers

import (
	"encoding/json"
	e "example/test/internal/errors"
	"example/test/internal/service"
	u "example/test/internal/utils"
	"net/http"
)

type H map[string]any

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	id, err := u.GetIDFromQuery(r)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.service.ServiceGetTask(id)
	if err != nil {
		u.RenderError(w, http.StatusNotFound, err.Error())
		return
	}

	u.RenderJSON(w, http.StatusOK, task)
}

func (h *Handler) HandleGetTasksList(w http.ResponseWriter, r *http.Request) {
	tasks := h.service.ServiceGetTasks()

	u.RenderJSON(w, http.StatusOK, tasks)
}

func (h *Handler) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr != "" {
		h.HandleGetTask(w, r)
		return
	}

	h.HandleGetTasksList(w, r)
}

func (h *Handler) HandlePostTask(w http.ResponseWriter, r *http.Request) {
	title, err := u.DecodeTask(r)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(title) == 0 {
		u.RenderError(w, http.StatusBadRequest, e.ErrInvalidTitleName.Error())
	}

	task, err := h.service.ServiceCreateTask(title)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, e.ErrInvalidTitleName.Error())
		return
	}

	u.RenderJSON(w, http.StatusCreated, task)
}

func (h *Handler) HandlePatchTask(w http.ResponseWriter, r *http.Request) {
	id, err := u.GetIDFromQuery(r)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req struct {
		Done bool `json:"done"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.RenderError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.service.ServiceMarkDoneTask(id, req.Done); err != nil {
		if err == e.ErrTaskNotFound {
			u.RenderError(w, http.StatusNotFound, err.Error())
			return
		}
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	u.RenderJSON(w, http.StatusOK, H{
		"updated": true,
	})
}

func (h *Handler) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := u.GetIDFromQuery(r)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.service.ServiceDeleteTask(id)
	if err != nil {
		u.RenderError(w, http.StatusNotFound, err.Error())
		return
	}

	u.RenderJSON(w, http.StatusOK, task)
}

package tasks

import (
	"encoding/json"
	"errors"
	e "example/test/internal/errors"
	"example/test/internal/service/tasks"
	u "example/test/internal/utils"
	"net/http"
)

type H map[string]any

type TaskHandler struct {
	Task *tasks.TaskService
}

func NewTaskHandler(task *tasks.TaskService) *TaskHandler {
	return &TaskHandler{Task: task}
}

func (h *TaskHandler) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	id, err := u.GetIDFromQuery(r)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.Task.ServiceGetTask(id)
	if err != nil {
		u.RenderError(w, http.StatusNotFound, err.Error())
		return
	}

	u.RenderJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) HandleGetTasksList(w http.ResponseWriter) {
	getTasks, err := h.Task.ServiceGetTasks()
	if err != nil {
		u.RenderError(w, http.StatusNotFound, err.Error())
	}

	u.RenderJSON(w, http.StatusOK, getTasks)
}

func (h *TaskHandler) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr != "" {
		h.HandleGetTask(w, r)
		return
	}

	h.HandleGetTasksList(w)
}

func (h *TaskHandler) HandlePostTask(w http.ResponseWriter, r *http.Request) {
	title, err := u.DecodeTask(r)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(title) == 0 {
		u.RenderError(w, http.StatusBadRequest, e.ErrInvalidTitleName.Error())
		return
	}

	task, err := h.Task.ServiceCreateTask(title)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	u.RenderJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) HandlePatchTask(w http.ResponseWriter, r *http.Request) {
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

	if err := h.Task.ServiceMarkDoneTask(id, req.Done); err != nil {
		if errors.Is(err, e.ErrTaskNotFound) {
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

func (h *TaskHandler) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := u.GetIDFromQuery(r)
	if err != nil {
		u.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.Task.ServiceDeleteTask(id)
	if err != nil {
		u.RenderError(w, http.StatusNotFound, err.Error())
		return
	}

	u.RenderJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	u.RenderJSON(w, http.StatusOK, H{
		"healthy": true,
	})
}

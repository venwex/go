package task

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"bool"`
}

func NewTask(id int, title string) Task {
	return Task{
		ID:    id,
		Title: title,
		Done:  false,
	}
}

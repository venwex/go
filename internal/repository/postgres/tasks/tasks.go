package tasks

import (
	"database/sql"
	"errors"
	e "example/test/internal/errors"
	m "example/test/internal/models"
	"example/test/internal/repository/postgres"
)

type TaskRepository struct {
	db *postgres.Dialect
}

func NewRepository(db *postgres.Dialect) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) GetTask(id int) (m.Task, error) {
	var task m.Task

	err := r.db.DB.Get(&task,
		`SELECT id, title, done FROM tasks WHERE id = $1`,
		id,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m.Task{}, e.ErrTaskNotFound
		}
		return m.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) GetTasks() ([]m.Task, error) {
	var tasks []m.Task

	err := r.db.DB.Select(&tasks,
		`SELECT id, title, done FROM tasks`,
	)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) CreateTask(title string) (m.Task, error) {
	var task m.Task

	err := r.db.DB.Get(&task,
		`INSERT INTO tasks (title)
		 VALUES ($1)
		 RETURNING id, title, done`,
		title,
	)

	if err != nil {
		return m.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) MarkDoneTask(id int, done bool) error {
	res, err := r.db.DB.Exec(
		`UPDATE tasks SET done = $1 WHERE id = $2`,
		done, id,
	)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return e.ErrTaskNotFound
	}

	return nil
}

func (r *TaskRepository) DeleteTask(id int) (m.Task, error) {
	var task m.Task

	err := r.db.DB.Get(&task,
		`DELETE FROM tasks
		 WHERE id = $1
		 RETURNING id, title, done`,
		id,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m.Task{}, e.ErrTaskNotFound
		}
		return m.Task{}, err
	}

	return task, nil
}

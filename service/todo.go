package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// if len(subject) == 0 {
	// 	// log.Println("subject must be 1 character more.")
	// 	// err := errors.New("subject must be 1 character more")
	// 	var err error
	// 	return nil, err
	// }

	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var createdAt time.Time
	var updatedAt time.Time

	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&subject, &description, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	todo := &model.TODO{
		ID:          id,
		Subject:     subject,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	if size == 0 {
		return []*model.TODO{}, nil
	}

	var rows *sql.Rows
	var err error

	if prevID > 0 {
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	} else {
		rows, err = s.db.QueryContext(ctx, read, size)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*model.TODO
	for rows.Next() {
		var todo model.TODO
		err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if todos == nil {
		todos = []*model.TODO{}
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// if len(subject) == 0 {
	// 	var err error
	// 	return nil, err
	// }

	result, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}

	affectedRowCount, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affectedRowCount == 0 {
		return nil, &model.ErrNotFound{When: time.Now(), What: "Todo Not Found."}
	}

	id, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var createdAt time.Time
	var updatedAt time.Time

	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&subject, &description, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	todo := &model.TODO{
		ID:          id,
		Subject:     subject,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	placeholder := strings.Repeat("?,", len(ids)-1) + "?"
	query := fmt.Sprintf(`DELETE FROM todos WHERE id IN (%s)`, placeholder)

	anyIDs := make([]interface{}, len(ids))
	for i, id := range ids {
		anyIDs[i] = id
	}

	res, err := s.db.ExecContext(ctx, query, anyIDs...)
	if err != nil {
		return fmt.Errorf("failed to delete todos: %w", err)
	}

	deletedCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	if deletedCount == 0 {
		return &model.ErrNotFound{}
	}

	return nil
}

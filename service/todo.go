package service

import (
	"context"
	"database/sql"
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

	return nil, nil
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
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}

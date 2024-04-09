package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/enchik0reo/commandApi/internal/models"
)

type CommandStoage struct {
	db *sql.DB
}

func NewCommandStorage(db *sql.DB) *CommandStoage {
	return &CommandStoage{db: db}
}

// CreateNew adds new command to db ...
func (c *CommandStoage) CreateNew(ctx context.Context, comandName string) (int64, error) {
	stmt, err := c.db.PrepareContext(ctx, "INSERT INTO commands (command_name) VALUES ($1) RETURNING command_id")
	if err != nil {
		return 0, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, comandName)

	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("can't insert source: %w", err)
	}

	var id int64

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("can't get last insert id: %w", err)
	}

	return id, nil
}

// GetList returns n latest commands ...
func (c *CommandStoage) GetList(ctx context.Context, n int64) ([]models.Command, error) {
	stmt, err := c.db.PrepareContext(ctx, `SELECT command_id, command_name, started_at, is_working 
	FROM commands ORDER BY command_id DESC LIMIT $1`)
	if err != nil {
		return nil, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, n)
	if err != nil {
		return nil, fmt.Errorf("can't get command's list: %w", err)
	}
	defer rows.Close()

	cmds := []models.Command{}

	for rows.Next() {
		cmd := models.Command{}
		var created time.Time

		if err := rows.Scan(&cmd.ID, &cmd.Name, &created, &cmd.IsWorking); err != nil {
			return nil, fmt.Errorf("can't scan row: %w", err)
		}

		cmd.StartedAt = created.UTC().Format(time.StampMilli)

		cmds = append(cmds, cmd)
	}

	return cmds, nil
}

// GetOne returns description of one command by command id ...
func (c *CommandStoage) GetOne(ctx context.Context, id int64) (*models.Command, error) {
	stmt, err := c.db.PrepareContext(ctx, `SELECT c.command_id, c.command_name, c.started_at, c.is_working, o.output 
	FROM commands c 
	INNER JOIN outputs o ON c.command_id = o.command_id 
	WHERE c.command_id = $1`)
	if err != nil {
		return nil, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't get command: %w", err)
	}
	defer rows.Close()

	outputs := []string{}
	cmd := models.Command{}
	var created time.Time

	for rows.Next() {
		var output string

		if err := rows.Scan(&cmd.ID, &cmd.Name, &created, &cmd.IsWorking, &output); err != nil {
			return nil, fmt.Errorf("can't scan row: %w", err)
		}

		cmd.StartedAt = created.UTC().Format(time.StampMilli)

		outputs = append(outputs, output)
	}

	cmd.Output = outputs

	return &cmd, nil
}

// StopOne stops the command by command id ...
func (c *CommandStoage) StopOne(ctx context.Context, id int64) (int64, error) {
	stmt, err := c.db.PrepareContext(ctx, `UPDATE commands SET is_working = false
	WHERE command_id = $1 RETURNING command_id`)
	if err != nil {
		return 0, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)

	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("can't stop command: %w", err)
	}

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("can't get stoped id: %w", err)
	}

	return id, nil
}

// SaveOutput saves command's output by command id ...
func (c *CommandStoage) SaveOutput(ctx context.Context, id int64, output string) (int64, error) {
	stmt, err := c.db.PrepareContext(ctx, "INSERT INTO outputs (command_id, output) VALUES ($1, $2) RETURNING output_id")
	if err != nil {
		return 0, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id, output)

	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("can't insert source: %w", err)
	}

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("can't get last insert id: %w", err)
	}

	return id, nil
}

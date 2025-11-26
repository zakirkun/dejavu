package repository

import (
	"database/sql"

	"github.com/dejavu/backend/internal/domain"
	"github.com/dejavu/backend/pkg/database"
)

type ProjectRepository struct {
	db *database.DB
}

func NewProjectRepository(db *database.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *domain.Project) error {
	query := `
		INSERT INTO projects (user_id, name, repo_url, build_command, output_dir)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query,
		project.UserID,
		project.Name,
		project.RepoURL,
		project.BuildCommand,
		project.OutputDir,
	).Scan(&project.ID, &project.CreatedAt)
}

func (r *ProjectRepository) GetByID(id string) (*domain.Project, error) {
	project := &domain.Project{}
	query := `
		SELECT id, user_id, name, repo_url, build_command, output_dir, created_at
		FROM projects
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&project.ID,
		&project.UserID,
		&project.Name,
		&project.RepoURL,
		&project.BuildCommand,
		&project.OutputDir,
		&project.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return project, err
}

func (r *ProjectRepository) ListByUserID(userID string) ([]*domain.Project, error) {
	query := `
		SELECT id, user_id, name, repo_url, build_command, output_dir, created_at
		FROM projects
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*domain.Project
	for rows.Next() {
		project := &domain.Project{}
		if err := rows.Scan(
			&project.ID,
			&project.UserID,
			&project.Name,
			&project.RepoURL,
			&project.BuildCommand,
			&project.OutputDir,
			&project.CreatedAt,
		); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *ProjectRepository) Update(project *domain.Project) error {
	query := `
		UPDATE projects
		SET name = $1, repo_url = $2, build_command = $3, output_dir = $4
		WHERE id = $5 AND user_id = $6
	`
	result, err := r.db.Exec(
		query,
		project.Name,
		project.RepoURL,
		project.BuildCommand,
		project.OutputDir,
		project.ID,
		project.UserID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ProjectRepository) Delete(id, userID string) error {
	query := `DELETE FROM projects WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}


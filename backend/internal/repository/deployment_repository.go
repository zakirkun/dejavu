package repository

import (
	"database/sql"

	"github.com/dejavu/backend/internal/domain"
	"github.com/dejavu/backend/pkg/database"
)

type DeploymentRepository struct {
	db *database.DB
}

func NewDeploymentRepository(db *database.DB) *DeploymentRepository {
	return &DeploymentRepository{db: db}
}

func (r *DeploymentRepository) Create(deployment *domain.Deployment) error {
	query := `
		INSERT INTO deployments (project_id, status, subdomain, commit_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		deployment.ProjectID,
		deployment.Status,
		deployment.Subdomain,
		deployment.CommitHash,
	).Scan(&deployment.ID, &deployment.CreatedAt, &deployment.UpdatedAt)
}

func (r *DeploymentRepository) GetByID(id string) (*domain.Deployment, error) {
	deployment := &domain.Deployment{}
	query := `
		SELECT id, project_id, status, subdomain, 
		       COALESCE(image_url, '') as image_url, 
		       COALESCE(commit_hash, '') as commit_hash, 
		       COALESCE(build_logs, '') as build_logs, 
		       created_at, updated_at
		FROM deployments
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&deployment.ID,
		&deployment.ProjectID,
		&deployment.Status,
		&deployment.Subdomain,
		&deployment.ImageURL,
		&deployment.CommitHash,
		&deployment.BuildLogs,
		&deployment.CreatedAt,
		&deployment.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return deployment, err
}

func (r *DeploymentRepository) ListByProjectID(projectID string) ([]*domain.Deployment, error) {
	query := `
		SELECT id, project_id, status, subdomain, 
		       COALESCE(image_url, '') as image_url, 
		       COALESCE(commit_hash, '') as commit_hash, 
		       COALESCE(build_logs, '') as build_logs, 
		       created_at, updated_at
		FROM deployments
		WHERE project_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deployments []*domain.Deployment
	for rows.Next() {
		deployment := &domain.Deployment{}
		if err := rows.Scan(
			&deployment.ID,
			&deployment.ProjectID,
			&deployment.Status,
			&deployment.Subdomain,
			&deployment.ImageURL,
			&deployment.CommitHash,
			&deployment.BuildLogs,
			&deployment.CreatedAt,
			&deployment.UpdatedAt,
		); err != nil {
			return nil, err
		}
		deployments = append(deployments, deployment)
	}
	return deployments, nil
}

func (r *DeploymentRepository) UpdateStatus(id string, status domain.DeploymentStatus) error {
	query := `
		UPDATE deployments
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *DeploymentRepository) UpdateImageURL(id, imageURL string) error {
	query := `
		UPDATE deployments
		SET image_url = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := r.db.Exec(query, imageURL, id)
	return err
}

func (r *DeploymentRepository) UpdateLogs(id, logs string) error {
	query := `
		UPDATE deployments
		SET build_logs = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := r.db.Exec(query, logs, id)
	return err
}

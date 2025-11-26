package domain

import "time"

type Project struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	RepoURL      string    `json:"repo_url"`
	BuildCommand string    `json:"build_command"`
	OutputDir    string    `json:"output_dir"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateProjectRequest struct {
	Name         string `json:"name" validate:"required"`
	RepoURL      string `json:"repo_url" validate:"required,url"`
	BuildCommand string `json:"build_command"`
	OutputDir    string `json:"output_dir"`
}

type UpdateProjectRequest struct {
	Name         string `json:"name"`
	RepoURL      string `json:"repo_url"`
	BuildCommand string `json:"build_command"`
	OutputDir    string `json:"output_dir"`
}


package domain

import "time"

type DeploymentStatus string

const (
	StatusPending   DeploymentStatus = "pending"
	StatusBuilding  DeploymentStatus = "building"
	StatusDeploying DeploymentStatus = "deploying"
	StatusReady     DeploymentStatus = "ready"
	StatusError     DeploymentStatus = "error"
)

type Deployment struct {
	ID         string           `json:"id"`
	ProjectID  string           `json:"project_id"`
	Status     DeploymentStatus `json:"status"`
	Subdomain  string           `json:"subdomain"`
	ImageURL   string           `json:"image_url"`
	CommitHash string           `json:"commit_hash"`
	BuildLogs  string           `json:"build_logs"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

type TriggerDeployRequest struct {
	ProjectID  string `json:"project_id" validate:"required"`
	CommitHash string `json:"commit_hash"`
}

type DeploymentEvent struct {
	DeploymentID string `json:"deployment_id"`
	ProjectID    string `json:"project_id"`
	RepoURL      string `json:"repo_url"`
	BuildCommand string `json:"build_command"`
	OutputDir    string `json:"output_dir"`
	CommitHash   string `json:"commit_hash"`
}

type BuildCompleteEvent struct {
	DeploymentID string `json:"deployment_id"`
	ImageURL     string `json:"image_url"`
	Success      bool   `json:"success"`
	Logs         string `json:"logs"`
}

type DeployCompleteEvent struct {
	DeploymentID string `json:"deployment_id"`
	Success      bool   `json:"success"`
	Subdomain    string `json:"subdomain"`
}

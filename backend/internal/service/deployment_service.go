package service

import (
	"errors"
	"fmt"

	"github.com/dejavu/backend/internal/domain"
	"github.com/dejavu/backend/internal/repository"
	"github.com/dejavu/backend/pkg/queue"
	"github.com/google/uuid"
)

type DeploymentService struct {
	deployRepo  *repository.DeploymentRepository
	projectRepo *repository.ProjectRepository
	queue       *queue.Queue
}

func NewDeploymentService(
	deployRepo *repository.DeploymentRepository,
	projectRepo *repository.ProjectRepository,
	queue *queue.Queue,
) *DeploymentService {
	return &DeploymentService{
		deployRepo:  deployRepo,
		projectRepo: projectRepo,
		queue:       queue,
	}
}

func (s *DeploymentService) Trigger(userID string, req *domain.TriggerDeployRequest) (*domain.Deployment, error) {
	// Verify project ownership
	project, err := s.projectRepo.GetByID(req.ProjectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}
	if project.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Generate subdomain
	subdomain := s.generateSubdomain()

	// Create deployment record
	deployment := &domain.Deployment{
		ProjectID:  req.ProjectID,
		Status:     domain.StatusPending,
		Subdomain:  subdomain,
		CommitHash: req.CommitHash,
	}

	if err := s.deployRepo.Create(deployment); err != nil {
		return nil, err
	}

	// Publish to NATS for builder
	event := domain.DeploymentEvent{
		DeploymentID: deployment.ID,
		ProjectID:    project.ID,
		RepoURL:      project.RepoURL,
		BuildCommand: project.BuildCommand,
		OutputDir:    project.OutputDir,
		CommitHash:   req.CommitHash,
	}

	if err := s.queue.Publish("DEPLOYMENTS.request", event); err != nil {
		return nil, err
	}

	return deployment, nil
}

func (s *DeploymentService) GetStatus(id string) (*domain.Deployment, error) {
	deployment, err := s.deployRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if deployment == nil {
		return nil, errors.New("deployment not found")
	}
	return deployment, nil
}

func (s *DeploymentService) ListByProject(projectID string) ([]*domain.Deployment, error) {
	return s.deployRepo.ListByProjectID(projectID)
}

func (s *DeploymentService) generateSubdomain() string {
	return fmt.Sprintf("app-%s", uuid.New().String()[:8])
}


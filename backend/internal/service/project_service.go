package service

import (
	"errors"

	"github.com/dejavu/backend/internal/domain"
	"github.com/dejavu/backend/internal/repository"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(userID string, req *domain.CreateProjectRequest) (*domain.Project, error) {
	buildCmd := req.BuildCommand
	if buildCmd == "" {
		buildCmd = "npm run build"
	}

	outputDir := req.OutputDir
	if outputDir == "" {
		outputDir = "dist"
	}

	project := &domain.Project{
		UserID:       userID,
		Name:         req.Name,
		RepoURL:      req.RepoURL,
		BuildCommand: buildCmd,
		OutputDir:    outputDir,
	}

	if err := s.repo.Create(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) GetByID(id, userID string) (*domain.Project, error) {
	project, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}
	if project.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	return project, nil
}

func (s *ProjectService) List(userID string) ([]*domain.Project, error) {
	return s.repo.ListByUserID(userID)
}

func (s *ProjectService) Update(id, userID string, req *domain.UpdateProjectRequest) error {
	project, err := s.GetByID(id, userID)
	if err != nil {
		return err
	}

	if req.Name != "" {
		project.Name = req.Name
	}
	if req.RepoURL != "" {
		project.RepoURL = req.RepoURL
	}
	if req.BuildCommand != "" {
		project.BuildCommand = req.BuildCommand
	}
	if req.OutputDir != "" {
		project.OutputDir = req.OutputDir
	}

	return s.repo.Update(project)
}

func (s *ProjectService) Delete(id, userID string) error {
	return s.repo.Delete(id, userID)
}

package service

import (
	"database/sql"
	"errors"

	"github.com/dejavu/backend/internal/repository"
)

type BillingService struct {
	repo *repository.BillingRepository
}

const (
	// Usage costs
	CostPerBuildMinute = 0.10  // $0.10 per minute
	CostPerGBBandwidth = 0.15  // $0.15 per GB
	CostPerDeployment  = 0.50  // $0.50 per deployment
)

func NewBillingService(repo *repository.BillingRepository) *BillingService {
	return &BillingService{repo: repo}
}

func (s *BillingService) GetBalance(userID string) (float64, error) {
	account, err := s.repo.GetOrCreateAccount(userID)
	if err != nil {
		return 0, err
	}
	return account.Credits, nil
}

func (s *BillingService) AddCredits(userID string, amount float64) error {
	return s.repo.AddCredits(userID, amount)
}

func (s *BillingService) ChargeForDeployment(userID, deploymentID string) error {
	// Check if user has sufficient credits
	balance, err := s.GetBalance(userID)
	if err != nil {
		return err
	}

	if balance < CostPerDeployment {
		return errors.New("insufficient credits")
	}

	// Deduct credits
	if err := s.repo.DeductCredits(userID, CostPerDeployment); err != nil {
		return err
	}

	// Record usage
	record := &repository.UsageRecord{
		UserID:       userID,
		DeploymentID: sql.NullString{String: deploymentID, Valid: true},
		Type:         "deployment",
		Amount:       CostPerDeployment,
		Description:  "Deployment charge",
	}
	return s.repo.RecordUsage(record)
}

func (s *BillingService) ChargeForBuildTime(userID, deploymentID string, minutes float64) error {
	cost := minutes * CostPerBuildMinute

	// Deduct credits
	if err := s.repo.DeductCredits(userID, cost); err != nil {
		return err
	}

	// Record usage
	record := &repository.UsageRecord{
		UserID:       userID,
		DeploymentID: sql.NullString{String: deploymentID, Valid: true},
		Type:         "build_time",
		Amount:       cost,
		Description:  "Build time charge",
	}
	return s.repo.RecordUsage(record)
}

func (s *BillingService) ChargeForBandwidth(userID string, gigabytes float64) error {
	cost := gigabytes * CostPerGBBandwidth

	// Deduct credits
	if err := s.repo.DeductCredits(userID, cost); err != nil {
		return err
	}

	// Record usage
	record := &repository.UsageRecord{
		UserID:      userID,
		Type:        "bandwidth",
		Amount:      cost,
		Description: "Bandwidth charge",
	}
	return s.repo.RecordUsage(record)
}

func (s *BillingService) GetUsageHistory(userID string) ([]*repository.UsageRecord, error) {
	return s.repo.GetUsageHistory(userID, 50)
}


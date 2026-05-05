package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"lakoo/backend/internal/domain"
	"lakoo/backend/internal/dto"
	"lakoo/backend/internal/repository"
	"lakoo/backend/pkg/auth"
	"lakoo/backend/pkg/config"
)

type TenantService interface {
	Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	ForgotPassword(req *dto.ForgotPasswordRequest) error
	ResetPassword(req *dto.ResetPasswordRequest) error
	UpdateProfile(tenantID string, userID string, req *dto.UpdateTenantRequest) error
	ChangePassword(userID string, req *dto.ChangePasswordRequest) error
	ListStaff(tenantID string) ([]*dto.StaffMemberResponse, error)
	AddStaff(tenantID string, req *dto.AddStaffRequest) error
	RemoveStaff(tenantID string, userID string) error
}

type tenantService struct {
	tenantRepo repository.TenantRepository
	userRepo   repository.UserRepository
	redisCli   *redis.Client
	cfg        *config.Config
}

func NewTenantService(tr repository.TenantRepository, ur repository.UserRepository, redisCli *redis.Client, cfg *config.Config) TenantService {
	return &tenantService{
		tenantRepo: tr,
		userRepo:   ur,
		redisCli:   redisCli,
		cfg:        cfg,
	}
}

func (u *tenantService) Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	existingTenant, _ := u.tenantRepo.GetBySlug(req.Slug)
	if existingTenant != nil {
		return nil, errors.New("slug is already in use")
	}

	existingUser, _ := u.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email is already in use")
	}

	newTenantID := uuid.New().String()
	newUserID := uuid.New().String()

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	trialEnd := time.Now().AddDate(0, 0, 14)
	now := time.Now()

	tenant := &domain.Tenant{
		ID:          newTenantID,
		Slug:        req.Slug,
		Name:        req.TenantName,
		Plan:        "free",
		Status:      "active",
		OwnerID:     newUserID,
		TrialEndsAt: &trialEnd,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	user := &domain.User{
		ID:        newUserID,
		TenantID:  newTenantID,
		Name:      req.OwnerName,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      "owner",
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = u.tenantRepo.Create(tenant)
	if err != nil {
		return nil, err
	}

	err = u.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		TenantID: tenant.ID,
		UserID:   user.ID,
		Slug:     tenant.Slug,
	}, nil
}

func (u *tenantService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !auth.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	tenant, err := u.tenantRepo.GetByID(user.TenantID)
	if err != nil || tenant == nil {
		return nil, errors.New("tenant not found for user")
	}

	token, err := auth.GenerateTokens(user.ID, user.TenantID, user.Role, u.cfg)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserObj{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
		Tenant: dto.TenantObj{
			ID:            tenant.ID,
			Name:          tenant.Name,
			Slug:          tenant.Slug,
			PaymentConfig: tenant.PaymentConfig,
			LogoURL:       tenant.LogoURL,
		},
	}, nil
}

func (u *tenantService) ForgotPassword(req *dto.ForgotPasswordRequest) error {
	user, err := u.userRepo.GetByEmail(req.Email)
	if err != nil || user == nil {
		// Silent drop to defend against email enumeration attacks
		return nil
	}

	token := uuid.New().String()
	ctx := context.Background()

	err = u.redisCli.Set(ctx, "reset_password:"+token, user.ID, 15*time.Minute).Err()
	if err != nil {
		return err
	}

	// Simulated Email Transporter
	log.Printf("[SIMULASI EMAIL] Klik tautan berikut untuk Reset Password: http://localhost:5174/reset-password?token=%s\n", token)
	return nil
}

func (u *tenantService) ResetPassword(req *dto.ResetPasswordRequest) error {
	ctx := context.Background()
	userID, err := u.redisCli.Get(ctx, "reset_password:"+req.Token).Result()
	if err != nil {
		return errors.New("Token pemulihan tidak valid atau sudah kedaluwarsa")
	}

	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	err = u.userRepo.UpdatePassword(userID, hashedPassword)
	if err != nil {
		return err
	}

	// Discard token
	_ = u.redisCli.Del(ctx, "reset_password:"+req.Token)
	return nil
}

func (u *tenantService) UpdateProfile(tenantID string, userID string, req *dto.UpdateTenantRequest) error {
	tenant, err := u.tenantRepo.GetByID(tenantID)
	if err != nil || tenant == nil { return errors.New("tenant not found") }

	existingSlug, _ := u.tenantRepo.GetBySlug(req.Slug)
	if existingSlug != nil && existingSlug.ID != tenantID { return errors.New("slug is already in use") }

	user, err := u.userRepo.GetByID(userID)
	if err != nil || user == nil { return errors.New("user not found") }

	existingEmail, _ := u.userRepo.GetByEmail(req.Email)
	if existingEmail != nil && existingEmail.ID != userID { return errors.New("email is already in use") }

	tenant.Name = req.Name
	tenant.Slug = req.Slug
	tenant.PaymentConfig = &req.PaymentConfig
	tenant.LogoURL = &req.LogoURL
	tenant.UpdatedAt = time.Now()
	if err := u.tenantRepo.Update(tenant); err != nil { return err }

	user.Name = req.OwnerName
	user.Email = req.Email
	user.UpdatedAt = time.Now()
	return u.userRepo.UpdateProfile(user)
}

func (u *tenantService) ChangePassword(userID string, req *dto.ChangePasswordRequest) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil || user == nil { return errors.New("user not found") }

	if !auth.CheckPasswordHash(req.OldPassword, user.Password) {
		return errors.New("password lama tidak sesuai")
	}

	hashed, err := auth.HashPassword(req.NewPassword)
	if err != nil { return err }

	return u.userRepo.UpdatePassword(userID, hashed)
}

func (u *tenantService) ListStaff(tenantID string) ([]*dto.StaffMemberResponse, error) {
	users, err := u.userRepo.GetByTenantID(tenantID)
	if err != nil { return nil, err }

	res := make([]*dto.StaffMemberResponse, 0)
	for _, user := range users {
		res = append(res, &dto.StaffMemberResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		})
	}
	return res, nil
}

func (u *tenantService) AddStaff(tenantID string, req *dto.AddStaffRequest) error {
	existing, _ := u.userRepo.GetByEmail(req.Email)
	if existing != nil { return errors.New("email is already in use") }

	hashed, err := auth.HashPassword(req.Password)
	if err != nil { return err }

	now := time.Now()
	user := &domain.User{
		ID:        uuid.New().String(),
		TenantID:  tenantID,
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashed,
		Role:      req.Role,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return u.userRepo.Create(user)
}

func (u *tenantService) RemoveStaff(tenantID string, userID string) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil || user == nil { return errors.New("user not found") }

	if user.TenantID != tenantID { return errors.New("unauthorized") }
	if user.Role == "owner" { return errors.New("cannot remove owner") }

	return u.userRepo.Delete(userID)
}

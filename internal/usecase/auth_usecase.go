package usecase

import (
	"context"

	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/internal/repository"
	"github.com/template-be-commerce/pkg/apperrors"
	jwtpkg "github.com/template-be-commerce/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	RegisterCustomer(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error)
	RegisterAdmin(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error)
	LoginCustomer(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error)
	LoginAdmin(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error)
	RefreshToken(ctx context.Context, req domain.RefreshTokenRequest) (*domain.AuthResponse, error)
}

type authUseCase struct {
	userRepo     repository.UserRepository
	customerRepo repository.CustomerRepository
	jwtManager   *jwtpkg.Manager
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	customerRepo repository.CustomerRepository,
	jwtManager *jwtpkg.Manager,
) AuthUseCase {
	return &authUseCase{
		userRepo:     userRepo,
		customerRepo: customerRepo,
		jwtManager:   jwtManager,
	}
}

func (uc *authUseCase) RegisterCustomer(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check email availability
	existing, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, apperrors.ErrEmailExists
	}
	if err != nil && !apperrors.Is(err, 404) {
		return nil, err
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.InternalServer("failed to hash password", err)
	}

	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         domain.RoleCustomer,
		IsActive:     true,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	customer := &domain.Customer{
		UserID:    user.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	}
	if err := uc.customerRepo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return uc.buildAuthResponse(user)
}

func (uc *authUseCase) RegisterAdmin(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check email availability
	existing, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, apperrors.ErrEmailExists
	}
	if err != nil && !apperrors.Is(err, 404) {
		return nil, err
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.InternalServer("failed to hash password", err)
	}

	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         domain.RoleAdmin,
		IsActive:     true,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Note: Admin users don't have a customer record

	return uc.buildAuthResponse(user)
}

func (uc *authUseCase) LoginCustomer(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	return uc.login(ctx, req, domain.RoleCustomer)
}

func (uc *authUseCase) LoginAdmin(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	return uc.login(ctx, req, domain.RoleAdmin)
}

func (uc *authUseCase) login(ctx context.Context, req domain.LoginRequest, expectedRole domain.UserRole) (*domain.AuthResponse, error) {
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.ErrInvalidPassword
	}

	if user.Role != expectedRole {
		return nil, apperrors.ErrInvalidPassword
	}

	if !user.IsActive {
		return nil, apperrors.ErrInactiveAccount
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrInvalidPassword
	}

	return uc.buildAuthResponse(user)
}

func (uc *authUseCase) RefreshToken(ctx context.Context, req domain.RefreshTokenRequest) (*domain.AuthResponse, error) {
	claims, err := uc.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, apperrors.ErrInvalidToken
	}

	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, apperrors.ErrInvalidToken
	}

	if !user.IsActive {
		return nil, apperrors.ErrInactiveAccount
	}

	return uc.buildAuthResponse(user)
}

func (uc *authUseCase) buildAuthResponse(user *domain.User) (*domain.AuthResponse, error) {
	tokenPair, err := uc.jwtManager.GenerateTokenPair(user)
	if err != nil {
		return nil, apperrors.InternalServer("failed to generate tokens", err)
	}

	return &domain.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		User: domain.UserInfo{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

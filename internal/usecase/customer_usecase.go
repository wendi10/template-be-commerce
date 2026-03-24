package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/internal/repository"
	"github.com/template-be-commerce/pkg/apperrors"
)

type CustomerUseCase interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.Customer, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req domain.UpdateProfileRequest) (*domain.Customer, error)
	ListAddresses(ctx context.Context, userID uuid.UUID) ([]domain.Address, error)
	CreateAddress(ctx context.Context, userID uuid.UUID, req domain.CreateAddressRequest) (*domain.Address, error)
	UpdateAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID, req domain.UpdateAddressRequest) (*domain.Address, error)
	DeleteAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) error
	SetDefaultAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) error
	// Admin
	ListAll(ctx context.Context, page, limit int) ([]domain.Customer, int64, error)
}

type customerUseCase struct {
	customerRepo repository.CustomerRepository
	addressRepo  repository.AddressRepository
}

func NewCustomerUseCase(customerRepo repository.CustomerRepository, addressRepo repository.AddressRepository) CustomerUseCase {
	return &customerUseCase{customerRepo: customerRepo, addressRepo: addressRepo}
}

func (uc *customerUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.Customer, error) {
	customer, err := uc.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	addresses, err := uc.addressRepo.ListByCustomerID(ctx, customer.ID)
	if err != nil {
		return nil, err
	}
	customer.Addresses = addresses
	return customer, nil
}

func (uc *customerUseCase) UpdateProfile(ctx context.Context, userID uuid.UUID, req domain.UpdateProfileRequest) (*domain.Customer, error) {
	customer, err := uc.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.FirstName != "" {
		customer.FirstName = req.FirstName
	}
	if req.LastName != "" {
		customer.LastName = req.LastName
	}
	if req.Phone != "" {
		customer.Phone = req.Phone
	}
	if req.Avatar != "" {
		customer.Avatar = req.Avatar
	}

	if err := uc.customerRepo.Update(ctx, customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (uc *customerUseCase) ListAddresses(ctx context.Context, userID uuid.UUID) ([]domain.Address, error) {
	customer, err := uc.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return uc.addressRepo.ListByCustomerID(ctx, customer.ID)
}

func (uc *customerUseCase) CreateAddress(ctx context.Context, userID uuid.UUID, req domain.CreateAddressRequest) (*domain.Address, error) {
	customer, err := uc.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	address := &domain.Address{
		CustomerID:    customer.ID,
		Label:         req.Label,
		RecipientName: req.RecipientName,
		Phone:         req.Phone,
		AddressLine1:  req.AddressLine1,
		AddressLine2:  req.AddressLine2,
		City:          req.City,
		Province:      req.Province,
		PostalCode:    req.PostalCode,
		IsDefault:     req.IsDefault,
	}

	if err := uc.addressRepo.Create(ctx, address); err != nil {
		return nil, err
	}

	// If this is the first address or marked as default, set it
	if req.IsDefault {
		_ = uc.addressRepo.SetDefault(ctx, customer.ID, address.ID)
	}

	return address, nil
}

func (uc *customerUseCase) UpdateAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID, req domain.UpdateAddressRequest) (*domain.Address, error) {
	customer, err := uc.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	address, err := uc.addressRepo.GetByID(ctx, addressID)
	if err != nil {
		return nil, err
	}

	// Ensure ownership
	if address.CustomerID != customer.ID {
		return nil, apperrors.ErrForbidden
	}

	if req.Label != "" {
		address.Label = req.Label
	}
	if req.RecipientName != "" {
		address.RecipientName = req.RecipientName
	}
	if req.Phone != "" {
		address.Phone = req.Phone
	}
	if req.AddressLine1 != "" {
		address.AddressLine1 = req.AddressLine1
	}
	if req.AddressLine2 != "" {
		address.AddressLine2 = req.AddressLine2
	}
	if req.City != "" {
		address.City = req.City
	}
	if req.Province != "" {
		address.Province = req.Province
	}
	if req.PostalCode != "" {
		address.PostalCode = req.PostalCode
	}
	if req.IsDefault != nil {
		address.IsDefault = *req.IsDefault
	}

	if err := uc.addressRepo.Update(ctx, address); err != nil {
		return nil, err
	}

	if address.IsDefault {
		_ = uc.addressRepo.SetDefault(ctx, customer.ID, address.ID)
	}

	return address, nil
}

func (uc *customerUseCase) DeleteAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) error {
	customer, err := uc.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	address, err := uc.addressRepo.GetByID(ctx, addressID)
	if err != nil {
		return err
	}

	if address.CustomerID != customer.ID {
		return apperrors.ErrForbidden
	}

	return uc.addressRepo.Delete(ctx, addressID)
}

func (uc *customerUseCase) ListAll(ctx context.Context, page, limit int) ([]domain.Customer, int64, error) {
	return uc.customerRepo.List(ctx, page, limit)
}

func (uc *customerUseCase) SetDefaultAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) error {
	customer, err := uc.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	address, err := uc.addressRepo.GetByID(ctx, addressID)
	if err != nil {
		return err
	}

	if address.CustomerID != customer.ID {
		return apperrors.ErrForbidden
	}

	return uc.addressRepo.SetDefault(ctx, customer.ID, addressID)
}

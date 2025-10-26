package repository

import (
	"errors"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"grubzo/internal/utils/ce"
	"grubzo/internal/utils/random"

	"gorm.io/gorm"
)

type TenantUserRepository interface {
	CreateTenantUser(dto *dto.CreateTenantUser) (*entity.TenantUser, error)
	UpdateTenantUser(dto *dto.UpdateTenantUser) (*entity.TenantUser, error)
	FindTenantUser(query *query.TenantUserQuery) (*entity.TenantUser, error)
	FindAllTenantUsers(query *query.TenantUserQuery) ([]*entity.TenantUser, error)
	CheckTenantUserPassword(user *entity.TenantUser, password string) bool
}

func populateHash(usr *entity.TenantUser) {
	salt := random.SecureAlphaNumeric(16)
	pass := random.HashPassword(usr.Password, salt)
	usr.Password = pass
	usr.Salt = salt
}

func validateTenantUser(usr *entity.TenantUser, db *gorm.DB) error {
	// Check unique email across tenant
	var count int64
	sess := db.Session(&gorm.Session{}).Model(entity.TenantUser{})
	sess.Where("email = ? AND tenant_id = ?", usr.Email, usr.TenantID)
	if usr.ID != 0 {
		sess.Not("id = ?", usr.ID)
	}
	sess.Count(&count)
	if count > 0 {
		return ce.New("User with same email already exists")
	}
	return nil
}

func (r *Repository) CreateTenantUser(dto *dto.CreateTenantUser) (*entity.TenantUser, error) {
	tenantUser := &entity.TenantUser{
		TenantID:   dto.TenantID,
		Email:      dto.Email,
		Password:   dto.Password,
		Name:       dto.Name,
		LocationID: dto.LocationID,
	}
	if dto.Role != nil {
		tenantUser.Role = *dto.Role
	}

	if err := validateTenantUser(tenantUser, r.db); err != nil {
		return nil, err
	}
	populateHash(tenantUser)
	if err := r.db.Create(tenantUser).Error; err != nil {
		return nil, err
	}

	return tenantUser, nil
}

func (r *Repository) UpdateTenantUser(dto *dto.UpdateTenantUser) (*entity.TenantUser, error) {
	tenantUser, err := r.FindTenantUser(query.NewTenantUserQuery(dto.TenantID).WithID(dto.ID))
	if err != nil {
		return nil, err
	}

	if dto.Email != nil {
		tenantUser.Email = *dto.Email
	}
	if dto.Password != nil {
		tenantUser.Password = *dto.Password
		populateHash(tenantUser)
	}
	if dto.Name != nil {
		tenantUser.Name = *dto.Name
	}
	if dto.Role != nil {
		tenantUser.Role = *dto.Role
	}
	if dto.LocationID != nil {
		tenantUser.LocationID = dto.LocationID
	}

	if err := r.db.Save(&tenantUser).Error; err != nil {
		return nil, err
	}

	return tenantUser, nil
}

func (r *Repository) FindTenantUser(q *query.TenantUserQuery) (*entity.TenantUser, error) {
	tenantUser := &entity.TenantUser{}
	sess := r.db.Session(&gorm.Session{}).Model(&entity.TenantUser{})

	sess = sess.Where("tenant_id = ?", q.TenantID)
	if q.WithPreload {
		for _, preload := range tenantUser.GetPreloads() {
			sess = sess.Preload(preload)
		}
	}
	if q.ID != nil {
		sess = sess.Where("id = ?", *q.ID)
	}
	if q.Email != nil {
		sess = sess.Where("email = ?", *q.Email)
	}
	if q.Role != nil {
		sess = sess.Where("role = ?", *q.Role)
	}

	if err := sess.First(&tenantUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ce.New("User with the specified ID was not found")
		}
		return nil, err
	}
	return tenantUser, nil
}

func (r *Repository) FindAllTenantUsers(q *query.TenantUserQuery) ([]*entity.TenantUser, error) {
	var tenantUsers []*entity.TenantUser
	sess := r.db.Session(&gorm.Session{}).Model(&entity.TenantUser{})

	sess = sess.Where("tenant_id = ?", q.TenantID)
	if q.ID != nil {
		sess = sess.Where("id = ?", *q.ID)
	}
	if q.Email != nil {
		sess = sess.Where("email = ?", *q.Email)
	}
	if q.Role != nil {
		sess = sess.Where("role = ?", *q.Role)
	}

	if err := sess.Find(&tenantUsers).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ce.New("User with the specified details was not found")
		}
		return nil, err
	}
	return tenantUsers, nil
}

func (r *Repository) CheckTenantUserPassword(user *entity.TenantUser, password string) bool {
	return random.HashPassword(password, user.Salt) == user.Password
}

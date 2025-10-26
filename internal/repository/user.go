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

type UserRepository interface {
	CreateUser(dto *dto.CreateUser) (*entity.User, error)
	UpdateUser(dto *dto.UpdateUser) (*entity.User, error)
	FindUser(query *query.UserQuery) (*entity.User, error)
	FindAllUsers(query *query.UserQuery) ([]*entity.User, error)
	CheckUserPassword(user *entity.User, password string) bool
}

func populateUserHash(usr *entity.User) {
	salt := random.SecureAlphaNumeric(16)
	pass := random.HashPassword(usr.Password, salt)
	usr.Password = pass
	usr.Salt = salt
}

func validateUser(usr *entity.User, db *gorm.DB) error {
	// Check unique email across tenant
	var count int64
	sess := db.Session(&gorm.Session{}).Model(entity.User{})
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

func (r *Repository) CreateUser(dto *dto.CreateUser) (*entity.User, error) {
	user := &entity.User{
		TenantID: dto.TenantID,
		UserID:   dto.UserID,
		Email:    dto.Email,
		Password: dto.Password,
		Name:     dto.Name,
	}

	if err := validateUser(user, r.db); err != nil {
		return nil, err
	}
	populateUserHash(user)
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

const (
	UserNotFound = "User with provided details not found"
)

func (r *Repository) CheckUserPassword(user *entity.User, password string) bool {
	return random.HashPassword(password, user.Salt) == user.Password
}

func (r *Repository) UpdateUser(dto *dto.UpdateUser) (*entity.User, error) {
	user, err := r.FindUser(query.NewUserQuery(dto.TenantID).WithID(dto.ID))
	if err != nil {
		return nil, err
	}
	if dto.Email != nil {
		user.Email = *dto.Email
	}
	if dto.Password != nil {
		user.Password = *dto.Password
		populateUserHash(user)
	}
	if dto.Name != nil {
		user.Name = *dto.Name
	}
	if err := validateUser(user, r.db); err != nil {
		return nil, err
	}
	if err := r.db.Save(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) FindUser(filter *query.UserQuery) (*entity.User, error) {
	user := &entity.User{}
	sess := r.db.Session(&gorm.Session{}).Model(&entity.User{})

	sess = sess.Where("tenant_id = ?", filter.TenantID)
	if filter.ID != nil {
		sess = sess.Where("id = ?", *filter.ID)
	}
	if filter.UserId != nil {
		sess = sess.Where("user_id = ?", *filter.UserId)
	}
	if filter.Email != nil {
		sess = sess.Where("email = ?", *filter.Email)
	}

	if err := sess.First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ce.New(UserNotFound)
		}
		return nil, err
	}
	return user, nil
}

func (r *Repository) FindAllUsers(filter *query.UserQuery) ([]*entity.User, error) {
	var users []*entity.User
	sess := r.db.Session(&gorm.Session{}).Model(&entity.User{})

	sess.Where("tenant_id = ?", filter.TenantID)
	if filter.ID != nil {
		sess = sess.Where("id = ?", *filter.ID)
	}
	if filter.UserId != nil {
		sess = sess.Where("user_id = ?", *filter.UserId)
	}
	if filter.Email != nil {
		sess = sess.Where("email = ?", *filter.Email)
	}

	if err := sess.Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ce.New("User with provided details not found")
		}
		return nil, err
	}
	return users, nil
}

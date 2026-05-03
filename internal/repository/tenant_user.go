package repository

//go:generate go run ../../cmd/injecttrace -file tenant_user.go -receiver Repository -service Repository
import (
	"context"
	"errors"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"grubzo/internal/router/ext"
	"grubzo/internal/utils/random"

	"github.com/lib/pq"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TenantUserRepository interface {
	CreateTenantUser(ctx context.Context, dto *dto.CreateTenantUser) (*entity.TenantUser, error)
	UpdateTenantUser(ctx context.Context, dto *dto.UpdateTenantUser) (*entity.TenantUser, error)
	FindTenantUser(ctx context.Context, query *query.TenantUserQuery) (*entity.TenantUser, error)
	FindAllTenantUsers(ctx context.Context, query *query.TenantUserQuery) ([]*entity.TenantUser, error)
	CheckTenantUserPassword(user *entity.TenantUser, password string) bool
}

func populateHash(usr *entity.TenantUser) {
	salt := random.SecureAlphaNumeric(16)
	pass := random.HashPassword(usr.Password, salt)
	usr.Password = pass
	usr.Salt = salt
}

func (r *Repository) validateTenantUser(usr *entity.TenantUser, db *gorm.DB) error {
	// check unique email across tenant
	var count int64
	sess := db.Session(&gorm.Session{}).Model(entity.TenantUser{})
	sess.Where("email = ? AND tenant_id = ?", usr.Email, usr.TenantID)
	if usr.ID != 0 {
		sess.Not("id = ?", usr.ID)
	}
	sess.Count(&count)
	if count > 0 {
		return ext.Error("User with same email already exists")
	}

	// check roles provided exists
	sess = db.Session(&gorm.Session{}).Model(entity.UserRole{})
	err := sess.Where("tenant_id = ? AND name=ANY(?)", usr.TenantID, usr.Roles).Count(&count).Error
	if err != nil {
		return err
	}
	if count != int64(len(usr.Roles)) {
		r.logger.Debug("COUNT, USER_ROLES_COUNT", zap.Any("COUNT", count), zap.Any("USER_ROLES_COUNT", len(usr.Roles)))
		return ext.Error("Invalid user role")
	}

	return nil
}

func (r *Repository) CreateTenantUser(ctx context.Context, dto *dto.CreateTenantUser) (*entity.TenantUser, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.CreateTenantUser")
	defer span.End()

	tenantUser := &entity.TenantUser{
		TenantID:   dto.TenantID,
		Email:      dto.Email,
		Password:   dto.Password,
		Name:       dto.Name,
		LocationID: dto.LocationID,
	}
	if len(dto.Roles) != 0 {
		tenantUser.Roles = pq.StringArray(dto.Roles)
	}

	db := r.dbWithContext(ctx)
	if err := r.validateTenantUser(tenantUser, db); err != nil {
		return nil, err
	}
	populateHash(tenantUser)
	if err := db.Create(tenantUser).Error; err != nil {
		return nil, err
	}

	return tenantUser, nil
}

func (r *Repository) UpdateTenantUser(ctx context.Context, dto *dto.UpdateTenantUser) (*entity.TenantUser, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.UpdateTenantUser")
	defer span.End()

	tenantUser, err := r.FindTenantUser(ctx, query.NewTenantUserQuery(dto.TenantID).WithID(dto.ID))
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
	if len(dto.Roles) != 0 {
		tenantUser.Roles = pq.StringArray(dto.Roles)
	}
	if dto.LocationID != nil {
		tenantUser.LocationID = *dto.LocationID
	}

	if err := r.dbWithContext(ctx).Save(&tenantUser).Error; err != nil {
		return nil, err
	}

	return tenantUser, nil
}

func (r *Repository) FindTenantUser(ctx context.Context, q *query.TenantUserQuery) (*entity.TenantUser, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.FindTenantUser")
	defer span.End()

	tenantUser := &entity.TenantUser{}
	sess := r.dbWithContext(ctx).Session(&gorm.Session{}).Model(&entity.TenantUser{})

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
	if len(q.Roles) != 0 {
		sess = sess.Where("role IN ?", q.Roles)
	}

	if err := sess.First(&tenantUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ext.Error("User with the specified ID was not found")
		}
		return nil, err
	}
	return tenantUser, nil
}

func (r *Repository) FindAllTenantUsers(ctx context.Context, q *query.TenantUserQuery) ([]*entity.TenantUser, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.FindAllTenantUsers")
	defer span.End()

	var tenantUsers []*entity.TenantUser
	sess := r.dbWithContext(ctx).Session(&gorm.Session{}).Model(&entity.TenantUser{})

	sess = sess.Where("tenant_id = ?", q.TenantID)
	if q.ID != nil {
		sess = sess.Where("id = ?", *q.ID)
	}
	if q.Email != nil {
		sess = sess.Where("email = ?", *q.Email)
	}
	if len(q.Roles) != 0 {
		sess = sess.Where("role IN ?", pq.StringArray(q.Roles))
	}
	if q.LocationID != nil {
		sess = sess.Where("location_id = ?", q.LocationID)
	}

	if err := sess.Find(&tenantUsers).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ext.Error("User with the specified details was not found")
		}
		return nil, err
	}
	return tenantUsers, nil
}

func (r *Repository) CheckTenantUserPassword(user *entity.TenantUser, password string) bool {
	return random.HashPassword(password, user.Salt) == user.Password
}

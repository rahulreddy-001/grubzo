package store

import (
	"grubzo/internal/config"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"
	"grubzo/internal/repository"
	"grubzo/internal/services/file"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type StoreService interface {
	CreateItem(*dto.CreateMenuItem) (*dto.CreateMenuItemResponse, error)
	UpdateItem(*dto.UpdateMenuItem) (*dto.UpdateMenuItemResponse, error)
	GetItem(*query.MenuItemQuery) (*dto.GetMenuItemResponse, error)
	GetItems(*query.MenuItemQuery) (*dto.GetMenuItemsResponse, error)
}

func Init(repository *repository.Repository, config *config.Config, fm file.Manager, logger *zap.Logger) (*storeServiceImpl, error) {
	return &storeServiceImpl{
		repository:  repository,
		config:      config,
		fileManager: fm,
		logger:      logger.Named("store_service"),
	}, nil
}

type storeServiceImpl struct {
	repository  *repository.Repository
	config      *config.Config
	fileManager file.Manager
	logger      *zap.Logger
}

func (ss *storeServiceImpl) GetItem(query *query.MenuItemQuery) (*dto.GetMenuItemResponse, error) {
	eItem, err := ss.repository.GetItem(query)
	if err != nil {
		return nil, err
	}
	item := dto.MenuItem{}
	if copier.Copy(&item, eItem) != nil {
		ss.logger.Error("[copier.Copy] failed to copy eItem to item", zap.Any("eItem", eItem), zap.Any("item", item))
	}
	item.Files = ss.fileManager.GetFileMetas(eItem.Files)
	return &dto.GetMenuItemResponse{
		Message:  `Item created successfully`,
		MenuItem: item,
	}, nil
}

func (ss *storeServiceImpl) GetItems(query *query.MenuItemQuery) (*dto.GetMenuItemsResponse, error) {
	eItems, err := ss.repository.GetItems(query)
	if err != nil {
		return nil, err
	}
	items := make([]dto.MenuItem, len(eItems))
	for i, eItem := range eItems {
		if copier.Copy(&items[i], eItem) != nil {
			ss.logger.Error("[copier.Copy] failed to copy eItem to item", zap.Any("eItem", eItem), zap.Any("item", items[i]))
		}
		items[i].Files = ss.fileManager.GetFileMetas(eItem.Files)
	}
	return &dto.GetMenuItemsResponse{
		MenuItems: items,
		Message:   `Items fetched successfully`,
	}, nil
}

func (ss *storeServiceImpl) CreateItem(args *dto.CreateMenuItem) (*dto.CreateMenuItemResponse, error) {
	eItem, err := ss.repository.CreateItem(args)
	if err != nil {
		return nil, err
	}
	item := dto.MenuItem{}
	if copier.Copy(&item, eItem) != nil {
		ss.logger.Error("[copier.Copy] failed to copy eItem to item", zap.Any("eItem", eItem), zap.Any("item", item))
	}
	item.Files = ss.fileManager.GetFileMetas(eItem.Files)
	return &dto.CreateMenuItemResponse{
		MenuItem: item,
		Message:  `Item created successfully`,
	}, nil
}

func (ss *storeServiceImpl) UpdateItem(args *dto.UpdateMenuItem) (*dto.UpdateMenuItemResponse, error) {
	eItem, err := ss.repository.UpdateItem(args)
	if err != nil {
		return nil, err
	}
	item := dto.MenuItem{}
	if copier.Copy(&item, eItem) != nil {
		ss.logger.Error("[copier.Copy] failed to copy eItem to item", zap.Any("eItem", eItem), zap.Any("item", item))
	}
	item.Files = ss.fileManager.GetFileMetas(eItem.Files)
	return &dto.UpdateMenuItemResponse{
		MenuItem: item,
		Message:  `Item updated successfully`,
	}, nil
}

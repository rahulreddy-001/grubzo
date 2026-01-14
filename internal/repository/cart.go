package repository

import (
	"context"
	"encoding/json"
	"grubzo/internal/models/dto"

	"go.uber.org/zap"
)

type CartRepository interface {
	GetCart(key string) *dto.Cart
	SetItemQuantity(key string, action *dto.UpdateItemQuantity) (*dto.Cart, error)
	ClearCart(key string) *dto.Cart
}

func (repo *Repository) SetCart(cart *dto.Cart) bool {
	data, _ := json.Marshal(cart)
	_, err := repo.rdb.Do(context.Background(), "JSON.SET", cart.Key, ".", string(data)).Result()
	if err != nil {
		repo.logger.Error("failed to set cart in redis", zap.Error(err))
		return false
	}
	return true
}

func (repo *Repository) GetCart(key string) *dto.Cart {
	res, err := repo.rdb.Do(context.Background(), "JSON.GET", key).Result()
	if err != nil || res == nil {
		return &dto.Cart{
			Key:   key,
			Items: []dto.Item{},
		}
	}

	var cart dto.Cart
	json.Unmarshal([]byte(res.(string)), &cart)

	return &cart
}

func (repo *Repository) SetItemQuantity(key string, action *dto.UpdateItemQuantity) (*dto.Cart, error) {
	cart := repo.GetCart(key)

	updatedItems := []dto.Item{}
	found := false
	for _, item := range cart.Items {
		if item.Item == action.Item {
			found = true
			if *action.Quantity > 0 {
				updatedItems = append(updatedItems, dto.Item{
					Item:     action.Item,
					Quantity: *action.Quantity,
				})
			}
		} else {
			updatedItems = append(updatedItems, item)
		}
	}

	if !found && *action.Quantity > 0 {
		updatedItems = append(updatedItems, dto.Item{
			Item:     action.Item,
			Quantity: *action.Quantity,
		})
	}
	updatedCart := dto.Cart{
		Key:   key,
		Items: updatedItems,
	}
	data, _ := json.Marshal(updatedCart)
	_, err := repo.rdb.Do(context.Background(), "JSON.SET", key, ".", string(data)).Result()
	if err != nil {
		return nil, err
	}
	return &updatedCart, nil
}

func (repo *Repository) ClearCart(key string) *dto.Cart {
	repo.rdb.Del(context.Background(), key)
	return &dto.Cart{
		Key:   key,
		Items: []dto.Item{},
	}
}

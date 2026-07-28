package repository

//go:generate go run ../../cmd/injecttrace -file cart.go -receiver Repository -service Repository
import (
	"context"
	"encoding/json"
	"grubzo/internal/models/dto"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type CartRepository interface {
	SetCart(ctx context.Context, cart *dto.Cart) bool
	GetCart(ctx context.Context, key string) *dto.Cart
	SetItemQuantity(ctx context.Context, key string, action *dto.UpdateItemQuantity) (*dto.Cart, error)
	ClearCart(ctx context.Context, key string) *dto.Cart
}

func (repo *Repository) SetCart(ctx context.Context, cart *dto.Cart) bool {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.SetCart")
	defer span.End()

	data, _ := json.Marshal(cart)
	_, err := repo.rdb.Do(ctx, "JSON.SET", cart.Key, ".", string(data)).Result()
	if err != nil {
		repo.logger.Error("failed to set cart in redis", zap.Error(err))
		return false
	}
	return true
}

func (repo *Repository) GetCart(ctx context.Context, key string) *dto.Cart {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.GetCart")
	defer span.End()

	res, err := repo.rdb.Do(ctx, "JSON.GET", key).Result()
	if err != nil || res == nil {
		return &dto.Cart{
			Key:   key,
			Items: []dto.Item{},
		}
	}

	var cart dto.Cart
	if err := json.Unmarshal([]byte(res.(string)), &cart); err != nil {
		repo.logger.Error("failed to decode cart from redis", zap.Error(err))
		return &dto.Cart{
			Key:   key,
			Items: []dto.Item{},
		}
	}

	return &cart
}

func (repo *Repository) SetItemQuantity(ctx context.Context, key string, action *dto.UpdateItemQuantity) (*dto.Cart, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.SetItemQuantity")
	defer span.End()

	cart := repo.GetCart(ctx, key)

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
	_, err := repo.rdb.Do(ctx, "JSON.SET", key, ".", string(data)).Result()
	if err != nil {
		return nil, err
	}
	return &updatedCart, nil
}

func (repo *Repository) ClearCart(ctx context.Context, key string) *dto.Cart {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.ClearCart")
	defer span.End()

	repo.rdb.Del(ctx, key)
	return &dto.Cart{
		Key:   key,
		Items: []dto.Item{},
	}
}

package service

import (
    "context"
    "encoding/json" // Добавляем импорт
    "vet-app-clinic-management/internal/models"
    "vet-app-clinic-management/internal/repository/mySQL"
    "vet-app-clinic-management/internal/repository/redis"
)

type InventoryService struct {
    repo  *mySQL.InventoryRepo
    cache *redis.InventoryCache
}

func NewInventoryService(repo *mySQL.InventoryRepo, cache *redis.InventoryCache) *InventoryService {
    return &InventoryService{repo: repo, cache: cache}
}

func (s *InventoryService) GetAll(ctx context.Context) ([]models.Inventory, error) {
    // Проверяем кеш
    if data, err := s.cache.Get(ctx, "inventory"); err == nil {
        var inventory []models.Inventory
        json.Unmarshal([]byte(data), &inventory)
        return inventory, nil
    }

    // Если кеша нет, запрашиваем из MySQL
    inventory, err := s.repo.GetAll()
    if err != nil {
        return nil, err
    }

    // Сохраняем в кеш
    s.cache.Set(ctx, "inventory", inventory)
    return inventory, nil
}

func (s *InventoryService) UpdateQuantity(id uint, quantity int) error {
    return s.repo.UpdateQuantity(id, quantity)
}
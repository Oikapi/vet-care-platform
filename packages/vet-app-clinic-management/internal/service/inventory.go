package service

import (
    "context"
    "encoding/json"
    "vet-app-clinic-management/internal/models"
    "vet-app-clinic-management/internal/repository/mySQL"
    "vet-app-clinic-management/internal/repository/redis"
    "strconv"
    "fmt"
)

type InventoryService struct {
    repo  *mySQL.InventoryRepo
    cache *redis.InventoryCache
}

func NewInventoryService(repo *mySQL.InventoryRepo, cache *redis.InventoryCache) *InventoryService {
    return &InventoryService{repo: repo, cache: cache}
}

func (s *InventoryService) GetAll(ctx context.Context) ([]models.Inventory, error) {
    // Проверяем кэш
    if s.cache != nil {
        data, err := s.cache.Get(ctx, "inventory")
        if err == nil && data != nil {
            dataStr, ok := data.(string)
            if !ok {
                return nil, fmt.Errorf("cached data is not a string")
            }
            if dataStr != "" {
                var inventory []models.Inventory
                if err := json.Unmarshal([]byte(dataStr), &inventory); err == nil {
                    return inventory, nil
                }
            }
        }
    }

    // Если кэша нет, запрашиваем из MySQL
    inventory, err := s.repo.GetAll(ctx)
    if err != nil {
        return nil, err
    }

    // Сохраняем в кэш (сериализуем в JSON)
    if s.cache != nil {
        inventoryBytes, err := json.Marshal(inventory)
        if err == nil {
            if err := s.cache.Set(ctx, "inventory", string(inventoryBytes)); err != nil {
                return nil, err
            }
        }
    }

    return inventory, nil
}

func (s *InventoryService) GetByID(ctx context.Context, medicineID int) (*models.Inventory, error) {
    // Проверяем кэш
    if s.cache != nil {
        data, err := s.cache.Get(ctx, strconv.Itoa(medicineID))
        if err == nil && data != "" {
            dataStr, ok := data.(string)
            if !ok {
                return nil, fmt.Errorf("cached data is not a string")
            }
            var inventory models.Inventory
            if err := json.Unmarshal([]byte(dataStr), &inventory); err == nil {
                return &inventory, nil
            }
        }
    }

    // Если кэша нет, запрашиваем из базы
    inventory, err := s.repo.GetByID(medicineID)
    if err != nil {
        return nil, err
    }

    // Обновляем кэш
    if s.cache != nil {
        inventoryBytes, err := json.Marshal(inventory)
        if err == nil {
            if err := s.cache.Set(ctx, strconv.Itoa(medicineID), string(inventoryBytes)); err != nil {
                return nil, err
            }
        }
    }

    return inventory, nil
}

func (s *InventoryService) AutoOrder(ctx context.Context) error {
    inventory, err := s.GetAll(ctx)
    if err != nil {
        return err
    }

    ordered := false
    for _, item := range inventory {
        if item.Quantity < item.Threshold {
            // Здесь можно добавить логику реального заказа
            if s.cache != nil {
                if err := s.cache.Set(ctx, "order_"+item.MedicineName, "ordered"); err != nil {
                    return err
                }
            }
            ordered = true
        }
    }

    if ordered {
        // Инвалидация кэша после заказа
        if s.cache != nil {
            if err := s.cache.Del(ctx, "inventory"); err != nil {
                return err
            }
        }
    }

    return nil
}

func (s *InventoryService) UpdateQuantity(ctx context.Context, medicineID, quantity int) (*models.Inventory, error) {
    inventory, err := s.repo.GetByID(medicineID)
    if err != nil {
        return nil, err
    }
    inventory.Quantity = quantity
    if err := s.repo.Update(inventory); err != nil {
        return nil, err
    }

    if s.cache != nil {
        inventoryBytes, err := json.Marshal(inventory)
        if err == nil {
            if err := s.cache.Set(ctx, strconv.Itoa(medicineID), string(inventoryBytes)); err != nil {
                return nil, err
            }
            // Инвалидация общего кэша
            if err := s.cache.Del(ctx, "inventory"); err != nil {
                return nil, err
            }
        }
    }

    return inventory, nil
}

func (s *InventoryService) Create(ctx context.Context, inventory *models.Inventory) error {
    if err := s.repo.Create(inventory); err != nil {
        return err
    }
    // Инвалидация кэша
    if s.cache != nil {
        if err := s.cache.Del(ctx, "inventory"); err != nil {
            return err
        }
    }
    return nil
}

func (s *InventoryService) Delete(ctx context.Context, medicineID int) error {
    if err := s.repo.Delete(medicineID); err != nil {
        return err
    }
    // Инвалидация кэша
    if s.cache != nil {
        if err := s.cache.Del(ctx, strconv.Itoa(medicineID)); err != nil {
            return err
        }
        if err := s.cache.Del(ctx, "inventory"); err != nil {
            return err
        }
    }
    return nil
}
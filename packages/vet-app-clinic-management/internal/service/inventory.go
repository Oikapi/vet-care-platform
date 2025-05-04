package service

import (
    "context"
    "encoding/json"
    "vet-app-clinic-management/internal/models"
    "vet-app-clinic-management/internal/repository/mySQL"
    "vet-app-clinic-management/internal/repository/redis"
    "strconv"
    "fmt"
    "time"
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
            orderQty := item.Threshold - item.Quantity
            if orderQty <= 0 {
                orderQty = item.Threshold
            }

            // Симуляция заказа
            orderID := fmt.Sprintf("order_%s_%d", item.MedicineName, time.Now().UnixNano())
            orderData := map[string]interface{}{
                "medicine_id": item.ID,
                "name":        item.MedicineName,
                "quantity":    orderQty,
                "status":      "pending",
                "ordered_at":  time.Now().Format(time.RFC3339),
            }

            if s.cache != nil {
                orderJSON, err := json.Marshal(orderData)
                if err != nil {
                    return err
                }
                if err := s.cache.Set(ctx, orderID, string(orderJSON)); err != nil {
                    return err
                }
                if err := s.cache.Set(ctx, "order_"+item.MedicineName, orderID); err != nil {
                    return err
                }
            }
            
            updatedItem, err := s.UpdateQuantity(ctx, item.ID, item.Quantity+orderQty)
            if err != nil {
                return fmt.Errorf("failed to update inventory after auto-order: %v", err)
            }
            // Симуляция уведомления
            fmt.Printf("Auto-order initiated for %s (ID: %d), quantity: %d, orderID: %s, new quantity: %d\n",
                item.MedicineName, item.ID, orderQty, orderID, updatedItem.Quantity)

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
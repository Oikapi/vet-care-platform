package service

import (
    "context"
    "encoding/json"
    "vet-app-clinic-management/internal/models"
    "vet-app-clinic-management/internal/repository/mySQL"
    "vet-app-clinic-management/internal/repository/redis"
    "fmt"
    "time"
    "log"
)

type InventoryService struct {
    repo  *mySQL.InventoryRepo
    cache *redis.InventoryCache
}

func NewInventoryService(repo *mySQL.InventoryRepo, cache *redis.InventoryCache) *InventoryService {
    return &InventoryService{repo: repo, cache: cache}
}

func (s *InventoryService) GetAll(ctx context.Context, clinicID int) ([]*models.Inventory, error) {
    // Проверяем кэш
    cacheKey := fmt.Sprintf("inventory_clinic_%d", clinicID)
    if s.cache != nil {
        data, err := s.cache.Get(ctx, cacheKey)
        if err == nil && data != nil {
            dataStr, ok := data.(string)
            if !ok {
                return nil, fmt.Errorf("cached data is not a string")
            }
            if dataStr != "" {
                var inventory []*models.Inventory
                if err := json.Unmarshal([]byte(dataStr), &inventory); err == nil {
                    log.Printf("Fetched %d inventory items from cache for clinic ID %d", len(inventory), clinicID)
                    return inventory, nil
                }
            }
        }
    }

    // Если кэша нет, запрашиваем из MySQL
    inventory, err := s.repo.GetAll(clinicID)
    if err != nil {
        log.Println("Failed to fetch all inventory:", err)
        return nil, err
    }

    // Сохраняем в кэш (сериализуем в JSON)
    if s.cache != nil {
        inventoryBytes, err := json.Marshal(inventory)
        if err == nil {
            if err := s.cache.Set(ctx, cacheKey, string(inventoryBytes)); err != nil {
                log.Println("Failed to set cache for inventory:", err)
            }
        }
    }

    log.Printf("Fetched %d inventory items for clinic ID %d", len(inventory), clinicID)
    return inventory, nil
}

func (s *InventoryService) GetByID(ctx context.Context, clinicID, id int) (*models.Inventory, error) {
    // Проверяем кэш
    cacheKey := fmt.Sprintf("inventory_%d_clinic_%d", id, clinicID)
    if s.cache != nil {
        data, err := s.cache.Get(ctx, cacheKey)
        if err == nil && data != nil {
            dataStr, ok := data.(string)
            if !ok {
                return nil, fmt.Errorf("cached data is not a string")
            }
            if dataStr != "" {
                var inventory models.Inventory
                if err := json.Unmarshal([]byte(dataStr), &inventory); err == nil {
                    log.Printf("Fetched inventory item ID %d from cache for clinic ID %d", id, clinicID)
                    return &inventory, nil
                }
            }
        }
    }

    // Если кэша нет, запрашиваем из базы
    inventory, err := s.repo.GetByID(clinicID, id)
    if err != nil {
        log.Println("Failed to fetch inventory by ID:", err)
        return nil, err
    }

    // Обновляем кэш
    if s.cache != nil {
        inventoryBytes, err := json.Marshal(inventory)
        if err == nil {
            if err := s.cache.Set(ctx, cacheKey, string(inventoryBytes)); err != nil {
                log.Println("Failed to set cache for inventory item:", err)
            }
        }
    }

    log.Printf("Fetched inventory item ID %d for clinic ID %d", id, clinicID)
    return inventory, nil
}

func (s *InventoryService) AutoOrder(ctx context.Context, clinicID int) error {
    inventory, err := s.GetAll(ctx, clinicID)
    if err != nil {
        log.Println("Failed to fetch inventory for auto-order:", err)
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
                    log.Println("Failed to marshal order data:", err)
                    return err
                }
                if err := s.cache.Set(ctx, orderID, string(orderJSON)); err != nil {
                    log.Println("Failed to set cache for order:", err)
                    return err
                }
                if err := s.cache.Set(ctx, "order_"+item.MedicineName, orderID); err != nil {
                    log.Println("Failed to set cache for order reference:", err)
                    return err
                }
            }

            updatedItem, err := s.UpdateQuantity(ctx, clinicID, item.ID, item.Quantity+orderQty)
            if err != nil {
                log.Println("Failed to update inventory after auto-order:", err)
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
            cacheKey := fmt.Sprintf("inventory_clinic_%d", clinicID)
            if err := s.cache.Del(ctx, cacheKey); err != nil {
                log.Println("Failed to invalidate cache after auto-order:", err)
                return err
            }
        }
    }
    return nil
}

func (s *InventoryService) UpdateQuantity(ctx context.Context, clinicID, medicineID, quantity int) (*models.Inventory, error) {
    inventory, err := s.repo.GetByID(clinicID, medicineID)
    if err != nil {
        log.Println("Failed to fetch inventory for update quantity:", err)
        return nil, err
    }
    inventory.Quantity = quantity
    if err := s.repo.Update(inventory); err != nil {
        log.Println("Failed to update inventory quantity:", err)
        return nil, err
    }

    if s.cache != nil {
        cacheKey := fmt.Sprintf("inventory_%d_clinic_%d", medicineID, clinicID)
        inventoryBytes, err := json.Marshal(inventory)
        if err == nil {
            if err := s.cache.Set(ctx, cacheKey, string(inventoryBytes)); err != nil {
                log.Println("Failed to set cache for updated inventory item:", err)
            }
            // Инвалидация общего кэша
            cacheKeyAll := fmt.Sprintf("inventory_clinic_%d", clinicID)
            if err := s.cache.Del(ctx, cacheKeyAll); err != nil {
                log.Println("Failed to invalidate general cache after update:", err)
            }
        }
    }

    log.Printf("Updated inventory quantity for item ID %d in clinic ID %d", medicineID, clinicID)
    return inventory, nil
}

func (s *InventoryService) Create(ctx context.Context, clinicID int, medicineName string, quantity, threshold int) (*models.Inventory, error) {
    log.Printf("Creating inventory: medicine_name=%s for clinic ID %d", medicineName, clinicID)
    inventory := &models.Inventory{
        MedicineName: medicineName,
        Quantity:     quantity,
        Threshold:    threshold,
        ClinicID:     clinicID,
    }
    if err := s.repo.Create(inventory); err != nil {
        log.Println("Failed to create inventory:", err)
        return nil, err
    }

    // Инвалидация общего кэша
    if s.cache != nil {
        cacheKey := fmt.Sprintf("inventory_clinic_%d", clinicID)
        if err := s.cache.Del(ctx, cacheKey); err != nil {
            log.Println("Failed to invalidate cache after create:", err)
        }
    }

    log.Printf("Inventory created successfully: ID=%d", inventory.ID)
    return inventory, nil
}

func (s *InventoryService) Update(ctx context.Context, clinicID, id int, medicineName string, quantity, threshold int) (*models.Inventory, error) {
    log.Printf("Updating inventory with ID: %d for clinic ID %d", id, clinicID)
    inventory, err := s.repo.GetByID(clinicID, id)
    if err != nil {
        log.Println("Failed to fetch inventory for update:", err)
        return nil, err
    }
    inventory.MedicineName = medicineName
    inventory.Quantity = quantity
    inventory.Threshold = threshold
    if err := s.repo.Update(inventory); err != nil {
        log.Println("Failed to update inventory:", err)
        return nil, err
    }

    if s.cache != nil {
        cacheKey := fmt.Sprintf("inventory_%d_clinic_%d", id, clinicID)
        inventoryBytes, err := json.Marshal(inventory)
        if err == nil {
            if err := s.cache.Set(ctx, cacheKey, string(inventoryBytes)); err != nil {
                log.Println("Failed to set cache for updated inventory item:", err)
            }
            // Инвалидация общего кэша
            cacheKeyAll := fmt.Sprintf("inventory_clinic_%d", clinicID)
            if err := s.cache.Del(ctx, cacheKeyAll); err != nil {
                log.Println("Failed to invalidate general cache after update:", err)
            }
        }
    }

    log.Printf("Inventory updated successfully: ID=%d", inventory.ID)
    return inventory, nil
}

func (s *InventoryService) Delete(ctx context.Context, clinicID, id int) error {
    log.Printf("Deleting inventory with ID: %d for clinic ID %d", id, clinicID)
    if err := s.repo.Delete(clinicID, id); err != nil {
        log.Println("Failed to delete inventory:", err)
        return err
    }

    // Инвалидация кэша
    if s.cache != nil {
        cacheKey := fmt.Sprintf("inventory_%d_clinic_%d", id, clinicID)
        if err := s.cache.Del(ctx, cacheKey); err != nil {
            log.Println("Failed to delete cache for inventory item:", err)
        }
        cacheKeyAll := fmt.Sprintf("inventory_clinic_%d", clinicID)
        if err := s.cache.Del(ctx, cacheKeyAll); err != nil {
            log.Println("Failed to delete general cache after delete:", err)
        }
    }

    log.Printf("Inventory deleted successfully: ID=%d", id)
    return nil
}
package service

import "context"

func (s *InventoryService) CheckAndAutoOrder(ctx context.Context) error {
    inventory, err := s.GetAll(ctx)
    if err != nil {
        return err
    }

    for _, item := range inventory {
        if item.Quantity < item.Threshold {
            // Здесь логика автозаказа (например, вызов API поставщика)
            // Для примера просто обновим количество
            s.UpdateQuantity(item.ID, item.Quantity+100)
        }
    }
    return nil
}
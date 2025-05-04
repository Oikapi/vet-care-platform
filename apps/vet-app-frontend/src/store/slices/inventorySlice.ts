import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { apiInstance } from '../../api/auth';

const PREFIX = 'management';

interface InventoryItem {
  id: number;
  medicine_name: string;
  quantity: number;
  threshold: number;
}

interface InventoryState {
  items: InventoryItem[];
  loading: boolean;
  error: string | null;
}

const initialState: InventoryState = {
  items: [],
  loading: false,
  error: null,
};

export const fetchInventory = createAsyncThunk(
  'inventory/fetchAll',
  async () => {
    const response = await apiInstance.get(`${PREFIX}management/inventory`);
    return response.data;
  }
);

export const createInventoryItem = createAsyncThunk(
  'inventory/create',
  async (itemData: Omit<InventoryItem, 'id'>) => {
    const response = await apiInstance.post(
      `${PREFIX}management/inventory`,
      itemData
    );
    return response.data;
  }
);

export const updateInventoryItem = createAsyncThunk(
  'inventory/update',
  async ({
    id,
    itemData,
  }: {
    id: number;
    itemData: Partial<InventoryItem>;
  }) => {
    const response = await apiInstance.put(
      `${PREFIX}management/inventory/${id}`,
      itemData
    );
    return response.data;
  }
);

export const updateQuantity = createAsyncThunk(
  'inventory/updateQuantity',
  async ({ id, quantity }: { id: number; quantity: number }) => {
    await apiInstance.put(`${PREFIX}management/inventory/${id}/quantity`, {
      quantity,
    });
    return { id, quantity };
  }
);

export const deleteInventoryItem = createAsyncThunk(
  'inventory/delete',
  async (id: number) => {
    await apiInstance.delete(`${PREFIX}management/inventory/${id}`);
    return id;
  }
);

export const autoOrder = createAsyncThunk('inventory/autoOrder', async () => {
  await apiInstance.post(`${PREFIX}management/inventory/autoorder`);
});

const inventorySlice = createSlice({
  name: 'inventory',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchInventory.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchInventory.fulfilled, (state, action) => {
        state.loading = false;
        state.items = action.payload;
      })
      .addCase(fetchInventory.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch inventory';
      })
      .addCase(createInventoryItem.fulfilled, (state, action) => {
        state.items.push(action.payload);
      })
      .addCase(updateInventoryItem.fulfilled, (state, action) => {
        const index = state.items.findIndex((i) => i.id === action.payload.id);
        if (index !== -1) {
          state.items[index] = action.payload;
        }
      })
      .addCase(updateQuantity.fulfilled, (state, action) => {
        const item = state.items.find((i) => i.id === action.payload.id);
        if (item) {
          item.quantity = action.payload.quantity;
        }
      })
      .addCase(deleteInventoryItem.fulfilled, (state, action) => {
        state.items = state.items.filter((i) => i.id !== action.payload);
      });
  },
});

export default inventorySlice.reducer;

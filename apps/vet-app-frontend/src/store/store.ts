import { configureStore } from '@reduxjs/toolkit';
import schedulesReducer from './slices/shedulesSlice';
import inventoryReducer from './slices/inventorySlice';

export const store = configureStore({
  reducer: {
    schedules: schedulesReducer,
    inventory: inventoryReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

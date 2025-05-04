import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { apiInstance } from '../../api/auth';
import dayjs from 'dayjs';

const PREFIX = 'management/';

interface Schedule {
  id: number;
  doctor_id: number;
  doctor_name: string;
  start_time: string;
  end_time: string;
}

interface SchedulesState {
  schedules: Schedule[];
  loading: boolean;
  error: string | null;
}

const initialState: SchedulesState = {
  schedules: [],
  loading: false,
  error: null,
};

export const fetchSchedules = createAsyncThunk(
  'schedules/fetchAll',
  async () => {
    const response = await apiInstance.get(
      `${PREFIX}management/schedules/doctor`
    );
    return response.data;
  }
);

export const createSchedule = createAsyncThunk(
  'schedules/create',
  async (scheduleData: {
    doctor_id: number;
    start_time: string;
    end_time: string;
  }) => {
    const response = await apiInstance.post(
      `${PREFIX}/management/schedules`,
      scheduleData
    );
    return response.data;
  }
);

export const updateSchedule = createAsyncThunk(
  'schedules/update',
  async ({ id, scheduleData }: { id: number; scheduleData: any }) => {
    const response = await apiInstance.put(
      `${PREFIX}management/schedules/${id}`,
      scheduleData
    );
    return response.data;
  }
);

export const deleteSchedule = createAsyncThunk(
  'schedules/delete',
  async (id: number) => {
    await apiInstance.delete(`${PREFIX}management/schedules/${id}`);
    return id;
  }
);

const schedulesSlice = createSlice({
  name: 'schedules',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchSchedules.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchSchedules.fulfilled, (state, action) => {
        state.loading = false;
        state.schedules = action.payload;
      })
      .addCase(fetchSchedules.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch schedules';
      })
      .addCase(createSchedule.fulfilled, (state, action) => {
        state.schedules.push(action.payload);
      })
      .addCase(updateSchedule.fulfilled, (state, action) => {
        const index = state.schedules.findIndex(
          (s) => s.id === action.payload.id
        );
        if (index !== -1) {
          state.schedules[index] = action.payload;
        }
      })
      .addCase(deleteSchedule.fulfilled, (state, action) => {
        state.schedules = state.schedules.filter(
          (s) => s.id !== action.payload
        );
      });
  },
});

export default schedulesSlice.reducer;

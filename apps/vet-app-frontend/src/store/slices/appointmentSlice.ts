import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { apiInstance } from '../../api/auth';
import dayjs from 'dayjs';
import { Appointment, Slot } from '../../types';

interface AppointmentsState {
  appointments: Appointment[];
  availableSlots: Slot[];
  loading: boolean;
  error: string | null;
}

const initialState: AppointmentsState = {
  appointments: [],
  availableSlots: [],
  loading: false,
  error: null,
};

const PREFIX = 'appointment/';

export const createAppointment = createAsyncThunk(
  'appointments/create',
  async (
    appointmentData: { slot_id: number; patient_id: number },
    { rejectWithValue }
  ) => {
    try {
      const response = await apiInstance.post(
        `${PREFIX}appointments/`,
        appointmentData
      );
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response.data);
    }
  }
);

export const fetchAvailableSlots = createAsyncThunk(
  'appointments/fetchSlots',
  async (clinic_id: number, { rejectWithValue }) => {
    try {
      const response = await apiInstance.get(
        `${PREFIX}appointments/slots/${clinic_id}`
      );
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response.data);
    }
  }
);

export const fetchAppointment = createAsyncThunk(
  'appointments/fetchOne',
  async (id: number, { rejectWithValue }) => {
    try {
      const response = await apiInstance.get(`${PREFIX}appointments/${id}`);
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response.data);
    }
  }
);

export const createSlot = createAsyncThunk(
  'appointments/createSlot',
  async (
    slotData: {
      doctor_id: number;
      slot_time: string;
    },
    { rejectWithValue }
  ) => {
    try {
      const response = await apiInstance.post(`${PREFIX}appointments/slots`, {
        ...slotData,
        start_time: dayjs(slotData.slot_time).toISOString(),
      });
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response.data);
    }
  }
);

const appointmentsSlice = createSlice({
  name: 'appointments',
  initialState,
  reducers: {
    clearAvailableSlots(state) {
      state.availableSlots = [];
    },
    resetAppointmentsError(state) {
      console.log(state);
    },
  },
  extraReducers: (builder) => {
    builder
      // Create Appointment
      .addCase(createAppointment.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(
        createAppointment.fulfilled,
        (state, action: PayloadAction<Appointment>) => {
          state.loading = false;
          state.appointments.push(action.payload);
          // Mark slot as unavailable
          state.availableSlots = state.availableSlots.map((slot) =>
            slot.id === action.payload.slot_id
              ? { ...slot, is_available: false }
              : slot
          );
        }
      )
      .addCase(createAppointment.rejected, (state, action) => {
        state.loading = false;
        state.error =
          (action.payload as string) || 'Failed to create appointment';
      })

      // Fetch Available Slots
      .addCase(fetchAvailableSlots.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(
        fetchAvailableSlots.fulfilled,
        (state, action: PayloadAction<Slot[]>) => {
          state.loading = false;
          state.availableSlots = action.payload;
        }
      )
      .addCase(fetchAvailableSlots.rejected, (state, action) => {
        state.loading = false;
        state.error =
          (action.payload.error as string) || 'Failed to fetch available slots';
      })

      // Fetch Single Appointment
      .addCase(fetchAppointment.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(
        fetchAppointment.fulfilled,
        (state, action: PayloadAction<Appointment>) => {
          state.loading = false;
          // Replace if exists, add if new
          const index = state.appointments.findIndex(
            (a) => a.id === action.payload.id
          );
          if (index !== -1) {
            state.appointments[index] = action.payload;
          } else {
            state.appointments.push(action.payload);
          }
        }
      )
      .addCase(fetchAppointment.rejected, (state, action) => {
        console.log(action.payload);
        state.loading = false;
        state.error =
          (action.payload.error as string) || 'Failed to fetch appointment';
      })

      // Create Slot
      .addCase(createSlot.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(createSlot.fulfilled, (state, action: PayloadAction<Slot>) => {
        state.loading = false;
        state.availableSlots.push(action.payload);
      })
      .addCase(createSlot.rejected, (state, action) => {
        state.loading = false;
        state.error = (action.payload as string) || 'Failed to create slot';
      });
  },
});

export const { clearAvailableSlots, resetAppointmentsError } =
  appointmentsSlice.actions;
export default appointmentsSlice.reducer;

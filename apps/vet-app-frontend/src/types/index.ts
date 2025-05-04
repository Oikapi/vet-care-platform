export interface Slot {
  id: number;
  clinic_id: number;
  doctor_id: number;
  start_time: string;
  end_time: string;
  is_available: boolean;
}

export interface Appointment {
  id: number;
  slot_id: number;
  patient_id: number;
  status: string;
  created_at: string;
  updated_at: string;
  slot: Slot;
}

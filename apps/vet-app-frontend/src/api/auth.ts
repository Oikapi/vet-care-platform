import axios from 'axios';

const apiInstance = axios.create({
  baseURL: 'http://localhost:3000',
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const authApi = {
  registerUser: async (userData) => {
    try {
      const response = await apiInstance.post('/auth/register/user', userData);
      return response.data;
    } catch (error) {
      throw error.response?.data || { message: 'Network Error' };
    }
  },

  registerClinic: async (clinicData) => {
    try {
      const response = await apiInstance.post(
        '/auth/register/clinic',
        clinicData
      );
      return response.data;
    } catch (error) {
      throw error.response?.data || { message: 'Network Error' };
    }
  },

  loginUser: async (credentials) => {
    try {
      const response = await apiInstance.post('/auth/login', credentials);
      return response.data;
    } catch (error) {
      throw error.response?.data || { message: 'Network Error' };
    }
  },

  loginClinic: async (credentials) => {
    try {
      const response = await apiInstance.post('/auth/login', credentials);
      return response.data;
    } catch (error) {
      throw error.response?.data || { message: 'Network Error' };
    }
  },
};

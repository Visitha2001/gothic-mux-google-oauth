import axios from 'axios';
import type { LoginRequest, RegisterRequest, AuthResponse, User } from '../types';

const API_BASE_URL = 'http://localhost:8082';

const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
});

export const authAPI = {
  login: async (credentials: LoginRequest): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>('/auth/login', credentials);
    return response.data;
  },

  register: async (userData: RegisterRequest): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>('/auth/register', userData);
    return response.data;
  },

  logout: async (): Promise<void> => {
    await api.post('/auth/log-out');
  },

  getCurrentUser: async (): Promise<User> => {
    const response = await api.get<User>('/auth/me');
    return response.data;
  },

  // OAuth endpoints
  oauthLogin: (provider: string) => {
    window.location.href = `${API_BASE_URL}/auth/${provider}`;
  },
};

export default api;
import { useEffect } from 'react';
import { useAuthStore } from '../store/authStore';
import { authAPI } from '../services/api';
import type { LoginRequest, RegisterRequest } from '../types';

export const useAuth = () => {
  const { user, isAuthenticated, isLoading, setUser, setLoading, logout: storeLogout } = useAuthStore();

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      const userData = await authAPI.getCurrentUser();
      setUser(userData);
    } catch (error) {
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  const login = async (credentials: LoginRequest) => {
    try {
      const response = await authAPI.login(credentials);
      setUser(response.user);
      return { success: true };
    } catch (error: any) {
      return { 
        success: false, 
        error: error.response?.data?.message || 'Login failed' 
      };
    }
  };

  const register = async (userData: RegisterRequest) => {
    try {
      const response = await authAPI.register(userData);
      setUser(response.user);
      return { success: true };
    } catch (error: any) {
      return { 
        success: false, 
        error: error.response?.data?.message || 'Registration failed' 
      };
    }
  };

  const logout = async () => {
    try {
      await authAPI.logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      storeLogout();
    }
  };

  const oauthLogin = (provider: string) => {
    authAPI.oauthLogin(provider);
  };

  return {
    user,
    isAuthenticated,
    isLoading,
    login,
    register,
    logout,
    oauthLogin,
    checkAuth,
  };
};
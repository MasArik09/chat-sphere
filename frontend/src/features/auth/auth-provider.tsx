import { useEffect } from 'react';
import type { ReactNode, FC } from 'react';
import { AuthContext } from './auth-context';
import { useAuthStore } from '../../store/auth-store';
import { authService } from './auth-service';
import type { LoginPayload, RegisterPayload } from './types';
import { webSocketService } from '../../services/websocket';

export const AuthProvider: FC<{ children: ReactNode }> = ({ children }) => {
  const {
    user,
    token,
    isAuthenticated,
    isLoading,
    error,
    setUser,
    setIsLoading,
    setError,
    login: storeLogin,
    logout: storeLogout,
  } = useAuthStore();

  // Manage WebSocket lifecycle reactively based on token and auth state
  useEffect(() => {
    if (token && isAuthenticated) {
      webSocketService.connect(token);
    } else {
      webSocketService.disconnect();
    }
    return () => {
      webSocketService.disconnect();
    };
  }, [token, isAuthenticated]);

  // Auto-login on mount/refresh if token exists
  useEffect(() => {
    const autoLogin = async () => {
      if (!token) return;
      
      setIsLoading(true);
      try {
        const response = await authService.getMe();
        if (response.success && response.user) {
          setUser(response.user);
        } else {
          storeLogout();
        }
      } catch (err) {
        storeLogout();
      } finally {
        setIsLoading(false);
      }
    };

    autoLogin();
  }, [token, setUser, setIsLoading, storeLogout]);

  const login = async (payload: LoginPayload) => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await authService.login(payload);
      if (response.success && response.token && response.user) {
        storeLogin(response.token, response.user);
      } else {
        setError(response.message || 'Login failed');
      }
    } catch (err: any) {
      const msg = err.response?.data?.message || 'Login failed. Please check your credentials.';
      setError(msg);
      throw new Error(msg);
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (payload: RegisterPayload) => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await authService.register(payload);
      if (!response.success) {
        setError(response.message || 'Registration failed');
      }
    } catch (err: any) {
      const msg = err.response?.data?.message || 'Registration failed. Please check your inputs.';
      setError(msg);
      throw new Error(msg);
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    storeLogout();
  };

  const clearError = () => {
    setError(null);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        isAuthenticated,
        isLoading,
        error,
        login,
        register,
        logout,
        clearError,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

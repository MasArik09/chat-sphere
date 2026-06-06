import { create } from 'zustand';
import type { User } from '../api/types';

interface AuthStoreState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  setUser: (user: User | null) => void;
  setToken: (token: string | null) => void;
  setIsLoading: (isLoading: boolean) => void;
  setError: (error: string | null) => void;
  login: (token: string, user: User) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthStoreState>((set) => ({
  user: null,
  token: localStorage.getItem('chatsphere_token'),
  isAuthenticated: !!localStorage.getItem('chatsphere_token'),
  isLoading: false,
  error: null,
  setUser: (user) => set({ user, isAuthenticated: !!user }),
  setToken: (token) => {
    if (token) {
      localStorage.setItem('chatsphere_token', token);
    } else {
      localStorage.removeItem('chatsphere_token');
    }
    set({ token, isAuthenticated: !!token });
  },
  setIsLoading: (isLoading) => set({ isLoading }),
  setError: (error) => set({ error }),
  login: (token, user) => {
    localStorage.setItem('chatsphere_token', token);
    set({ token, user, isAuthenticated: true, error: null });
  },
  logout: () => {
    localStorage.removeItem('chatsphere_token');
    set({ token: null, user: null, isAuthenticated: false, error: null });
  },
}));

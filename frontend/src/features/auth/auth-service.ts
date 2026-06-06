import { apiClient } from '../../api/axios';
import { ENDPOINTS } from '../../api/endpoints';
import type { AuthResponse, User } from '../../api/types';
import type { LoginPayload, RegisterPayload } from './types';

export const authService = {
  async register(payload: RegisterPayload): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>(ENDPOINTS.AUTH.REGISTER, payload);
    return response.data;
  },

  async login(payload: LoginPayload): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>(ENDPOINTS.AUTH.LOGIN, payload);
    return response.data;
  },

  async getMe(): Promise<{ success: boolean; user: User }> {
    const response = await apiClient.get<{ success: boolean; user: User }>(ENDPOINTS.AUTH.ME);
    return response.data;
  },
};

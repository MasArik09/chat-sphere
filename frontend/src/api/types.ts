export interface User {
  id: number;
  name: string;
  email: string;
  is_online: boolean;
  last_seen_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  success: boolean;
  token?: string;
  user?: User;
  message?: string;
}

export interface ApiErrorResponse {
  success: boolean;
  message: string;
}

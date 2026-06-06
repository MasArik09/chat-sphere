export const ENDPOINTS = {
  AUTH: {
    REGISTER: '/auth/register',
    LOGIN: '/auth/login',
    ME: '/auth/me',
  },
  CONVERSATIONS: {
    BASE: '/conversations',
    DETAIL: (id: number | string) => `/conversations/${id}`,
    PARTICIPANTS: (id: number | string) => `/conversations/${id}/participants`,
    PARTICIPANT: (id: number | string, userId: number | string) => `/conversations/${id}/participants/${userId}`,
    MESSAGES: (id: number | string) => `/conversations/${id}/messages`,
  },
};

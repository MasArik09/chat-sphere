import { create } from 'zustand';
import { apiClient } from '../api/axios';
import { ENDPOINTS } from '../api/endpoints';

export interface Message {
  id: number;
  conversation_id: number;
  sender_id: number;
  content: string;
  sent_at: string;
  created_at: string;
  updated_at: string;
}

interface MessageStoreState {
  messages: Record<number, Message[]>;
  isLoading: boolean;
  isSending: boolean;
  fetchMessages: (conversationId: number, page?: number, limit?: number) => Promise<void>;
  sendMessage: (conversationId: number, content: string) => Promise<Message>;
  addWebSocketMessage: (message: Message) => void;
}

export const useMessageStore = create<MessageStoreState>((set) => ({
  messages: {},
  isLoading: false,
  isSending: false,
  fetchMessages: async (conversationId, page = 1, limit = 100) => {
    set({ isLoading: true });
    try {
      const response = await apiClient.get<{ success: boolean; messages: Message[] }>(
        `${ENDPOINTS.CONVERSATIONS.MESSAGES(conversationId)}?page=${page}&limit=${limit}`
      );
      if (response.data.success) {
        set((state) => ({
          messages: {
            ...state.messages,
            [conversationId]: response.data.messages,
          },
        }));
      }
    } catch (err) {
      console.error('Failed to fetch messages:', err);
    } finally {
      set({ isLoading: false });
    }
  },
  sendMessage: async (conversationId, content) => {
    set({ isSending: true });
    try {
      const response = await apiClient.post<{ success: boolean; message: Message }>(
        ENDPOINTS.CONVERSATIONS.MESSAGES(conversationId),
        { content }
      );
      if (response.data.success) {
        const newMsg = response.data.message;
        set((state) => {
          const convMsgs = state.messages[conversationId] || [];
          if (convMsgs.some((m) => m.id === newMsg.id)) {
            return {};
          }
          return {
            messages: {
              ...state.messages,
              [conversationId]: [...convMsgs, newMsg],
            },
          };
        });
        return newMsg;
      }
      throw new Error('Failed to send message');
    } finally {
      set({ isSending: false });
    }
  },
  addWebSocketMessage: (message) => {
    const conversationId = message.conversation_id;
    set((state) => {
      const convMsgs = state.messages[conversationId] || [];
      if (convMsgs.some((m) => m.id === message.id)) {
        return {};
      }
      return {
        messages: {
          ...state.messages,
          [conversationId]: [...convMsgs, message],
        },
      };
    });
  },
}));

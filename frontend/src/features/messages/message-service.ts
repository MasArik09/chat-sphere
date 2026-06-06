import { apiClient } from '../../api/axios';
import { ENDPOINTS } from '../../api/endpoints';
import type { Message } from '../../store/message-store';

export const messageService = {
  async getMessages(conversationId: number, page = 1, limit = 100): Promise<Message[]> {
    const response = await apiClient.get<{ success: boolean; messages: Message[] }>(
      `${ENDPOINTS.CONVERSATIONS.MESSAGES(conversationId)}?page=${page}&limit=${limit}`
    );
    return response.data.messages;
  },

  async sendMessage(conversationId: number, content: string): Promise<Message> {
    const response = await apiClient.post<{ success: boolean; message: Message }>(
      ENDPOINTS.CONVERSATIONS.MESSAGES(conversationId),
      { content }
    );
    return response.data.message;
  },
};

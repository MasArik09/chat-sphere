import { apiClient } from '../../api/axios';
import { ENDPOINTS } from '../../api/endpoints';
import type { Conversation } from '../../store/conversation-store';

export const conversationService = {
  async getConversations(): Promise<Conversation[]> {
    const response = await apiClient.get<{ success: boolean; conversations: Conversation[] }>(
      ENDPOINTS.CONVERSATIONS.BASE
    );
    return response.data.conversations;
  },

  async getConversationDetail(id: number): Promise<Conversation> {
    const response = await apiClient.get<{ success: boolean; conversation: Conversation }>(
      ENDPOINTS.CONVERSATIONS.DETAIL(id)
    );
    return response.data.conversation;
  },

  async createConversation(participantIds: number[]): Promise<Conversation> {
    const response = await apiClient.post<{ success: boolean; conversation: Conversation }>(
      ENDPOINTS.CONVERSATIONS.BASE,
      { participant_ids: participantIds }
    );
    return response.data.conversation;
  },
};

import { create } from 'zustand';
import { apiClient } from '../api/axios';
import { ENDPOINTS } from '../api/endpoints';

export interface ParticipantUser {
  id: number;
  name: string;
  email: string;
  is_online: boolean;
  last_seen_at?: string;
}

export interface ConversationParticipant {
  id: number;
  conversation_id: number;
  user_id: number;
  created_at: string;
  user?: ParticipantUser;
  last_read_message_id?: number | null;
}

export interface MessagePreview {
  id?: number | null;
  content?: string | null;
  sender_id?: number | null;
  sent_at?: string | null;
}

export interface Conversation {
  id: number;
  participant_count: number;
  created_at: string;
  updated_at: string;
  participants?: ConversationParticipant[];
  last_message?: MessagePreview | null;
  unread_count?: number;
}

interface ConversationStoreState {
  conversations: Conversation[];
  selectedConversation: Conversation | null;
  isLoading: boolean;
  setConversations: (conversations: Conversation[]) => void;
  setSelectedConversation: (conversation: Conversation | null) => void;
  setIsLoading: (isLoading: boolean) => void;
  fetchConversations: (search?: string) => Promise<void>;
  createConversation: (participantIds: number[]) => Promise<Conversation>;
  markAsRead: (conversationId: number, lastReadMessageId: number) => Promise<void>;
}

export const useConversationStore = create<ConversationStoreState>((set, get) => ({
  conversations: [],
  selectedConversation: null,
  isLoading: false,
  setConversations: (conversations) => set({ conversations }),
  setSelectedConversation: (selectedConversation) => set({ selectedConversation }),
  setIsLoading: (isLoading) => set({ isLoading }),
  fetchConversations: async (search?: string) => {
    set({ isLoading: true });
    try {
      const url = search
        ? `${ENDPOINTS.CONVERSATIONS.BASE}?search=${encodeURIComponent(search)}`
        : ENDPOINTS.CONVERSATIONS.BASE;
      const response = await apiClient.get<{ success: boolean; conversations: Conversation[] }>(url);
      if (response.data.success) {
        const convList = response.data.conversations;
        
        // Dynamic import to prevent potential circular dependency issues
        const { setOnline } = await import('./presence-store').then((m) => m.usePresenceStore.getState());
        convList.forEach((c) => {
          c.participants?.forEach((p) => {
            if (p.user) {
              setOnline(p.user_id, p.user.is_online);
            }
          });
        });
        
        set({ conversations: convList });
      }
    } catch (err) {
      console.error('Failed to fetch conversations:', err);
    } finally {
      set({ isLoading: false });
    }
  },
  createConversation: async (participantIds: number[]) => {
    const response = await apiClient.post<{ success: boolean; conversation: Conversation }>(
      ENDPOINTS.CONVERSATIONS.BASE,
      { participant_ids: participantIds }
    );
    if (response.data.success) {
      const newConv = response.data.conversation;
      
      const { setOnline } = await import('./presence-store').then((m) => m.usePresenceStore.getState());
      newConv.participants?.forEach((p) => {
        if (p.user) {
          setOnline(p.user_id, p.user.is_online);
        }
      });
      
      const currentConvs = get().conversations;
      if (!currentConvs.some((c) => c.id === newConv.id)) {
        set({ conversations: [newConv, ...currentConvs] });
      }
      return newConv;
    }
    throw new Error('Failed to create conversation');
  },
  markAsRead: async (conversationId: number, lastReadMessageId: number) => {
    try {
      const response = await apiClient.post<{ success: boolean }>(
        `${ENDPOINTS.CONVERSATIONS.BASE}/${conversationId}/read`,
        { last_read_message_id: lastReadMessageId }
      );
      if (response.data.success) {
        const { useAuthStore } = await import('./auth-store');
        const currentUserId = useAuthStore.getState().user?.id;
        
        set((state) => {
          const updated = state.conversations.map((c) => {
            if (c.id === conversationId) {
              return {
                ...c,
                unread_count: 0,
                participants: c.participants?.map((p) => {
                  if (p.user_id === currentUserId) {
                    return { ...p, last_read_message_id: lastReadMessageId };
                  }
                  return p;
                }),
              };
            }
            return c;
          });

          const selected = state.selectedConversation;
          const updatedSelected = selected && selected.id === conversationId
            ? {
                ...selected,
                unread_count: 0,
                participants: selected.participants?.map((p) => {
                  if (p.user_id === currentUserId) {
                    return { ...p, last_read_message_id: lastReadMessageId };
                  }
                  return p;
                }),
              }
            : selected;

          return {
            conversations: updated,
            selectedConversation: updatedSelected,
          };
        });
      }
    } catch (err) {
      console.error('Failed to mark conversation as read:', err);
    }
  },
}));

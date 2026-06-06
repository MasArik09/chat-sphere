import { create } from 'zustand';

interface PresenceStoreState {
  onlineUsers: Set<number>;
  typingStates: Record<string, boolean>; // key format: `${conversationId}-${userId}`
  setOnline: (userId: number, online: boolean) => void;
  isOnline: (userId: number) => boolean;
  setTyping: (conversationId: number, userId: number, typing: boolean) => void;
  isTyping: (conversationId: number, userId: number) => boolean;
  clearUserTyping: (userId: number) => void;
}

export const usePresenceStore = create<PresenceStoreState>((set, get) => ({
  onlineUsers: new Set<number>(),
  typingStates: {},
  setOnline: (userId, online) =>
    set((state) => {
      const next = new Set(state.onlineUsers);
      if (online) {
        next.add(userId);
      } else {
        next.delete(userId);
      }
      return { onlineUsers: next };
    }),
  isOnline: (userId) => get().onlineUsers.has(userId),
  setTyping: (conversationId, userId, typing) =>
    set((state) => ({
      typingStates: {
        ...state.typingStates,
        [`${conversationId}-${userId}`]: typing,
      },
    })),
  isTyping: (conversationId, userId) => !!get().typingStates[`${conversationId}-${userId}`],
  clearUserTyping: (userId) =>
    set((state) => {
      const nextTyping = { ...state.typingStates };
      Object.keys(nextTyping).forEach((key) => {
        const parts = key.split('-');
        if (parts.length === 2 && parseInt(parts[1], 10) === userId) {
          delete nextTyping[key];
        }
      });
      return { typingStates: nextTyping };
    }),
}));

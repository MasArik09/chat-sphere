import { useEffect } from 'react';
import type { FC } from 'react';
import { ConversationSidebar } from '../features/conversations/conversation-sidebar';
import { ChatWindow } from '../features/messages/chat-window';
import { useConversationStore } from '../store/conversation-store';

export const ChatPage: FC = () => {
  const { selectedConversation, fetchConversations } = useConversationStore();

  useEffect(() => {
    fetchConversations();
  }, [fetchConversations]);

  return (
    <div className="flex-1 flex h-full overflow-hidden bg-gray-50 relative">
      {/* Sidebar - Visible on desktop always, visible on mobile only when no chat is active */}
      <div
        className={`h-full shrink-0 w-full md:w-80 border-r border-gray-200
          ${selectedConversation ? 'hidden md:block' : 'block'}`}
      >
        <ConversationSidebar />
      </div>

      {/* Chat Window - Visible on desktop always, visible on mobile only when a chat is active */}
      <div
        className={`h-full flex-1
          ${selectedConversation ? 'block' : 'hidden md:block'}`}
      >
        <ChatWindow />
      </div>
    </div>
  );
};

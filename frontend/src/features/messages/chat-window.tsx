import { useEffect } from 'react';
import type { FC } from 'react';
import { ChatHeader } from './chat-header';
import { MessageList } from './message-list';
import { MessageInput } from './message-input';
import { useConversationStore } from '../../store/conversation-store';
import { useMessageStore } from '../../store/message-store';
import { usePresenceStore } from '../../store/presence-store';
import { useAuth } from '../auth/use-auth';
import { webSocketService } from '../../services/websocket';

export const ChatWindow: FC = () => {
  const { selectedConversation, setSelectedConversation, markAsRead } = useConversationStore();
  const { messages, isLoading, isSending, fetchMessages, sendMessage } = useMessageStore();
  const { user: currentUser } = useAuth();

  const otherParticipant = selectedConversation?.participants?.find((p) => p.user_id !== currentUser?.id);
  const otherUserId = otherParticipant?.user_id || 0;
  const isPartnerTyping = usePresenceStore((state) =>
    selectedConversation ? state.isTyping(selectedConversation.id, otherUserId) : false
  );

  useEffect(() => {
    if (selectedConversation) {
      fetchMessages(selectedConversation.id);
    }
  }, [selectedConversation, fetchMessages]);

  const convMessages = selectedConversation ? (messages[selectedConversation.id] || []) : [];

  // Mark conversation as read on load or when new messages arrive
  useEffect(() => {
    if (!selectedConversation || convMessages.length === 0 || !currentUser) return;
    const lastMessage = convMessages[convMessages.length - 1];
    const myParticipant = selectedConversation.participants?.find((p) => p.user_id === currentUser.id);
    const myLastReadId = myParticipant?.last_read_message_id || 0;
    
    if (lastMessage && lastMessage.id > myLastReadId) {
      markAsRead(selectedConversation.id, lastMessage.id);
    }
  }, [selectedConversation, convMessages, currentUser, markAsRead]);

  const handleSend = async (content: string) => {
    if (!selectedConversation) return;
    await sendMessage(selectedConversation.id, content);
  };

  const handleTyping = (typing: boolean) => {
    if (!selectedConversation) return;
    const event = typing ? 'typing.start' : 'typing.stop';
    webSocketService.sendEvent(event, { conversation_id: selectedConversation.id });
  };

  const handleBack = () => {
    setSelectedConversation(null);
  };

  if (!selectedConversation) {
    return (
      <div className="flex-1 flex flex-col items-center justify-center p-6 text-center text-gray-400 bg-gray-50 h-full select-none">
        <div className="h-16 w-16 rounded-full bg-white flex items-center justify-center shadow-sm text-gray-300 mb-4 border border-gray-100">
          <svg
            className="h-8 w-8"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            strokeWidth="1.5"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M8.625 9.75a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H8.25m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H12m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0h-.375m-13.5 3.01c0 1.6 1.123 2.994 2.707 3.227 1.087.16 2.185.283 3.293.369V21l4.184-4.183a1.14 1.14 0 01.778-.332 48.294 48.294 0 005.83-.498c1.585-.233 2.708-1.626 2.708-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0012 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018z"
            />
          </svg>
        </div>
        <h5 className="text-base font-bold text-gray-800">No conversation selected</h5>
        <p className="text-xs text-gray-500 mt-1.5 max-w-xs leading-relaxed">
          Choose a conversation from the sidebar or click the new chat button to start messaging.
        </p>
      </div>
    );
  }

  return (
    <section className="flex-1 flex flex-col h-full bg-gray-50 relative overflow-hidden z-10">
      <ChatHeader
        conversation={selectedConversation}
        onBack={handleBack}
      />
      
      <MessageList
        messages={convMessages}
        isLoading={isLoading}
        isPartnerTyping={isPartnerTyping}
      />
      
      <MessageInput
        onSend={handleSend}
        disabled={isSending}
        onTyping={handleTyping}
      />
    </section>
  );
};

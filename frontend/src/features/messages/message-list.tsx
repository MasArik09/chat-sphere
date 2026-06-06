import { useRef, useEffect } from 'react';
import type { FC } from 'react';
import { MessageBubble } from './message-bubble';
import type { Message } from '../../store/message-store';
import { useAuth } from '../auth/use-auth';

interface MessageListProps {
  messages: Message[];
  isLoading: boolean;
  isPartnerTyping?: boolean;
}

export const MessageList: FC<MessageListProps> = ({ messages, isLoading, isPartnerTyping }) => {
  const { user } = useAuth();
  const scrollRef = useRef<HTMLDivElement>(null);
  const prevMessagesCountRef = useRef(messages.length);

  useEffect(() => {
    if (!scrollRef.current) return;

    const { scrollHeight, scrollTop, clientHeight } = scrollRef.current;
    const isAtBottom = scrollHeight - scrollTop - clientHeight < 150;
    const isFirstLoad = prevMessagesCountRef.current === 0 && messages.length > 0;
    const isNewMessageSent = messages.length > prevMessagesCountRef.current && 
      messages[messages.length - 1]?.sender_id === user?.id;

    if (isAtBottom || isFirstLoad || isNewMessageSent || isPartnerTyping) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }

    prevMessagesCountRef.current = messages.length;
  }, [messages, user?.id, isPartnerTyping]);

  if (isLoading && messages.length === 0) {
    return (
      <div className="flex-1 overflow-y-auto p-6 flex flex-col gap-4">
        {[...Array(3)].map((_, i) => (
          <div
            key={i}
            className={`flex w-full mb-3.5 animate-pulse ${
              i % 2 === 0 ? 'justify-end' : 'justify-start'
            }`}
          >
            <div className="max-w-[50%] h-12 bg-gray-200 rounded-2xl w-48" />
          </div>
        ))}
      </div>
    );
  }

  if (messages.length === 0) {
    return (
      <div className="flex-1 flex flex-col items-center justify-center p-6 text-center text-gray-500 select-none">
        <div className="h-14 w-14 rounded-full bg-gray-100 flex items-center justify-center text-gray-400 mb-3">
          <svg
            className="h-6.5 w-6.5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            strokeWidth="2"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z"
            />
          </svg>
        </div>
        <p className="text-sm font-semibold text-gray-700">No messages yet</p>
        <p className="text-xs text-gray-400 mt-1">Send a message to start conversation.</p>
      </div>
    );
  }

  return (
    <div
      ref={scrollRef}
      className="flex-1 overflow-y-auto p-6 flex flex-col scroll-smooth"
    >
      {messages.map((m) => (
        <MessageBubble
          key={m.id}
          message={m}
          isOwn={m.sender_id === user?.id}
        />
      ))}
      {isPartnerTyping && (
        <div className="flex justify-start mb-3.5">
          <div className="max-w-[70%] bg-white border border-gray-150 text-gray-800 rounded-2xl rounded-bl-none px-4 py-3 shadow-sm flex items-center gap-1 min-h-[36px]">
            <span className="w-1.5 h-1.5 rounded-full bg-gray-400 animate-bounce" style={{ animationDelay: '0ms' }} />
            <span className="w-1.5 h-1.5 rounded-full bg-gray-400 animate-bounce" style={{ animationDelay: '150ms' }} />
            <span className="w-1.5 h-1.5 rounded-full bg-gray-400 animate-bounce" style={{ animationDelay: '300ms' }} />
          </div>
        </div>
      )}
    </div>
  );
};

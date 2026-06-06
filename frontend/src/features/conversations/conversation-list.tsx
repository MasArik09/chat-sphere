import type { FC } from 'react';
import { ConversationItem } from './conversation-item';
import type { Conversation } from '../../store/conversation-store';

interface ConversationListProps {
  conversations: Conversation[];
  selectedId: number | null;
  onSelect: (conversation: Conversation) => void;
  isLoading: boolean;
}

export const ConversationList: FC<ConversationListProps> = ({
  conversations,
  selectedId,
  onSelect,
  isLoading,
}) => {
  if (isLoading) {
    return (
      <div className="flex flex-col gap-3 p-1">
        {[...Array(4)].map((_, i) => (
          <div key={i} className="p-3.5 rounded-xl border border-transparent flex items-center gap-3.5 animate-pulse bg-white">
            <div className="h-11 w-11 rounded-full bg-gray-200 shrink-0" />
            <div className="flex-1 min-w-0 text-left">
              <div className="h-4 bg-gray-200 rounded w-1/3 mb-2.5" />
              <div className="h-3 bg-gray-200 rounded w-2/3" />
            </div>
          </div>
        ))}
      </div>
    );
  }

  if (conversations.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center p-8 text-center text-gray-500 h-64 select-none">
        <svg
          className="h-12 w-12 text-gray-300 mb-3"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          strokeWidth="1.5"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M8.625 12a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H8.25m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H12m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0h-.375M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
          />
        </svg>
        <p className="text-sm font-semibold text-gray-700">No conversations</p>
        <p className="text-xs text-gray-400 mt-1">Start a conversation to begin chatting.</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-2 p-1">
      {conversations.map((c) => (
        <ConversationItem
          key={c.id}
          conversation={c}
          isActive={selectedId === c.id}
          onClick={() => onSelect(c)}
        />
      ))}
    </div>
  );
};

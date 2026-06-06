import type { FC } from 'react';
import type { Conversation } from '../../store/conversation-store';
import { useAuth } from '../auth/use-auth';
import { usePresenceStore } from '../../store/presence-store';
import { getUserDisplayName } from '../../utils/user';

interface ChatHeaderProps {
  conversation: Conversation;
  onBack: () => void;
}

export const ChatHeader: FC<ChatHeaderProps> = ({ conversation, onBack }) => {
  const { user: currentUser } = useAuth();
  const otherParticipant = conversation.participants?.find((p) => p.user_id !== currentUser?.id);
  const otherUserId = otherParticipant?.user_id || 0;
  const displayName = otherParticipant?.user?.name || getUserDisplayName(otherUserId, currentUser?.id);

  const isOnline = usePresenceStore((state) => state.isOnline(otherUserId));
  const isTyping = usePresenceStore((state) => state.isTyping(conversation.id, otherUserId));

  const formatLastSeen = (isoString?: string) => {
    if (!isoString) return 'Offline';
    try {
      const date = new Date(isoString);
      const now = new Date();
      const isToday = date.toDateString() === now.toDateString();
      const timeStr = date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      if (isToday) {
        return `Last seen today at ${timeStr}`;
      }
      const dateStr = date.toLocaleDateString([], { month: 'short', day: 'numeric' });
      return `Last seen on ${dateStr} at ${timeStr}`;
    } catch {
      return 'Offline';
    }
  };

  return (
    <div className="h-16 border-b border-gray-200 bg-white px-4 md:px-6 flex items-center gap-3 shrink-0 z-10 shadow-sm">
      <button
        onClick={onBack}
        className="md:hidden p-1 rounded-lg text-gray-500 hover:bg-gray-100 hover:text-gray-700 transition-colors"
        title="Back to conversations"
      >
        <svg
          className="h-6 w-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          strokeWidth="2"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M15 19l-7-7 7-7"
          />
        </svg>
      </button>

      <div className="relative">
        <div className="h-10 w-10 rounded-full bg-primary-100 flex items-center justify-center font-bold text-primary-700 text-sm select-none">
          {displayName.charAt(0).toUpperCase()}
        </div>
        <span
          className={`absolute bottom-0 right-0 h-3 w-3 rounded-full border-2 border-white transition-all
            ${isOnline ? 'bg-green-500' : 'bg-gray-300'}`}
        />
      </div>

      <div className="flex-1 min-w-0 text-left">
        <h4 className="text-sm font-semibold text-gray-800 truncate">{displayName}</h4>
        {isTyping ? (
          <span className="text-xs font-semibold text-green-500 animate-pulse flex items-center gap-1">
            typing
            <span className="flex gap-0.5 items-center inline-flex h-2">
              <span className="w-1 h-1 rounded-full bg-green-500 animate-bounce" style={{ animationDelay: '0ms' }} />
              <span className="w-1 h-1 rounded-full bg-green-500 animate-bounce" style={{ animationDelay: '150ms' }} />
              <span className="w-1 h-1 rounded-full bg-green-500 animate-bounce" style={{ animationDelay: '300ms' }} />
            </span>
          </span>
        ) : (
          <span
            className={`text-xs font-medium transition-colors
              ${isOnline ? 'text-green-500' : 'text-gray-400'}`}
          >
            {isOnline ? 'Online' : formatLastSeen(otherParticipant?.user?.last_seen_at)}
          </span>
        )}
      </div>
    </div>
  );
};

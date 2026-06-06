import type { FC } from 'react';
import type { Conversation } from '../../store/conversation-store';
import { useAuth } from '../auth/use-auth';
import { usePresenceStore } from '../../store/presence-store';
import { getUserDisplayName } from '../../utils/user';

interface ConversationItemProps {
  conversation: Conversation;
  isActive: boolean;
  onClick: () => void;
}

export const ConversationItem: FC<ConversationItemProps> = ({
  conversation,
  isActive,
  onClick,
}) => {
  const { user: currentUser } = useAuth();
  const isOnline = usePresenceStore((state) => {
    const other = conversation.participants?.find((p) => p.user_id !== currentUser?.id);
    return other ? state.isOnline(other.user_id) : false;
  });

  const otherParticipant = conversation.participants?.find((p) => p.user_id !== currentUser?.id);
  const otherUserId = otherParticipant?.user_id || 0;
  const displayName = otherParticipant?.user?.name || getUserDisplayName(otherUserId, currentUser?.id);

  const formatTime = (isoString: string) => {
    try {
      const date = new Date(isoString);
      return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    } catch {
      return '';
    }
  };

  const lastMsgTime = conversation.last_message?.sent_at || conversation.updated_at;
  const previewText = conversation.last_message?.content || 'No messages yet';
  const isLastMessageFromMe = conversation.last_message?.sender_id === currentUser?.id;
  const isRead = isLastMessageFromMe && (() => {
    const otherPart = conversation.participants?.find((p) => p.user_id !== currentUser?.id);
    const lastMsgId = conversation.last_message?.id || 0;
    const otherLastReadId = otherPart?.last_read_message_id || 0;
    return lastMsgId > 0 && otherLastReadId >= lastMsgId;
  })();

  return (
    <div
      onClick={onClick}
      className={`p-3.5 rounded-xl flex items-center gap-3.5 cursor-pointer transition-all border
        ${
          isActive
            ? 'bg-primary-50 border-primary-100 shadow-sm'
            : 'bg-white hover:bg-gray-50 border-transparent'
        }`}
    >
      <div className="relative shrink-0">
        <div
          className={`h-11 w-11 rounded-full flex items-center justify-center font-bold text-sm select-none
            ${
              isActive
                ? 'bg-primary-100 text-primary-700'
                : 'bg-gray-100 text-gray-600'
            }`}
        >
          {displayName.charAt(0).toUpperCase()}
        </div>
        <span
          className={`absolute bottom-0 right-0 h-3 w-3 rounded-full border-2 border-white transition-all
            ${isOnline ? 'bg-green-500' : 'bg-gray-300'}`}
        />
      </div>

      <div className="flex-1 min-w-0">
        <div className="flex justify-between items-baseline mb-0.5">
          <span
            className={`text-sm font-semibold truncate
              ${isActive ? 'text-primary-950' : 'text-gray-900'}`}
          >
            {displayName}
          </span>
          <span className="text-[10px] text-gray-400 shrink-0 ml-1">
            {formatTime(lastMsgTime)}
          </span>
        </div>

        <div className="flex items-center justify-between">
          <p
            className={`text-xs truncate max-w-[150px]
              ${isActive ? 'text-primary-700' : 'text-gray-500'}`}
          >
            {isLastMessageFromMe && conversation.last_message && (
              <span className="inline-flex mr-1 align-middle text-[11px] font-bold">
                {isRead ? (
                  <span className="text-blue-500" title="Read">✓✓</span>
                ) : (
                  <span className="text-gray-400" title="Delivered">✓</span>
                )}
              </span>
            )}
            {previewText}
          </p>
          
          {conversation.unread_count && conversation.unread_count > 0 ? (
            <span className="text-[9px] font-semibold bg-green-500 text-white min-w-5 h-5 px-1.5 rounded-full flex items-center justify-center shadow-sm shrink-0 animate-pulse">
              {conversation.unread_count}
            </span>
          ) : (
            <span className="text-[10px] font-medium bg-gray-100 text-gray-500 px-1.5 py-0.5 rounded-md shrink-0">
              {conversation.participant_count} members
            </span>
          )}
        </div>
      </div>
    </div>
  );
};

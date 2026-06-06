import type { FC } from 'react';
import type { Message } from '../../store/message-store';

interface MessageBubbleProps {
  message: Message;
  isOwn: boolean;
}

export const MessageBubble: FC<MessageBubbleProps> = ({ message, isOwn }) => {
  const formatTime = (isoString: string) => {
    try {
      const date = new Date(isoString);
      return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    } catch {
      return '';
    }
  };

  return (
    <div className={`flex w-full mb-3.5 ${isOwn ? 'justify-end' : 'justify-start'}`}>
      <div
        className={`max-w-[70%] rounded-2xl px-4 py-2.5 shadow-sm relative transition-all border
          ${
            isOwn
              ? 'bg-primary-600 border-primary-700 text-white rounded-br-none'
              : 'bg-white border-gray-150 text-gray-800 rounded-bl-none'
          }`}
      >
        <p className="text-sm leading-relaxed break-words whitespace-pre-wrap">
          {message.content}
        </p>
        <span
          className={`text-[9px] block text-right mt-1 font-medium select-none
            ${isOwn ? 'text-primary-200' : 'text-gray-400'}`}
        >
          {formatTime(message.sent_at)}
        </span>
      </div>
    </div>
  );
};

import React, { useState, useRef, useEffect } from 'react';
import type { FC } from 'react';

interface MessageInputProps {
  onSend: (content: string) => Promise<void>;
  disabled: boolean;
  onTyping?: (isTyping: boolean) => void;
}

export const MessageInput: FC<MessageInputProps> = ({ onSend, disabled, onTyping }) => {
  const [content, setContent] = useState('');
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const typingTimeoutRef = useRef<number | null>(null);
  const isTypingRef = useRef<boolean>(false);

  useEffect(() => {
    return () => {
      if (typingTimeoutRef.current) {
        window.clearTimeout(typingTimeoutRef.current);
      }
    };
  }, []);

  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto';
      textareaRef.current.style.height = `${Math.min(textareaRef.current.scrollHeight, 120)}px`;
    }
  }, [content]);

  const handleSend = async () => {
    const trimmed = content.trim();
    if (!trimmed || disabled) return;

    if (isTypingRef.current) {
      isTypingRef.current = false;
      onTyping?.(false);
      if (typingTimeoutRef.current) {
        window.clearTimeout(typingTimeoutRef.current);
        typingTimeoutRef.current = null;
      }
    }

    try {
      await onSend(trimmed);
      setContent('');
      // Refocus the textarea after DOM updates so the user can keep typing immediately
      setTimeout(() => {
        textareaRef.current?.focus();
      }, 20);
    } catch (err) {
      console.error('Failed to send:', err);
    }
  };

  const handleInputChange = (val: string) => {
    setContent(val);
    if (!onTyping) return;

    if (val.trim() === '') {
      if (isTypingRef.current) {
        isTypingRef.current = false;
        onTyping(false);
        if (typingTimeoutRef.current) {
          window.clearTimeout(typingTimeoutRef.current);
          typingTimeoutRef.current = null;
        }
      }
      return;
    }

    if (!isTypingRef.current) {
      isTypingRef.current = true;
      onTyping(true);
    }

    if (typingTimeoutRef.current) {
      window.clearTimeout(typingTimeoutRef.current);
    }

    typingTimeoutRef.current = window.setTimeout(() => {
      isTypingRef.current = false;
      onTyping(false);
      typingTimeoutRef.current = null;
    }, 2000);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const remainingChars = 2000 - content.length;

  return (
    <div className="p-4 bg-white border-t border-gray-200 shrink-0 flex flex-col gap-1.5 relative z-10">
      <div className="flex items-end gap-3">
        <div className="flex-1 relative">
          <textarea
            ref={textareaRef}
            rows={1}
            maxLength={2000}
            value={content}
            onChange={(e) => handleInputChange(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Type a message..."
            disabled={disabled}
            className="w-full pl-4 pr-10 py-2.5 border border-gray-300 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition-all resize-none max-h-32 disabled:bg-gray-50 disabled:text-gray-400"
          />
          {content.length > 0 && (
            <span
              className={`absolute right-3.5 bottom-3 text-[10px] font-semibold select-none transition-colors
                ${remainingChars < 200 ? 'text-red-500' : 'text-gray-400'}`}
            >
              {remainingChars}
            </span>
          )}
        </div>

        <button
          onClick={handleSend}
          disabled={disabled || !content.trim()}
          className="h-10 px-4 bg-primary-600 hover:bg-primary-700 disabled:bg-gray-100 disabled:text-gray-400 text-white rounded-xl text-sm font-semibold flex items-center justify-center gap-1.5 transition-all active:scale-95 disabled:scale-100 disabled:cursor-not-allowed shadow-sm shrink-0"
        >
          <span>Send</span>
          <svg
            className="h-4.5 w-4.5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            strokeWidth="2"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M14 5l7 7m0 0l-7 7m7-7H3"
            />
          </svg>
        </button>
      </div>
    </div>
  );
};

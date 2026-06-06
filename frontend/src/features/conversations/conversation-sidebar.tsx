import React, { useState } from 'react';
import type { FC } from 'react';
import { ConversationList } from './conversation-list';
import { useConversationStore } from '../../store/conversation-store';
import { Button } from '../../components/button';
import { Input } from '../../components/input';

export const ConversationSidebar: FC = () => {
  const { conversations, selectedConversation, isLoading, createConversation, setSelectedConversation, fetchConversations } = useConversationStore();
  const [showNewChatModal, setShowNewChatModal] = useState(false);
  const [partnerId, setPartnerId] = useState('');
  const [createLoading, setCreateLoading] = useState(false);
  const [createError, setCreateError] = useState('');
  const [searchQuery, setSearchQuery] = useState('');

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setSearchQuery(value);
    fetchConversations(value);
  };

  const handleCreateChat = async (e: React.FormEvent) => {
    e.preventDefault();
    setCreateError('');
    
    const parsedId = parseInt(partnerId, 10);
    if (isNaN(parsedId) || parsedId <= 0) {
      setCreateError('Invalid User ID');
      return;
    }

    setCreateLoading(true);
    try {
      const newConv = await createConversation([parsedId]);
      setSelectedConversation(newConv);
      setShowNewChatModal(false);
      setPartnerId('');
    } catch (err: any) {
      setCreateError(err.response?.data?.message || 'Failed to start conversation. Does the user exist?');
    } finally {
      setCreateLoading(false);
    }
  };

  return (
    <aside className="w-full md:w-80 border-r border-gray-200 bg-white flex flex-col shrink-0 h-full relative z-20">
      {/* Sidebar Header */}
      <div className="p-4 flex items-center justify-between shrink-0">
        <h3 className="font-bold text-lg text-gray-800">Conversations</h3>
        <button
          onClick={() => setShowNewChatModal(true)}
          className="h-9 w-9 rounded-lg bg-primary-50 hover:bg-primary-100 flex items-center justify-center text-primary-600 transition-colors shadow-sm active:scale-95"
          title="New Chat"
        >
          <svg
            className="h-5.5 w-5.5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            strokeWidth="2"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M12 4v16m8-8H4"
            />
          </svg>
        </button>
      </div>

      {/* Search Bar */}
      <div className="px-4 pb-3 shrink-0">
        <div className="relative">
          <input
            type="text"
            className="w-full pl-9 pr-4 py-2 border border-gray-200 rounded-lg text-xs focus:outline-none focus:ring-1 focus:ring-primary-500 focus:border-primary-500 bg-gray-50 text-gray-800"
            placeholder="Search conversations..."
            value={searchQuery}
            onChange={handleSearchChange}
          />
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400">
            <svg
              className="h-4 w-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              strokeWidth="2"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
              />
            </svg>
          </div>
        </div>
      </div>

      {/* Conversations List Scrollable */}
      <div className="flex-1 overflow-y-auto p-3">
        <ConversationList
          conversations={conversations}
          selectedId={selectedConversation?.id || null}
          onSelect={setSelectedConversation}
          isLoading={isLoading}
        />
      </div>

      {/* Inline Modal Overlay */}
      {showNewChatModal && (
        <div className="absolute inset-0 bg-white/95 z-30 p-6 flex flex-col justify-start">
          <div className="flex justify-between items-center mb-6">
            <h4 className="font-bold text-lg text-gray-800">New Conversation</h4>
            <button
              onClick={() => {
                setShowNewChatModal(false);
                setCreateError('');
                setPartnerId('');
              }}
              className="text-gray-400 hover:text-gray-600 transition-colors"
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
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>

          {createError && (
            <div className="p-3 bg-red-50 border border-red-200 text-red-700 text-xs rounded-lg mb-4 font-medium">
              {createError}
            </div>
          )}

          <form onSubmit={handleCreateChat} className="flex flex-col gap-4">
            <Input
              id="partnerId"
              label="User ID to Chat With"
              type="text"
              pattern="[0-9]*"
              inputMode="numeric"
              placeholder="e.g. 2 atau 11"
              value={partnerId}
              onChange={(e) => setPartnerId(e.target.value)}
              required
              autoFocus
            />

            <div className="flex gap-3 mt-2">
              <Button
                variant="outline"
                onClick={() => {
                  setShowNewChatModal(false);
                  setCreateError('');
                  setPartnerId('');
                }}
                className="flex-1"
                type="button"
              >
                Cancel
              </Button>
              <Button
                variant="primary"
                isLoading={createLoading}
                className="flex-1"
                type="submit"
              >
                Start Chat
              </Button>
            </div>
          </form>
        </div>
      )}
    </aside>
  );
};

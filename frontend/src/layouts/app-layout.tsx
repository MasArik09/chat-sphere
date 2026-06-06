import React from 'react';
import { Outlet } from 'react-router-dom';
import { useAuth } from '../features/auth/use-auth';

export const AppLayout: React.FC = () => {
  const { user, logout } = useAuth();

  return (
    <div className="h-screen bg-gray-50 flex flex-col font-sans overflow-hidden">
      {/* Header Navigation */}
      <header className="bg-white border-b border-gray-150 h-16 px-6 flex items-center justify-between shadow-sm shrink-0 z-30">
        <div className="flex items-center gap-3">
          <div className="h-9 w-9 rounded-xl bg-primary-600 flex items-center justify-center shadow-md shadow-primary-500/20">
            <svg
              className="h-5 w-5 text-white"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              strokeWidth="2.5"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
              />
            </svg>
          </div>
          <span className="font-bold text-xl text-gray-800 tracking-tight">ChatSphere</span>
        </div>

        {user && (
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-3">
              <div className="h-9 w-9 rounded-full bg-primary-100 flex items-center justify-center font-bold text-primary-700 text-sm">
                {user.name.charAt(0).toUpperCase()}
              </div>
              <div className="hidden md:flex flex-col text-left">
                <span className="text-sm font-semibold text-gray-800 leading-tight">{user.name}</span>
                <span className="text-xs text-gray-500">{user.email}</span>
              </div>
            </div>
            
            <button
              onClick={logout}
              className="inline-flex items-center gap-1.5 px-3 py-1.5 border border-gray-200 hover:bg-gray-50 text-gray-600 rounded-lg text-sm transition-all hover:text-red-600 hover:border-red-200 active:scale-[0.98]"
            >
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
                  d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                />
              </svg>
              Logout
            </button>
          </div>
        )}
      </header>

      {/* Main Content Area */}
      <main className="flex-1 flex overflow-hidden relative">
        <Outlet />
      </main>
    </div>
  );
};

import React from 'react';
import { Outlet } from 'react-router-dom';

export const AuthLayout: React.FC = () => {
  return (
    <div className="min-h-screen w-full flex items-center justify-center bg-gradient-to-tr from-violet-600 via-indigo-600 to-cyan-500 p-4 md:p-8 relative overflow-hidden">
      {/* Decorative backdrop blobs */}
      <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-white/10 rounded-full blur-3xl" />
      <div className="absolute bottom-[-10%] right-[-10%] w-[45%] h-[45%] bg-white/10 rounded-full blur-3xl" />

      <div className="w-full max-w-md relative z-10">
        <div className="flex flex-col items-center mb-8">
          <div className="h-12 w-12 rounded-2xl bg-white flex items-center justify-center shadow-lg shadow-indigo-600/30 mb-3">
            <svg
              className="h-7 w-7 text-indigo-600"
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
          <h1 className="text-3xl font-bold text-white tracking-tight">ChatSphere</h1>
          <p className="text-indigo-100 text-sm mt-1">Realtime Messaging Platform</p>
        </div>

        <Outlet />
      </div>
    </div>
  );
};

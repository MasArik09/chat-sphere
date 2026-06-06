import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { PublicRoute } from './public-route';
import { ProtectedRoute } from './protected-route';
import { AuthLayout } from '../layouts/auth-layout';
import { AppLayout } from '../layouts/app-layout';
import { LoginPage } from '../pages/login-page';
import { RegisterPage } from '../pages/register-page';
import { ChatPage } from '../pages/chat-page';
import { NotFoundPage } from '../pages/not-found-page';

export const AppRouter: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        {/* Base route redirection */}
        <Route path="/" element={<Navigate to="/chat" replace />} />

        {/* Public auth routes */}
        <Route element={<PublicRoute />}>
          <Route element={<AuthLayout />}>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
          </Route>
        </Route>

        {/* Protected application routes */}
        <Route element={<ProtectedRoute />}>
          <Route element={<AppLayout />}>
            <Route path="/chat" element={<ChatPage />} />
          </Route>
        </Route>

        {/* Fallback 404 Route */}
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </BrowserRouter>
  );
};

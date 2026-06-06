import React from 'react';

export interface AlertProps {
  type?: 'success' | 'error' | 'warning' | 'info';
  message: string;
  className?: string;
}

export const Alert: React.FC<AlertProps> = ({ type = 'info', message, className = '' }) => {
  const styles = {
    success: 'bg-green-50 border-green-200 text-green-800',
    error: 'bg-red-50 border-red-200 text-red-800',
    warning: 'bg-yellow-50 border-yellow-200 text-yellow-800',
    info: 'bg-blue-50 border-blue-200 text-blue-800',
  };

  return (
    <div
      className={`p-4 border rounded-lg text-sm font-medium ${styles[type]} ${className}`}
      role="alert"
    >
      {message}
    </div>
  );
};

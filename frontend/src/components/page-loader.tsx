import React from 'react';
import { Spinner } from './spinner';

export const PageLoader: React.FC = () => {
  return (
    <div className="fixed inset-0 flex items-center justify-center bg-gray-50/70 z-50 backdrop-blur-[2px]">
      <Spinner size="lg" />
    </div>
  );
};

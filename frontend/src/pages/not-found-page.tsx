import React from 'react';
import { Link } from 'react-router-dom';
import { Button } from '../components/button';

export const NotFoundPage: React.FC = () => {
  return (
    <div className="min-h-screen w-full flex flex-col items-center justify-center bg-gray-50 p-4">
      <div className="text-center max-w-md">
        <h1 className="text-9xl font-extrabold text-primary-600 tracking-widest">404</h1>
        <div className="bg-white px-4 py-2 border border-gray-150 rounded-lg shadow-sm text-sm font-semibold text-gray-500 uppercase -mt-4 relative z-10 inline-block mb-8">
          Page Not Found
        </div>
        <p className="text-gray-500 mb-6">
          The page you are looking for does not exist or has been moved.
        </p>
        <Link to="/">
          <Button variant="primary">
            Go Back Home
          </Button>
        </Link>
      </div>
    </div>
  );
};

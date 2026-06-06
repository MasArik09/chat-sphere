import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../features/auth/use-auth';
import { Card } from '../components/card';
import { Input } from '../components/input';
import { Button } from '../components/button';
import { Alert } from '../components/alert';

export const LoginPage: React.FC = () => {
  const { login, error, clearError } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [formError, setFormError] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError('');
    clearError();

    if (!email || !password) {
      setFormError('Please fill in all fields');
      return;
    }

    setIsLoading(true);
    try {
      await login({ email, password });
      navigate('/chat');
    } catch (err) {
      // Error handled by AuthProvider context state
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card className="backdrop-blur-md bg-white/95 border border-white/20 shadow-2xl rounded-2xl p-8">
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-gray-900">Welcome back</h2>
        <p className="text-gray-500 text-sm mt-1">Please sign in to your account</p>
      </div>

      {(formError || error) && (
        <Alert
          type="error"
          message={formError || error || ''}
          className="mb-5"
        />
      )}

      <form onSubmit={handleSubmit} className="flex flex-col gap-5">
        <Input
          id="email"
          label="Email address"
          type="email"
          placeholder="you@example.com"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />

        <Input
          id="password"
          label="Password"
          type="password"
          placeholder="••••••••"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />

        <Button
          type="submit"
          isLoading={isLoading}
          className="w-full mt-2"
        >
          Sign In
        </Button>
      </form>

      <div className="mt-6 text-center text-sm text-gray-500">
        Don't have an account?{' '}
        <Link
          to="/register"
          className="font-medium text-primary-600 hover:text-primary-700 transition-colors"
          onClick={clearError}
        >
          Sign up
        </Link>
      </div>
    </Card>
  );
};

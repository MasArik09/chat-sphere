import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../features/auth/use-auth';
import { Card } from '../components/card';
import { Input } from '../components/input';
import { Button } from '../components/button';
import { Alert } from '../components/alert';

export const RegisterPage: React.FC = () => {
  const { register, error, clearError } = useAuth();
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [formError, setFormError] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError('');
    clearError();

    if (!name || !email || !password) {
      setFormError('Please fill in all fields');
      return;
    }

    if (password.length < 8) {
      setFormError('Password must be at least 8 characters long');
      return;
    }

    setIsLoading(true);
    try {
      await register({ name, email, password });
      navigate('/login');
    } catch (err) {
      // Error handled by AuthProvider context state
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card className="backdrop-blur-md bg-white/95 border border-white/20 shadow-2xl rounded-2xl p-8">
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-gray-900">Create an account</h2>
        <p className="text-gray-500 text-sm mt-1">Join ChatSphere and start messaging</p>
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
          id="name"
          label="Full name"
          type="text"
          placeholder="John Doe"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />

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
          placeholder="Min. 8 characters"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />

        <Button
          type="submit"
          isLoading={isLoading}
          className="w-full mt-2"
        >
          Sign Up
        </Button>
      </form>

      <div className="mt-6 text-center text-sm text-gray-500">
        Already have an account?{' '}
        <Link
          to="/login"
          className="font-medium text-primary-600 hover:text-primary-700 transition-colors"
          onClick={clearError}
        >
          Sign in
        </Link>
      </div>
    </Card>
  );
};

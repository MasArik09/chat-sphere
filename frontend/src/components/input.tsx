import { forwardRef } from 'react';
import type { InputHTMLAttributes } from 'react';

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, className = '', id, ...props }, ref) => {
    return (
      <div className="w-full flex flex-col gap-1.5">
        {label && (
          <label htmlFor={id} className="text-sm font-medium text-gray-700">
            {label}
          </label>
        )}
        <input
          ref={ref}
          id={id}
          className={`px-3.5 py-2 border rounded-lg focus:outline-none focus:ring-2 transition-all text-sm
            ${
              error
                ? 'border-red-300 focus:ring-red-500 focus:border-red-500 bg-red-50/10'
                : 'border-gray-300 focus:ring-primary-500 focus:border-primary-500'
            }
            ${className}`}
          {...props}
        />
        {error && (
          <span className="text-xs text-red-600 mt-0.5" role="alert">
            {error}
          </span>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';

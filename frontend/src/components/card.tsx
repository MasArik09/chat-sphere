import type { HTMLAttributes, FC } from 'react';

export interface CardProps extends HTMLAttributes<HTMLDivElement> {}

export const Card: FC<CardProps> = ({ children, className = '', ...props }) => {
  return (
    <div
      className={`bg-white rounded-xl border border-gray-100 shadow-sm p-6 ${className}`}
      {...props}
    >
      {children}
    </div>
  );
};

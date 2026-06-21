import React from 'react';
import { Loader2 } from 'lucide-react';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  isLoading?: boolean;
  children: React.ReactNode;
}

const variantStyles = {
  primary:
    'bg-accent text-white hover:bg-accent-hover active:bg-accent-muted',
  secondary:
    'bg-surface-card text-zinc-300 border border-zinc-800 hover:bg-surface-hover hover:text-zinc-100 active:bg-zinc-800',
  ghost:
    'text-zinc-500 hover:text-zinc-300 hover:bg-surface-hover active:bg-zinc-800',
  danger:
    'text-red-400 hover:text-red-300 hover:bg-red-500/10 border border-transparent hover:border-red-500/20',
};

const sizeStyles = {
  sm: 'px-2.5 py-1.5 text-xs',
  md: 'px-3.5 py-2 text-sm',
  lg: 'px-5 py-2.5 text-sm',
};

export function Button({
  variant = 'primary',
  size = 'md',
  isLoading,
  children,
  className = '',
  disabled,
  ...props
}: ButtonProps) {
  return (
    <button
      className={`inline-flex items-center justify-center gap-2 rounded-lg font-medium 
        transition-all duration-150 select-none
        ${variantStyles[variant]} ${sizeStyles[size]} 
        ${disabled || isLoading ? 'opacity-50 cursor-not-allowed' : ''} 
        focus:outline-none focus:ring-2 focus:ring-accent/40
        ${className}`}
      disabled={disabled || isLoading}
      {...props}
    >
      {isLoading && <Loader2 className="w-3.5 h-3.5 animate-spin" />}
      {children}
    </button>
  );
}

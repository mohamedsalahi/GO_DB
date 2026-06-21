import React from 'react';

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  hint?: string;
}

export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, hint, className = '', ...props }, ref) => {
    return (
      <div className="space-y-1.5">
        {label && (
          <label className="block text-sm font-medium text-zinc-300">
            {label}
          </label>
        )}
        <input
          ref={ref}
          className={`w-full rounded-lg bg-zinc-900 border border-zinc-800 text-zinc-100 
            placeholder:text-zinc-600 transition-colors duration-150
            focus:border-zinc-700 focus:ring-0 focus:outline-none
            px-3.5 py-2.5 text-sm
            ${error ? 'border-red-500/50' : ''} 
            ${className}`}
          {...props}
        />
        {hint && !error && <p className="text-xs text-zinc-600">{hint}</p>}
        {error && <p className="text-xs text-red-400">{error}</p>}
      </div>
    );
  }
);

Input.displayName = 'Input';

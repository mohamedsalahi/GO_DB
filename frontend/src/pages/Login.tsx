import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import { LogIn } from 'lucide-react';
import { Input } from '../components/ui/Input';
import { Button } from '../components/ui/Button';
import { useAuth } from '../contexts/AuthContext';

export function Login() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);
    try {
      await login(email, password);
      navigate('/');
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Login failed');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-surface flex items-center justify-center p-4">
      <motion.div
        initial={{ opacity: 0, y: 4 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
        className="w-full max-w-sm"
      >
        <div className="text-center mb-8">
          <div className="w-10 h-10 rounded-lg bg-accent flex items-center justify-center mx-auto mb-4">
            <span className="text-sm font-bold text-white">G</span>
          </div>
          <h1 className="text-lg font-semibold text-zinc-100">Sign in to GoClean</h1>
          <p className="text-sm text-zinc-600 mt-1">Enter your credentials to continue</p>
        </div>

        <div className="bg-surface-card border border-zinc-800 rounded-xl p-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
              <div className="p-2.5 rounded-lg bg-red-500/5 border border-red-500/10 text-sm text-red-400">
                {error}
              </div>
            )}

            <Input
              label="Email"
              type="email"
              placeholder="you@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              autoFocus
            />

            <Input
              label="Password"
              type="password"
              placeholder="Enter your password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />

            <Button type="submit" isLoading={isLoading} className="w-full">
              <LogIn className="w-3.5 h-3.5" />
              Sign In
            </Button>
          </form>
        </div>

        <p className="text-center text-sm text-zinc-600 mt-6">
          Don't have an account?{' '}
          <Link to="/register" className="text-accent hover:text-accent-hover font-medium transition-colors">
            Create one
          </Link>
        </p>
      </motion.div>
    </div>
  );
}

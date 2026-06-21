import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import { UserPlus } from 'lucide-react';
import { Input } from '../components/ui/Input';
import { Button } from '../components/ui/Button';
import { useAuth } from '../contexts/AuthContext';

export function Register() {
  const { register } = useAuth();
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    if (password.length < 8) {
      setError('Password must be at least 8 characters');
      return;
    }
    setIsLoading(true);
    try {
      await register(name, email, password);
      navigate('/');
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Registration failed');
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
          <h1 className="text-lg font-semibold text-zinc-100">Create your account</h1>
          <p className="text-sm text-zinc-600 mt-1">Get started with GoClean</p>
        </div>

        <div className="bg-surface-card border border-zinc-800 rounded-xl p-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
              <div className="p-2.5 rounded-lg bg-red-500/5 border border-red-500/10 text-sm text-red-400">
                {error}
              </div>
            )}

            <Input
              label="Name"
              placeholder="John Doe"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              autoFocus
            />

            <Input
              label="Email"
              type="email"
              placeholder="you@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />

            <Input
              label="Password"
              type="password"
              placeholder="At least 8 characters"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              minLength={8}
            />

            <Button type="submit" isLoading={isLoading} className="w-full">
              <UserPlus className="w-3.5 h-3.5" />
              Create Account
            </Button>
          </form>
        </div>

        <p className="text-center text-sm text-zinc-600 mt-6">
          Already have an account?{' '}
          <Link to="/login" className="text-accent hover:text-accent-hover font-medium transition-colors">
            Sign in
          </Link>
        </p>
      </motion.div>
    </div>
  );
}

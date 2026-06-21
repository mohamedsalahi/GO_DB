import { motion } from 'framer-motion';
import { Mail, Calendar, Shield } from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';

export function Profile() {
  const { user } = useAuth();

  if (!user) return null;

  const joined = new Date(user.created_at).toLocaleDateString('en-US', {
    year: 'numeric', month: 'long', day: 'numeric',
  });

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <div>
        <h1 className="text-xl font-semibold text-zinc-100">Profile</h1>
        <p className="text-sm text-zinc-600 mt-0.5">Your account information</p>
      </div>

      <motion.div
        initial={{ opacity: 0, y: 4 }}
        animate={{ opacity: 1, y: 0 }}
        className="bg-surface-card border border-zinc-800 rounded-xl overflow-hidden"
      >
        <div className="h-20 bg-gradient-to-r from-accent/5 to-accent/10" />

        <div className="px-6 pb-6">
          <div className="flex items-end gap-4 -mt-8 mb-5">
            <div className="w-16 h-16 rounded-xl bg-zinc-800 border-2 border-surface flex items-center justify-center text-xl font-semibold text-zinc-400">
              {user.name.charAt(0).toUpperCase()}
            </div>
            <div className="pb-1">
              <h2 className="text-base font-semibold text-zinc-100">{user.name}</h2>
              <div className="flex items-center gap-1.5 mt-0.5">
                <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                <span className="text-xs text-zinc-600">Active</span>
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div className="flex items-center gap-3 p-3 rounded-lg bg-zinc-900/50">
              <div className="w-8 h-8 rounded-lg bg-accent/10 flex items-center justify-center">
                <Mail className="w-4 h-4 text-accent" />
              </div>
              <div className="min-w-0">
                <p className="text-xs text-zinc-600 uppercase tracking-wider font-medium">Email</p>
                <p className="text-sm text-zinc-300 truncate">{user.email}</p>
              </div>
            </div>
            <div className="flex items-center gap-3 p-3 rounded-lg bg-zinc-900/50">
              <div className="w-8 h-8 rounded-lg bg-zinc-800 flex items-center justify-center">
                <Calendar className="w-4 h-4 text-zinc-500" />
              </div>
              <div>
                <p className="text-xs text-zinc-600 uppercase tracking-wider font-medium">Joined</p>
                <p className="text-sm text-zinc-300">{joined}</p>
              </div>
            </div>
          </div>
        </div>
      </motion.div>

      <div className="bg-surface-card border border-zinc-800 rounded-xl p-5">
        <div className="flex items-start gap-3">
          <Shield className="w-4 h-4 text-zinc-600 mt-0.5" />
          <div>
            <p className="text-sm font-medium text-zinc-300">GoClean API v1.0.0</p>
            <p className="text-xs text-zinc-600 mt-0.5">
              Authenticated via JWT. Data stored in PostgreSQL with Redis-backed rate limiting.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

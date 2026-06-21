import { useState } from 'react';
import { NavLink, Outlet } from 'react-router-dom';
import { motion } from 'framer-motion';
import {
  LayoutDashboard,
  CheckSquare,
  User,
  Eye,
  LogOut,
  Menu,
  X,
} from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';

function useNavItems() {
  const { user } = useAuth();
  const items = [
    { to: '/', icon: LayoutDashboard, label: 'Dashboard' },
    { to: '/tasks', icon: CheckSquare, label: 'Tasks' },
    { to: '/profile', icon: User, label: 'Profile' },
  ];
  if (user?.role === 'admin') {
    items.push({ to: '/god-eye', icon: Eye, label: "God's Eye" });
  }
  return items;
}

export function Layout() {
  const { user, logout } = useAuth();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const navItems = useNavItems();

  return (
    <div className="min-h-screen bg-surface flex">
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-40 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      <aside
        className={`fixed lg:sticky top-0 lg:flex lg:flex-col h-screen z-50 w-56 bg-surface border-r border-zinc-800/50
          transition-transform duration-200 lg:translate-x-0 ${
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
      >
        <div className="flex items-center justify-between h-12 px-4 border-b border-zinc-800/50">
          <div className="flex items-center gap-2.5">
            <div className="w-6 h-6 rounded-md bg-accent flex items-center justify-center">
              <span className="text-[10px] font-bold text-white tracking-tight">G</span>
            </div>
            <span className="text-sm font-semibold text-zinc-100 tracking-tight">GoClean</span>
          </div>
          <button
            className="lg:hidden p-1 rounded text-zinc-600 hover:text-zinc-300"
            onClick={() => setSidebarOpen(false)}
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        <nav className="flex-1 p-2 space-y-0.5">
          {navItems.map(({ to, icon: Icon, label }) => (
            <NavLink
              key={to}
              to={to}
              end={to === '/'}
              onClick={() => setSidebarOpen(false)}
              className={({ isActive }) =>
                `flex items-center gap-2.5 px-3 py-2 rounded-md text-sm font-medium transition-colors duration-150
                ${isActive
                  ? 'bg-accent/10 text-accent-hover'
                  : 'text-zinc-500 hover:text-zinc-300 hover:bg-zinc-800/50'
                }`
              }
            >
              <Icon className="w-4 h-4" strokeWidth={1.5} />
              {label}
            </NavLink>
          ))}
        </nav>

        <div className="p-2 border-t border-zinc-800/50">
          <div className="flex items-center gap-2.5 px-3 py-2 rounded-md">
            <div className="w-7 h-7 rounded-md bg-zinc-800 flex items-center justify-center text-xs font-semibold text-zinc-400">
              {user?.name?.charAt(0)?.toUpperCase() || '?'}
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-zinc-300 truncate leading-tight">{user?.name || 'User'}</p>
              <p className="text-xs text-zinc-600 truncate">{user?.email}</p>
            </div>
            <button
              onClick={logout}
              className="p-1 rounded text-zinc-600 hover:text-red-400 hover:bg-zinc-800 transition-colors"
              title="Logout"
            >
              <LogOut className="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      </aside>

      <div className="flex-1 flex flex-col min-w-0">
        <header className="sticky top-0 z-30 bg-surface/80 backdrop-blur-md border-b border-zinc-800/50">
          <div className="flex items-center px-4 lg:px-6 h-12">
            <button
              className="lg:hidden p-1.5 rounded-md text-zinc-500 hover:text-zinc-300 hover:bg-zinc-800"
              onClick={() => setSidebarOpen(true)}
            >
              <Menu className="w-4 h-4" />
            </button>
          </div>
        </header>

        <main className="flex-1 p-6 lg:p-8">
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.2 }}
          >
            <Outlet />
          </motion.div>
        </main>
      </div>
    </div>
  );
}

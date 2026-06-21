import { useState, useEffect, useCallback, useMemo } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Eye, Users, ListTodo, Shield, ChevronDown, AlertCircle, Inbox, Search } from 'lucide-react';
import { adminService } from '../services/adminService';
import { Button } from '../components/ui/Button';
import { Spinner, Badge } from '../components/ui/Badge';
import type { User, Task } from '../types';

interface UserWithTasks extends User {
  tasks?: Task[];
  tasksLoading?: boolean;
}

function RoleBadge({ role }: { role?: string }) {
  const isAdmin = role === 'admin';
  return (
    <span
      className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-2xs font-medium ${
        isAdmin
          ? 'text-amber-400 bg-amber-500/5'
          : 'text-zinc-500 bg-zinc-800/50'
      }`}
    >
      <span className={`w-1 h-1 rounded-full ${isAdmin ? 'bg-amber-400' : 'bg-zinc-600'}`} />
      {role || 'user'}
    </span>
  );
}

function formatDate(dateStr: string) {
  const d = new Date(dateStr);
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
}

function UserRow({
  user,
  isExpanded,
  onToggle,
  onPromote,
  promoting,
}: {
  user: UserWithTasks;
  isExpanded: boolean;
  onToggle: () => void;
  onPromote: () => void;
  promoting: boolean;
}) {
  return (
    <>
      <tr
        className="group cursor-pointer transition-colors hover:bg-zinc-800/30"
        onClick={onToggle}
      >
        <td className="py-3 px-4 text-sm">
          <div className="flex items-center gap-2.5">
            <div className="w-7 h-7 rounded-md bg-zinc-800 flex items-center justify-center text-xs font-semibold text-zinc-400 group-hover:bg-zinc-700 transition-colors">
              {user.name.charAt(0).toUpperCase()}
            </div>
            <span className="text-zinc-200 font-medium">{user.name}</span>
          </div>
        </td>
        <td className="py-3 px-4 text-sm text-zinc-400">{user.email}</td>
        <td className="py-3 px-4">
          <RoleBadge role={user.role} />
        </td>
        <td className="py-3 px-4 text-sm text-zinc-500 tabular-nums">{formatDate(user.created_at)}</td>
        <td className="py-3 px-4 text-right">
          <div className="flex items-center justify-end gap-2" onClick={(e) => e.stopPropagation()}>
            {user.role !== 'admin' && (
              <Button
                variant="ghost"
                size="sm"
                isLoading={promoting}
                onClick={onPromote}
                className="text-amber-500/70 hover:text-amber-400 hover:bg-amber-500/5"
              >
                <Shield className="w-3 h-3" />
                Promote
              </Button>
            )}
            <ChevronDown
              className={`w-4 h-4 text-zinc-600 transition-transform duration-200 ${
                isExpanded ? 'rotate-180' : ''
              }`}
            />
          </div>
        </td>
      </tr>
      <AnimatePresence>
        {isExpanded && (
          <tr key={`${user.id}-tasks`}>
            <td colSpan={5} className="px-0 py-0">
              <motion.div
                initial={{ height: 0, opacity: 0 }}
                animate={{ height: 'auto', opacity: 1 }}
                exit={{ height: 0, opacity: 0 }}
                transition={{ duration: 0.2, ease: 'easeInOut' }}
                className="overflow-hidden"
              >
                <div className="border-t border-zinc-800/50 bg-zinc-900/30 px-4 py-4">
                  {user.tasksLoading ? (
                    <div className="flex items-center gap-2.5 py-2">
                      <Spinner className="w-4 h-4" />
                      <span className="text-xs text-zinc-600">Loading tasks...</span>
                    </div>
                  ) : !user.tasks || user.tasks.length === 0 ? (
                    <div className="flex items-center gap-2 py-2">
                      <Inbox className="w-3.5 h-3.5 text-zinc-700" />
                      <span className="text-xs text-zinc-600">No tasks for this user</span>
                    </div>
                  ) : (
                    <div className="space-y-1">
                      <p className="text-2xs text-zinc-700 font-medium uppercase tracking-wider mb-2">
                        Tasks ({user.tasks.length})
                      </p>
                      {user.tasks.map((task) => (
                        <div
                          key={task.id}
                          className="flex items-center justify-between py-1.5 px-3 rounded-md bg-zinc-800/20 border border-zinc-800/30"
                        >
                          <div className="flex items-center gap-2.5 min-w-0">
                            <span className="text-sm text-zinc-300 truncate">{task.title}</span>
                          </div>
                          <Badge status={task.status} />
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </motion.div>
            </td>
          </tr>
        )}
      </AnimatePresence>
    </>
  );
}

export function GodEye() {
  const [users, setUsers] = useState<UserWithTasks[]>([]);
  const [allTasks, setAllTasks] = useState<Task[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [expandedUserId, setExpandedUserId] = useState<string | null>(null);
  const [promotingId, setPromotingId] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [activeTab, setActiveTab] = useState<'users' | 'tasks'>('users');

  const fetchData = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const [usersData, tasksData] = await Promise.all([
        adminService.listUsers(),
        adminService.listAllTasks(),
      ]);
      setUsers(usersData);
      setAllTasks(tasksData);
    } catch (err: any) {
      setError(err?.response?.data?.error || err.message || 'Failed to load admin data');
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const userIndex = useMemo(() => {
    const map = new Map<string, User>();
    for (const u of users) {
      map.set(u.id, u);
    }
    return map;
  }, [users]);

  const indexedTasks = useMemo(() => {
    const tasks = allTasks.map((t) => ({
      ...t,
      userName: userIndex.get(t.user_id)?.name || 'Unknown',
      userEmail: userIndex.get(t.user_id)?.email || '',
    }));

    if (!searchQuery.trim()) return tasks;

    const q = searchQuery.toLowerCase();
    return tasks.filter(
      (t) =>
        t.title.toLowerCase().includes(q) ||
        (t.description && t.description.toLowerCase().includes(q)) ||
        t.userName.toLowerCase().includes(q) ||
        t.userEmail.toLowerCase().includes(q),
    );
  }, [allTasks, userIndex, searchQuery]);

  async function handleToggleUser(userId: string) {
    if (expandedUserId === userId) {
      setExpandedUserId(null);
      return;
    }
    setExpandedUserId(userId);

    const user = users.find((u) => u.id === userId);
    if (!user || user.tasks) return;

    setUsers((prev) => prev.map((u) => (u.id === userId ? { ...u, tasksLoading: true } : u)));
    try {
      const tasks = await adminService.listUserTasks(userId);
      setUsers((prev) => prev.map((u) => (u.id === userId ? { ...u, tasks, tasksLoading: false } : u)));
    } catch {
      setUsers((prev) => prev.map((u) => (u.id === userId ? { ...u, tasks: [], tasksLoading: false } : u)));
    }
  }

  async function handlePromote(userId: string) {
    setPromotingId(userId);
    try {
      const updated = await adminService.promoteUser(userId);
      setUsers((prev) => prev.map((u) => (u.id === userId ? { ...u, ...updated } : u)));
    } catch {
      // error handled silently
    } finally {
      setPromotingId(null);
    }
  }

  if (isLoading) {
    return (
      <div className="max-w-5xl mx-auto space-y-6">
        <div className="flex items-center gap-3 py-8">
          <div className="w-8 h-8 rounded-lg bg-zinc-800 animate-pulse" />
          <div className="space-y-2">
            <div className="h-4 w-28 bg-zinc-800 rounded animate-pulse" />
            <div className="h-3 w-40 bg-zinc-800/50 rounded animate-pulse" />
          </div>
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div className="h-24 bg-surface-card border border-zinc-800 rounded-xl animate-pulse" />
          <div className="h-24 bg-surface-card border border-zinc-800 rounded-xl animate-pulse" />
        </div>
        <div className="bg-surface-card border border-zinc-800 rounded-xl p-6 space-y-4">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="flex items-center gap-3">
              <div className="w-7 h-7 rounded-md bg-zinc-800 animate-pulse" />
              <div className="flex-1 space-y-1.5">
                <div className="h-3 w-32 bg-zinc-800 rounded animate-pulse" />
                <div className="h-2.5 w-48 bg-zinc-800/50 rounded animate-pulse" />
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-5xl mx-auto">
        <div className="bg-surface-card border border-zinc-800 rounded-xl p-12 text-center">
          <div className="w-12 h-12 rounded-xl bg-red-500/5 flex items-center justify-center mx-auto mb-4">
            <AlertCircle className="w-6 h-6 text-red-400" />
          </div>
          <p className="text-sm text-zinc-400 font-medium">Failed to load admin panel</p>
          <p className="text-xs text-zinc-600 mt-1 mb-4">{error}</p>
          <Button variant="secondary" size="sm" onClick={fetchData}>
            Retry
          </Button>
        </div>
      </div>
    );
  }

  const totalUsers = users.length;
  const totalTasks = allTasks.length;

  return (
    <div className="max-w-5xl mx-auto space-y-6">
      <motion.div
        initial={{ opacity: 0, y: -8 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
        className="flex items-center justify-between"
      >
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 rounded-lg bg-accent/10 border border-accent/20 flex items-center justify-center">
            <Eye className="w-4 h-4 text-accent" />
          </div>
          <div>
            <h1 className="text-lg font-semibold text-zinc-100 tracking-tight">
              God's Eye
            </h1>
            <p className="text-xs text-zinc-600 mt-0.5">Admin oversight panel</p>
          </div>
        </div>
        <Button variant="secondary" size="sm" isLoading={isLoading} onClick={fetchData}>
          Refresh
        </Button>
      </motion.div>

      <div className="grid grid-cols-2 gap-4">
        <motion.div
          initial={{ opacity: 0, y: 8 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3, delay: 0.05 }}
          className="bg-surface-card border border-zinc-800 rounded-xl p-4"
        >
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-lg bg-indigo-500/5 border border-indigo-500/10 flex items-center justify-center">
              <Users className="w-4 h-4 text-indigo-400" />
            </div>
            <div>
              <p className="text-2xs text-zinc-600 font-medium uppercase tracking-wider">Total Users</p>
              <p className="text-xl font-semibold text-zinc-100 tabular-nums mt-0.5">{totalUsers}</p>
            </div>
          </div>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 8 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3, delay: 0.1 }}
          className="bg-surface-card border border-zinc-800 rounded-xl p-4"
        >
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-lg bg-emerald-500/5 border border-emerald-500/10 flex items-center justify-center">
              <ListTodo className="w-4 h-4 text-emerald-400" />
            </div>
            <div>
              <p className="text-2xs text-zinc-600 font-medium uppercase tracking-wider">Total Tasks</p>
              <p className="text-xl font-semibold text-zinc-100 tabular-nums mt-0.5">{totalTasks}</p>
            </div>
          </div>
        </motion.div>
      </div>

      <div className="flex items-center gap-2 border-b border-zinc-800/50">
        <button
          onClick={() => setActiveTab('users')}
          className={`px-3 py-2 text-xs font-medium transition-colors border-b-2 -mb-[1px] ${
            activeTab === 'users'
              ? 'text-zinc-200 border-accent'
              : 'text-zinc-600 border-transparent hover:text-zinc-400'
          }`}
        >
          <Users className="w-3.5 h-3.5 inline mr-1.5" />
          Users
        </button>
        <button
          onClick={() => setActiveTab('tasks')}
          className={`px-3 py-2 text-xs font-medium transition-colors border-b-2 -mb-[1px] ${
            activeTab === 'tasks'
              ? 'text-zinc-200 border-accent'
              : 'text-zinc-600 border-transparent hover:text-zinc-400'
          }`}
        >
          <ListTodo className="w-3.5 h-3.5 inline mr-1.5" />
          All Tasks
          {allTasks.length > 0 && (
            <span className="ml-1.5 text-2xs text-zinc-600">({allTasks.length})</span>
          )}
        </button>
      </div>

      {activeTab === 'users' ? (
        users.length === 0 ? (
          <div className="bg-surface-card border border-zinc-800 rounded-xl p-12 text-center">
            <div className="w-10 h-10 rounded-lg bg-zinc-800 flex items-center justify-center mx-auto mb-3">
              <Users className="w-5 h-5 text-zinc-600" />
            </div>
            <p className="text-sm text-zinc-400 font-medium">No users found</p>
            <p className="text-xs text-zinc-600 mt-1">No users are registered yet</p>
          </div>
        ) : (
          <motion.div
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3, delay: 0.15 }}
            className="bg-surface-card border border-zinc-800 rounded-xl overflow-hidden"
          >
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-zinc-800/50">
                    <th className="text-left py-3 px-4 text-2xs text-zinc-700 font-medium uppercase tracking-wider">Name</th>
                    <th className="text-left py-3 px-4 text-2xs text-zinc-700 font-medium uppercase tracking-wider">Email</th>
                    <th className="text-left py-3 px-4 text-2xs text-zinc-700 font-medium uppercase tracking-wider">Role</th>
                    <th className="text-left py-3 px-4 text-2xs text-zinc-700 font-medium uppercase tracking-wider">Created</th>
                    <th className="py-3 px-4" />
                  </tr>
                </thead>
                <tbody className="divide-y divide-zinc-800/30">
                  {users.map((user) => (
                    <UserRow
                      key={user.id}
                      user={user}
                      isExpanded={expandedUserId === user.id}
                      onToggle={() => handleToggleUser(user.id)}
                      onPromote={() => handlePromote(user.id)}
                      promoting={promotingId === user.id}
                    />
                  ))}
                </tbody>
              </table>
            </div>
          </motion.div>
        )
      ) : (
        <motion.div
          initial={{ opacity: 0, y: 8 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3, delay: 0.15 }}
          className="space-y-3"
        >
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-600" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search tasks by title, description, or user..."
              className="w-full pl-9 pr-3 py-2 bg-surface-card border border-zinc-800 rounded-lg text-sm text-zinc-300 placeholder-zinc-600 focus:outline-none focus:border-accent/40 focus:ring-1 focus:ring-accent/20 transition-colors"
            />
          </div>

          {indexedTasks.length === 0 ? (
            <div className="bg-surface-card border border-zinc-800 rounded-xl p-12 text-center">
              <div className="w-10 h-10 rounded-lg bg-zinc-800 flex items-center justify-center mx-auto mb-3">
                <ListTodo className="w-5 h-5 text-zinc-600" />
              </div>
              <p className="text-sm text-zinc-400 font-medium">
                {searchQuery ? 'No tasks match your search' : 'No tasks found'}
              </p>
              <p className="text-xs text-zinc-600 mt-1">
                {searchQuery ? 'Try a different search term' : 'No tasks have been created yet'}
              </p>
            </div>
          ) : (
            <div className="space-y-1">
              {indexedTasks.map((task) => (
                <div
                  key={task.id}
                  className="flex items-center justify-between py-2 px-4 bg-surface-card border border-zinc-800 rounded-lg hover:border-zinc-700 transition-colors"
                >
                  <div className="flex items-center gap-3 min-w-0">
                    <div className="flex-shrink-0 w-6 h-6 rounded bg-zinc-800 flex items-center justify-center text-2xs font-semibold text-zinc-500">
                      {task.userName.charAt(0).toUpperCase()}
                    </div>
                    <div className="min-w-0">
                      <p className="text-sm text-zinc-300 truncate">{task.title}</p>
                      <p className="text-2xs text-zinc-600 mt-0.5">
                        {task.userName} &middot; {task.userEmail}
                      </p>
                    </div>
                  </div>
                  <Badge status={task.status} />
                </div>
              ))}
            </div>
          )}
        </motion.div>
      )}
    </div>
  );
}

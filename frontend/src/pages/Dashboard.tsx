import { useState, useMemo } from 'react';
import { motion } from 'framer-motion';
import { Plus, ListTodo, CheckCircle2, Clock, TrendingUp } from 'lucide-react';
import { StatsCard } from '../components/StatsCard';
import { TaskCard } from '../components/TaskCard';
import { TaskModal } from '../components/TaskModal';
import { Button } from '../components/ui/Button';
import { Spinner } from '../components/ui/Badge';
import { useTasks } from '../hooks/useTasks';
import type { Task } from '../types';

export function Dashboard() {
  const { tasks, isLoading, createTask, updateTask, deleteTask } = useTasks();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);

  const stats = useMemo(() => ({
    total: tasks.length,
    completed: tasks.filter((t) => t.status === 'completed').length,
    pending: tasks.filter((t) => t.status === 'pending').length,
    inProgress: tasks.filter((t) => t.status === 'in_progress').length,
  }), [tasks]);

  const recentTasks = tasks.slice(0, 5);

  const progress = stats.total > 0 ? Math.round((stats.completed / stats.total) * 100) : 0;

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner className="w-5 h-5" />
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold text-zinc-100">Dashboard</h1>
          <p className="text-sm text-zinc-600 mt-0.5">Overview of your tasks</p>
        </div>
        <Button onClick={() => setIsModalOpen(true)} size="sm">
          <Plus className="w-3.5 h-3.5" />
          New Task
        </Button>
      </div>

      <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
        <StatsCard title="Total" value={stats.total} icon={ListTodo} delay={0} />
        <StatsCard title="Completed" value={stats.completed} icon={CheckCircle2} delay={0.05} />
        <StatsCard title="In Progress" value={stats.inProgress} icon={TrendingUp} delay={0.1} />
        <StatsCard title="Pending" value={stats.pending} icon={Clock} delay={0.15} />
      </div>

      {stats.total > 0 && (
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-xs text-zinc-600 font-medium uppercase tracking-wider">Progress</span>
            <span className="text-xs text-zinc-500 tabular-nums">{progress}%</span>
          </div>
          <div className="h-1.5 bg-zinc-800 rounded-full overflow-hidden">
            <motion.div
              initial={{ width: 0 }}
              animate={{ width: `${progress}%` }}
              transition={{ duration: 0.8, ease: 'easeOut', delay: 0.2 }}
              className="h-full bg-accent rounded-full"
            />
          </div>
        </div>
      )}

      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <h2 className="text-sm font-medium text-zinc-300">Recent tasks</h2>
          {tasks.length > 5 && (
            <a href="/tasks" className="text-xs text-zinc-600 hover:text-zinc-400 transition-colors">
              View all
            </a>
          )}
        </div>

        {recentTasks.length === 0 ? (
          <div className="bg-surface-card border border-zinc-800 rounded-xl p-12 text-center">
            <div className="w-10 h-10 rounded-lg bg-zinc-800 flex items-center justify-center mx-auto mb-3">
              <ListTodo className="w-5 h-5 text-zinc-600" />
            </div>
            <p className="text-sm text-zinc-400 font-medium">No tasks yet</p>
            <p className="text-xs text-zinc-600 mt-1">Create your first task to get started</p>
            <Button className="mt-4" size="sm" onClick={() => setIsModalOpen(true)}>
              <Plus className="w-3.5 h-3.5" />
              Create Task
            </Button>
          </div>
        ) : (
          <div className="space-y-1.5">
            {recentTasks.map((task, i) => (
              <TaskCard
                key={task.id}
                task={task}
                onEdit={(t) => { setEditingTask(t); setIsModalOpen(true); }}
                onDelete={deleteTask}
                index={i}
              />
            ))}
          </div>
        )}
      </div>

      <TaskModal
        isOpen={isModalOpen}
        onClose={() => { setIsModalOpen(false); setEditingTask(null); }}
        onSubmit={async (data) => {
          if (editingTask) await updateTask(editingTask.id, data);
          else await createTask(data as any);
        }}
        task={editingTask}
      />
    </div>
  );
}

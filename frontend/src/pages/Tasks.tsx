import { useState, useMemo } from 'react';
import { motion } from 'framer-motion';
import { Search, Plus, Inbox } from 'lucide-react';
import { TaskCard } from '../components/TaskCard';
import { TaskModal } from '../components/TaskModal';
import { Button } from '../components/ui/Button';
import { Spinner } from '../components/ui/Badge';
import { useTasks } from '../hooks/useTasks';
import type { Task, TaskStatus } from '../types';

const filters: Array<{ label: string; value: TaskStatus | 'all' }> = [
  { label: 'All', value: 'all' },
  { label: 'Pending', value: 'pending' },
  { label: 'In Progress', value: 'in_progress' },
  { label: 'Completed', value: 'completed' },
];

export function Tasks() {
  const { tasks, isLoading, createTask, updateTask, deleteTask } = useTasks();
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState<TaskStatus | 'all'>('all');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);

  const filtered = useMemo(() =>
    tasks.filter((t) => {
      const q = search.toLowerCase();
      const matchSearch = !q || t.title.toLowerCase().includes(q) || (t.description?.toLowerCase() || '').includes(q);
      const matchStatus = statusFilter === 'all' || t.status === statusFilter;
      return matchSearch && matchStatus;
    }),
    [tasks, search, statusFilter]
  );

  return (
    <div className="max-w-4xl mx-auto space-y-5">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold text-zinc-100">Tasks</h1>
          <p className="text-sm text-zinc-600 mt-0.5">{filtered.length} of {tasks.length} tasks</p>
        </div>
        <Button onClick={() => setIsModalOpen(true)} size="sm">
          <Plus className="w-3.5 h-3.5" />
          New Task
        </Button>
      </div>

      <div className="flex flex-col sm:flex-row gap-2">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-zinc-600" />
          <input
            type="text"
            placeholder="Search tasks..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full rounded-lg bg-surface-card border border-zinc-800 text-zinc-100 placeholder:text-zinc-600
              focus:border-zinc-700 focus:ring-0 focus:outline-none pl-9 pr-3.5 py-2 text-sm transition-colors"
          />
        </div>
        <div className="flex gap-1">
          {filters.map((f) => (
            <button
              key={f.value}
              onClick={() => setStatusFilter(f.value)}
              className={`px-2.5 py-1.5 rounded-md text-xs font-medium border transition-colors ${
                statusFilter === f.value
                  ? 'border-accent/30 bg-accent/10 text-accent-hover'
                  : 'border-zinc-800 text-zinc-600 hover:text-zinc-400 hover:border-zinc-700'
              }`}
            >
              {f.label}
            </button>
          ))}
        </div>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center h-64">
          <Spinner className="w-5 h-5" />
        </div>
      ) : filtered.length === 0 ? (
        <div className="bg-surface-card border border-zinc-800 rounded-xl p-16 text-center">
          <div className="w-10 h-10 rounded-lg bg-zinc-800 flex items-center justify-center mx-auto mb-3">
            <Inbox className="w-5 h-5 text-zinc-600" />
          </div>
          <p className="text-sm text-zinc-400 font-medium">
            {search || statusFilter !== 'all' ? 'No matching tasks' : 'No tasks yet'}
          </p>
          <p className="text-xs text-zinc-600 mt-1">
            {search || statusFilter !== 'all' ? 'Try adjusting your search or filters' : 'Create your first task'}
          </p>
          {!search && statusFilter === 'all' && (
            <Button className="mt-4" size="sm" onClick={() => setIsModalOpen(true)}>
              <Plus className="w-3.5 h-3.5" />
              Create Task
            </Button>
          )}
        </div>
      ) : (
        <div className="space-y-1.5">
          {filtered.map((task, i) => (
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

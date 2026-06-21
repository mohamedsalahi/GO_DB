import { motion } from 'framer-motion';
import { Calendar, Trash2, Edit3 } from 'lucide-react';
import { Badge } from './ui/Badge';
import type { Task } from '../types';

interface TaskCardProps {
  task: Task;
  onEdit: (task: Task) => void;
  onDelete: (id: string) => void;
  index: number;
}

export function TaskCard({ task, onEdit, onDelete, index }: TaskCardProps) {
  const date = task.created_at
    ? new Date(task.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
    : '';

  return (
    <motion.div
      initial={{ opacity: 0, y: 6 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.2, delay: index * 0.03 }}
      className="group bg-surface-card border border-zinc-800 rounded-lg px-4 py-3 card-hover cursor-default"
    >
      <div className="flex items-start justify-between gap-3">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <Badge status={task.status} />
          </div>
          <h3 className="text-sm font-medium text-zinc-100 truncate">{task.title}</h3>
          {task.description && (
            <p className="text-xs text-zinc-600 mt-1 line-clamp-2">{task.description}</p>
          )}
          <div className="flex items-center gap-3 mt-2">
            <div className="flex items-center gap-1 text-2xs text-zinc-600">
              <Calendar className="w-3 h-3" />
              {date}
            </div>
          </div>
        </div>
        <div className="flex items-center gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity duration-150 pt-1">
          <button
            onClick={() => onEdit(task)}
            className="p-1.5 rounded text-zinc-600 hover:text-zinc-300 hover:bg-zinc-800 transition-colors"
            title="Edit"
          >
            <Edit3 className="w-3.5 h-3.5" />
          </button>
          <button
            onClick={() => onDelete(task.id)}
            className="p-1.5 rounded text-zinc-600 hover:text-red-400 hover:bg-zinc-800 transition-colors"
            title="Delete"
          >
            <Trash2 className="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
    </motion.div>
  );
}

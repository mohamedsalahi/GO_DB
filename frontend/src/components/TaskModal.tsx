import { useState, useEffect } from 'react';
import { Modal } from './ui/Modal';
import { Input } from './ui/Input';
import { Button } from './ui/Button';
import type { Task, CreateTaskRequest, UpdateTaskRequest, TaskStatus } from '../types';

interface TaskModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: CreateTaskRequest | UpdateTaskRequest) => Promise<void>;
  task?: Task | null;
}

export function TaskModal({ isOpen, onClose, onSubmit, task }: TaskModalProps) {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [status, setStatus] = useState<TaskStatus>('pending');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const isEditing = !!task;

  useEffect(() => {
    if (task) {
      setTitle(task.title);
      setDescription(task.description || '');
      setStatus(task.status);
    } else {
      setTitle('');
      setDescription('');
      setStatus('pending');
    }
  }, [task, isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;
    setIsSubmitting(true);
    try {
      if (isEditing) {
        await onSubmit({ title: title.trim(), description: description.trim() || null, status } as UpdateTaskRequest);
      } else {
        await onSubmit({ title: title.trim(), description: description.trim() || null } as CreateTaskRequest);
      }
      onClose();
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={isEditing ? 'Edit task' : 'New task'}>
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Title"
          placeholder="What needs to be done?"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          autoFocus
        />

        <div className="space-y-1.5">
          <label className="block text-sm font-medium text-zinc-300">Description</label>
          <textarea
            className="w-full rounded-lg bg-zinc-900 border border-zinc-800 text-zinc-100 placeholder:text-zinc-600 
              focus:border-zinc-700 focus:ring-0 focus:outline-none px-3.5 py-2.5 text-sm min-h-[80px] resize-none transition-colors"
            placeholder="Add details (optional)"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
        </div>

        {isEditing && (
          <div className="space-y-1.5">
            <label className="block text-sm font-medium text-zinc-300">Status</label>
            <div className="flex gap-1.5">
              {(['pending', 'in_progress', 'completed'] as TaskStatus[]).map((s) => (
                <button
                  key={s}
                  type="button"
                  onClick={() => setStatus(s)}
                  className={`px-2.5 py-1.5 rounded-md text-xs font-medium border transition-colors ${
                    status === s
                      ? 'border-accent/50 bg-accent/10 text-accent-hover'
                      : 'border-zinc-800 text-zinc-600 hover:text-zinc-400 hover:border-zinc-700'
                  }`}
                >
                  {s === 'pending' ? 'Pending' : s === 'in_progress' ? 'In Progress' : 'Completed'}
                </button>
              ))}
            </div>
          </div>
        )}

        <div className="flex gap-2 pt-2">
          <Button type="button" variant="secondary" onClick={onClose} className="flex-1">
            Cancel
          </Button>
          <Button type="submit" isLoading={isSubmitting} disabled={!title.trim()} className="flex-1">
            {isEditing ? 'Save' : 'Create'}
          </Button>
        </div>
      </form>
    </Modal>
  );
}

import type { TaskStatus } from '../../types';

const styles: Record<TaskStatus, string> = {
  pending: 'text-amber-400 bg-amber-500/5',
  in_progress: 'text-blue-400 bg-blue-500/5',
  completed: 'text-emerald-400 bg-emerald-500/5',
};

const labels: Record<TaskStatus, string> = {
  pending: 'Pending',
  in_progress: 'In Progress',
  completed: 'Completed',
};

export function Badge({ status }: { status: TaskStatus }) {
  return (
    <span className={`inline-flex items-center gap-1.5 px-2 py-0.5 rounded-md text-2xs font-medium ${styles[status]}`}>
      <span className={`w-1 h-1 rounded-full ${status === 'pending' ? 'bg-amber-400' : status === 'in_progress' ? 'bg-blue-400' : 'bg-emerald-400'}`} />
      {labels[status]}
    </span>
  );
}

export function Spinner({ className = '' }: { className?: string }) {
  return (
    <svg className={`animate-spin text-zinc-500 ${className}`} xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="3" />
      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
    </svg>
  );
}

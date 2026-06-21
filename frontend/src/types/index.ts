export interface User {
  id: string;
  name: string;
  email: string;
  role?: string;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

export interface Task {
  id: string;
  user_id: string;
  title: string;
  description: string | null;
  status: TaskStatus;
  due_date: string | null;
  created_at: string;
  updated_at: string;
}

export type TaskStatus = 'pending' | 'in_progress' | 'completed';

export interface CreateTaskRequest {
  title: string;
  description?: string | null;
  due_date?: string | null;
}

export interface UpdateTaskRequest {
  title?: string;
  description?: string | null;
  status?: TaskStatus;
  due_date?: string | null;
}

export interface UserTasksResponse {
  user_id: string;
  tasks: Task[];
}

export interface APIResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

export const TASK_STATUS_LABELS: Record<TaskStatus, string> = {
  pending: 'Pending',
  in_progress: 'In Progress',
  completed: 'Completed',
};

export const TASK_STATUS_COLORS: Record<TaskStatus, string> = {
  pending: 'amber',
  in_progress: 'blue',
  completed: 'emerald',
};

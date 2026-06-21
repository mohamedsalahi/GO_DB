import { useState, useCallback, useEffect } from 'react';
import { tasksAPI } from '../api/client';
import type { Task, CreateTaskRequest, UpdateTaskRequest } from '../types';

export function useTasks() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchTasks = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await tasksAPI.list();
      setTasks(data);
    } catch (err: any) {
      setError(err.message || 'Failed to load tasks');
    } finally {
      setIsLoading(false);
    }
  }, []);

  const createTask = useCallback(async (data: CreateTaskRequest) => {
    const task = await tasksAPI.create(data);
    setTasks((prev) => [task, ...prev]);
    return task;
  }, []);

  const updateTask = useCallback(async (id: string, data: UpdateTaskRequest) => {
    const updated = await tasksAPI.update(id, data);
    setTasks((prev) => prev.map((t) => (t.id === id ? updated : t)));
    return updated;
  }, []);

  const deleteTask = useCallback(async (id: string) => {
    await tasksAPI.delete(id);
    setTasks((prev) => prev.filter((t) => t.id !== id));
  }, []);

  useEffect(() => { fetchTasks(); }, [fetchTasks]);

  return { tasks, isLoading, error, setError, fetchTasks, createTask, updateTask, deleteTask };
}

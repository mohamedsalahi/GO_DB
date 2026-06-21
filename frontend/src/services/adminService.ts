import { api, extractData } from '../api/client';
import type { User, Task, UserTasksResponse } from '../types';

export const adminService = {
  listUsers: () =>
    api.get('/admin/users').then(extractData<User[]>),

  listAllTasks: () =>
    api.get('/admin/tasks').then(extractData<Task[]>),

  listUserTasks: async (userId: string): Promise<Task[]> => {
    const res = await api.get(`/admin/users/${userId}/tasks`).then(extractData<UserTasksResponse>);
    return res.tasks;
  },

  promoteUser: (userId: string) =>
    api.put(`/admin/users/${userId}/promote`).then(extractData<User>),
};

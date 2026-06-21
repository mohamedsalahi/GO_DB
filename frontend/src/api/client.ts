import axios from 'axios';
import type {
  APIResponse,
  AuthResponse,
  User,
  Task,
  CreateTaskRequest,
  UpdateTaskRequest,
} from '../types';

export const api = axios.create({
  baseURL: '/api/v1',
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (res) => res,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      if (window.location.pathname !== '/login') {
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

export function extractData<T>(response: { data: APIResponse<T> }): T {
  if (!response.data.success) {
    throw new Error(response.data.error || 'An error occurred');
  }
  return response.data.data as T;
}

export const authAPI = {
  register: (name: string, email: string, password: string) =>
    api.post<APIResponse<AuthResponse>>('/auth/register', { name, email, password })
      .then(extractData),

  login: (email: string, password: string) =>
    api.post<APIResponse<AuthResponse>>('/auth/login', { email, password })
      .then(extractData),

  getProfile: () =>
    api.get<APIResponse<User>>('/users/me').then(extractData),
};

export const tasksAPI = {
  list: () =>
    api.get<APIResponse<Task[]>>('/tasks').then(extractData),

  get: (id: string) =>
    api.get<APIResponse<Task>>(`/tasks/${id}`).then(extractData),

  create: (data: CreateTaskRequest) =>
    api.post<APIResponse<Task>>('/tasks', data).then(extractData),

  update: (id: string, data: UpdateTaskRequest) =>
    api.put<APIResponse<Task>>(`/tasks/${id}`, data).then(extractData),

  delete: (id: string) =>
    api.delete(`/tasks/${id}`),
};

import { apiClient } from './client';

interface LoginRequest {
  username: string;
  password: string;
}

interface DeliveryOption {
  display_value: string;
  hashed_value: string;
}

interface TwoFactorMethod {
  type: string;
  delivery_options: DeliveryOption[];
}

interface LoginResponse {
  id?: string;
  email?: string;
  username?: string;
  name?: string | null;
  created_at?: string;
  last_modified_at?: string;
  deleted_at?: string | null;
  created_by?: string | null;
  roles?: Array<{
    id?: string;
    name?: string;
  }> | null;
  // 2FA fields
  status?: string;
  message?: string;
  temp_token?: string;
  two_factor_methods?: TwoFactorMethod[];
}

interface CreateUserRequest {
  email: string;
  username: string;
  name?: string | null;
  role_ids?: string[];
}

interface UpdateUserRequest {
  name?: string | null;
  username?: string;
  password?: string;
  role_ids?: string[];
}

interface FindUsernameRequest {
  email: string;
}

interface User {
  id?: string;
  email?: string;
  username?: string;
  name?: string | null;
  created_at?: string;
  last_modified_at?: string;
  deleted_at?: string | null;
  created_by?: string | null;
  roles?: Array<{
    id?: string;
    name?: string;
  }> | null;
}

export const userApi = {
  login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    const response = await apiClient.post('/auth/login', credentials, { skipAuth: true });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText || 'Login failed');
    }

    return response.json();
  },

  listUsers: async (): Promise<User[]> => {
    const response = await apiClient.get('/idm/users');

    if (!response.ok) {
      throw new Error('Failed to fetch users');
    }

    return response.json();
  },

  getUser: async (id: string): Promise<User> => {
    const response = await apiClient.get(`/idm/users/${id}`);

    if (!response.ok) {
      throw new Error('Failed to fetch user');
    }

    return response.json();
  },

  createUser: async (user: CreateUserRequest): Promise<User> => {
    const response = await apiClient.post('/idm/users', user);

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      if (errorData && errorData.message) {
        throw new Error(errorData.message);
      }
      throw new Error('Failed to create user');
    }

    return response.json();
  },

  updateUser: async (id: string, user: UpdateUserRequest): Promise<User> => {
    const response = await apiClient.put(`/idm/users/${id}`, user);

    if (!response.ok) {
      throw new Error('Failed to update user');
    }

    return response.json();
  },

  deleteUser: async (id: string): Promise<void> => {
    const response = await apiClient.delete(`/idm/users/${id}`);

    if (!response.ok) {
      throw new Error('Failed to delete user');
    }
  },

  findUsername: async (email: string): Promise<void> => {
    const response = await apiClient.post('/auth/find-username', { email }, { skipAuth: true });

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new Error(errorData?.message || 'Failed to find username');
    }
  }
};

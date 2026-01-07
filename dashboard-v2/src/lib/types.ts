export interface User {
  id: string;
  name: string;
  email: string;
  provider: string;
  providerId: string;
  avatar: string;
  labels: string[];
  createdAt: string;
  banned: boolean;
}

export interface CreateUserRequest {
  id?: string;
  name: string;
  email: string;
  password: string;
  provider: string;
  providerId: string;
  avatar: string;
  labels: string[];
}

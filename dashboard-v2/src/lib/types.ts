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

export interface Label {
  id: string;
  name: string;
  color: string;
  createdAt: string;
}

export interface Provider {
  ProviderCode: string;
  ClientID: string;
  ClientSecret: string;
  CallbackURL: string;
  IsEnabled: number;
}

export interface Session {
  id: string;
  userId: string;
  userName: string;
  userEmail: string;
  ipAddress: string;
  userAgent: string;
  createdAt: number;
  expiresAt: number;
  isActive: boolean;
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
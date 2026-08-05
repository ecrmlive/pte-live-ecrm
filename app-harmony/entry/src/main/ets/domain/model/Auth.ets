export interface AuthUser {
  id: string;
  nickname: string;
  avatarUrl?: string;
}

export interface AuthSession {
  accessToken: string;
  refreshToken: string;
  user: AuthUser;
}

export interface ApiEnvelope<T> {
  data: T;
  message?: string;
}

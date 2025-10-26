export interface User {
  Type: string;
  Name: string;
  Email: string;
  ID: number;
}

export interface LoginRequest {
  Type: string;
  Email: string;
  Password: string;
}
export interface LoginResponse {
  Message: string;
  User: User;
}

export interface SignupRequest {
  name: string;
  email: string;
  password: string;
}
export interface SignupResponse {
  Message: string;
  User: User;
}

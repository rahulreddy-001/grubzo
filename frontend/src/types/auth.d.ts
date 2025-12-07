import type { Location } from "./common";

export interface User {
  Type: string;
  Name: string;
  Email: string;
  Roles: string[] | null;
  Permisssions: string[] | null;
  Location: Location;
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

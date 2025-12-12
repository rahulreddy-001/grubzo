import { z } from "zod";
export interface ErrorResponse {
  Error: string;
}

export interface FileInfo {
  ID: string;
  Mime: string;
  Name: string;
  OwnerID: number;
  OwnerType: string;
  Size: number;
  Type: string;
  URL: string;
}

export const LocationSchema = z.object({
  ID: z.number().default(-1),
  Code: z.string().trim().min(3, "Code is required."),
  Address: z.string().trim().min(3, "Address is required."),
  City: z.string().trim().min(3, "City is required."),
  State: z.string().trim().min(3, "State is required."),
  Country: z.string().trim().min(3, "Country is required."),
  ZipCode: z.string().trim().min(6, "Enter valid zipcode."),
  IsPrimary: z.boolean().default(false),
});
export type Location = z.infer<typeof LocationSchema>;

export interface FetchLocationResponse {
  Message: string;
  Locations: Location[];
}
export interface LocationResponse {
  Message: string;
}

export const EmployeeSchema = z.object({
  ID: z.number().default(-1),
  Name: z.string().trim().min(1, "Employee name is required."),
  Email: z.email("Email is required."),
  Roles: z.array(z.string()).min(1, "Roles are required."),
  LocationID: z.number().default(-1),
});
export type Employee = z.infer<typeof EmployeeSchema>;

export interface FetchEmployeeResponse {
  Message: string;
  Users: Employee[];
}

export interface EmployeeResponse {
  Message: string;
}

export interface RBACResponse {
  Message: string;
  Data: {
    Roles: string[];
    Permissions: string[];
    Grid: Record<string, string[]> | null;
  };
}

export interface UpdateRole {
  Name: string;
  Permissions: string[];
  Action: number;
}

export const CreateRolePayloadSchema = z.object({
  Name: z.string().trim().min(1, "Role name is required."),
  Permissions: z.array(z.string()).min(1, "Select at least one permission."),
});
export type CreateRolePayload = z.infer<typeof CreateRolePayloadSchema>;

export interface CreateRoleResponse {
  Message: string;
}

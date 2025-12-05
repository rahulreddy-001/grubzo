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

  export interface Location {
    ID: number;
    Code: string;
    Address: string;
    City: string;
    State: string;
    Country: string;
    ZipCode: string;
    IsPrimary: boolean;
  }

  export interface FetchLocationResponse {
    Message: string;
    Locations: Location[];
  }
  export interface LocationResponse {
    Message: string;
  }

  export interface Employee {
    ID: number;
    Email: string;
    Name: string;
    Roles: string[];
    LocationID: number;
  }

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

  export interface CreateRolePayload {
    Name: string;
    Permissions: string[];
  }
  export interface CreateRoleResponse {
    Message: string;
  }
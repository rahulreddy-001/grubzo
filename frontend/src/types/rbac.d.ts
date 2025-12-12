export const PERMISSIONS = {
  DASHBOARD: "dashboard",
  ORDERS: "orders",
  ITEMS: "items",
  EMPLOYEE: "employee",
  LOCATION: "location",
  RBAC: "rbac",
} as const;

export type Permission = (typeof PERMISSIONS)[keyof typeof PERMISSIONS];

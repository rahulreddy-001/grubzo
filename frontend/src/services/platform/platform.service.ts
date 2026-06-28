import axios from "axios";

export type PlatformTenant = {
  ID: number;
  Name: string;
  Code: string;
  SubDomain: string;
};

export type PlatformTenantInput = {
  Tenant: {
    Name: string;
    Code: string;
    SubDomain: string;
  };
  Location: {
    Code: string;
    Address: string;
    City: string;
    State: string;
    Country: string;
    ZipCode: string;
  };
  Admin: {
    Email: string;
    Password: string;
    Name: string;
  };
};

export type PlatformTenantUpdate = {
  Name?: string;
};

const BASE = "/platform/v1";

const PlatformService = {
  async login(Username: string, Password: string): Promise<void> {
    await axios.post(`${BASE}/login`, { Username, Password });
  },

  async logout(): Promise<void> {
    await axios.post(`${BASE}/logout`);
  },

  async me(): Promise<{ User: string }> {
    const response = await axios.get<{ User: string }>(`${BASE}/me`);
    return response.data;
  },

  async listTenants(): Promise<PlatformTenant[]> {
    const response = await axios.get<{ Tenants: PlatformTenant[] }>(
      `${BASE}/tenants`
    );
    return response.data.Tenants ?? [];
  },

  async createTenant(input: PlatformTenantInput): Promise<PlatformTenant> {
    const response = await axios.post<{ Tenant: PlatformTenant }>(
      `${BASE}/tenants`,
      input
    );
    return response.data.Tenant;
  },

  async updateTenant(
    tenantID: number,
    input: PlatformTenantUpdate
  ): Promise<PlatformTenant> {
    const response = await axios.put<{ Tenant: PlatformTenant }>(
      `${BASE}/tenants/${tenantID}`,
      input
    );
    return response.data.Tenant;
  },
};

export default PlatformService;

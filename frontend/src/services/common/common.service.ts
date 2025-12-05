import type {
  CreateRolePayload,
  CreateRoleResponse,
  Employee,
  EmployeeResponse,
  FileInfo,
  LocationResponse,
  RBACResponse,
  Location,
} from "../../types/common";
import store from "../store";

import {
  fetchRBACInfo,
  createRole,
  syncRBACInfo,
  type SyncRBACInfoPayload,
} from "./rbac.slice";

import {
  fetchLocations,
  createLocation,
  updateLocation,
  setUserLocation,
} from "./common.slice";

import {
  fetchEmployees,
  createEmployee,
  updateEmployee,
} from "./employee.slice";

const CommonService = {
  async fetchRBACInfo(): Promise<void> {
    await store.dispatch(fetchRBACInfo()).unwrap();
  },

  async syncRBACInfo(body: SyncRBACInfoPayload): Promise<RBACResponse> {
    return store.dispatch(syncRBACInfo(body)).unwrap();
  },

  async createRole(body: CreateRolePayload): Promise<CreateRoleResponse> {
    return store.dispatch(createRole(body)).unwrap();
  },

  async setUserLocation(LocationID: number) {
    await setUserLocation(LocationID);
    location.reload();
  },

  async fetchLocations(): Promise<void> {
    await store.dispatch(fetchLocations()).unwrap();
  },

  async createLocation(body: Location): Promise<LocationResponse> {
    return store.dispatch(createLocation(body)).unwrap();
  },

  async updateLocation(body: Location): Promise<LocationResponse> {
    return store.dispatch(updateLocation(body)).unwrap();
  },

  async fetchEmployees(): Promise<void> {
    await store.dispatch(fetchEmployees()).unwrap();
  },

  async createEmployee(body: Employee): Promise<EmployeeResponse> {
    return store.dispatch(createEmployee(body)).unwrap();
  },

  async updateEmployee(body: Employee): Promise<EmployeeResponse> {
    return store.dispatch(updateEmployee(body)).unwrap();
  },

  async uploadWithProgress(
    formData: FormData,
    onProgress: (p: number) => void
  ): Promise<FileInfo[]> {
    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest();

      xhr.open("POST", "/api/v1/files/upload");

      xhr.upload.onprogress = (event) => {
        if (event.lengthComputable) {
          const percent = Math.round((event.loaded / event.total) * 100);
          onProgress(percent);
        }
      };

      xhr.onload = () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          try {
            const json = JSON.parse(xhr.responseText);
            resolve(json.filesMeta as FileInfo[]);
          } catch {
            reject(new Error("Invalid JSON response"));
          }
        } else {
          reject(new Error("File upload failed"));
        }
      };

      xhr.onerror = () => reject(new Error("Network error"));

      xhr.send(formData);
    });
  },
};

export default CommonService;

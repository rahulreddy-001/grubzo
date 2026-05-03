import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import axios from "axios";
import type {
  CreateRolePayload,
  CreateRoleResponse,
  ErrorResponse,
  RBACResponse,
  UpdateRole,
} from "../../types/common";
import { handleApiError } from "../api";

const FETCH_ROLES_PERMS_GRID = `/api/v1/rbac/roles_perms_grid`;
const CREATE_ROLE = `/api/v1/rbac/add_role`;
const UPDATE_ROLE_PERMS = `/api/v1/rbac/update_role_perms`;

export const fetchRBACInfo = createAsyncThunk<
  RBACResponse,
  void,
  { rejectValue: ErrorResponse }
>("rbac/fetchInfo", async (_, { rejectWithValue }) => {
  try {
    const response = await axios.get<RBACResponse>(FETCH_ROLES_PERMS_GRID);
    return response.data;
  } catch (error) {
    return handleApiError<ErrorResponse>(error, rejectWithValue);
  }
});

export interface SyncRBACInfoPayload {
  Data: UpdateRole[];
}

export const syncRBACInfo = createAsyncThunk<
  RBACResponse,
  SyncRBACInfoPayload,
  { rejectValue: ErrorResponse }
>("rbac/syncRBACInfo", async (body, { rejectWithValue }) => {
  try {
    const response = await axios.put<RBACResponse>(UPDATE_ROLE_PERMS, body);
    return response.data;
  } catch (error) {
    return handleApiError<ErrorResponse>(error, rejectWithValue);
  }
});

export const createRole = createAsyncThunk<
  CreateRoleResponse,
  CreateRolePayload,
  { rejectValue: ErrorResponse }
>("rbac/createRole", async (body, { rejectWithValue }) => {
  try {
    const response = await axios.post<CreateRoleResponse>(CREATE_ROLE, body);
    return response.data;
  } catch (error) {
    return handleApiError<ErrorResponse>(error, rejectWithValue);
  }
});

export interface RABCState {
  Roles: string[];
  Permissions: string[];
  Grid: Record<string, string[]> | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: RABCState = {
  Roles: [],
  Permissions: [],
  Grid: null,
  isLoading: false,
  error: null,
};

const rbacSlice = createSlice({
  name: "rbac",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchRBACInfo.pending, (state) => {
      state.isLoading = true;
      state.error = "";
    });
    builder.addCase(fetchRBACInfo.fulfilled, (state, action) => {
      state.isLoading = false;
      state.error = "";
      state.Grid = action.payload.Data.Grid;
      state.Permissions = action.payload.Data.Permissions;
      state.Roles = action.payload.Data.Roles;
    });
    builder.addCase(fetchRBACInfo.rejected, (state, action) => {
      state.isLoading = false;
      state.error = action.payload?.Error || null;
    });
  },
});

export default rbacSlice.reducer;

import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import axios from "axios";
import { handleApiError } from "../api";
import type {
  Employee,
  EmployeeResponse,
  FetchEmployeeResponse,
  ErrorResponse,
} from "../../types/common";

const EMPLOYEE_FETCH_ALL = `/api/v1/employee/all`;
const EMPLOYEE_CREATE = `/api/v1/employee/create`;
const EMPLOYEE_UPDATE = `/api/v1/employee/update`;

export const fetchEmployees = createAsyncThunk<
  FetchEmployeeResponse,
  void,
  { rejectValue: ErrorResponse }
>("employee/fetchEmployees", async (_, { rejectWithValue }) => {
  try {
    const response = await axios.get<FetchEmployeeResponse>(EMPLOYEE_FETCH_ALL);
    return response.data;
  } catch (error) {
    return handleApiError<ErrorResponse>(error, rejectWithValue);
  }
});

export const createEmployee = createAsyncThunk<
  EmployeeResponse,
  Employee,
  { rejectValue: ErrorResponse }
>("employee/createEmployee", async (body, { rejectWithValue }) => {
  try {
    const res = await axios.post<Employee>(EMPLOYEE_CREATE, body);
    return res.data;
  } catch (error) {
    return handleApiError<ErrorResponse>(error, rejectWithValue);
  }
});

export const updateEmployee = createAsyncThunk<
  EmployeeResponse,
  Employee,
  { rejectValue: ErrorResponse }
>("employee/updateEmployee", async (body, { rejectWithValue }) => {
  try {
    const res = await axios.put<Employee>(EMPLOYEE_UPDATE, body);
    return res.data;
  } catch (error) {
    return handleApiError<ErrorResponse>(error, rejectWithValue);
  }
});

export interface employeeState {
  Employees: Employee[] | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: employeeState = {
  Employees: [],
  isLoading: false,
  error: null,
};

const employeeSlice = createSlice({
  name: "employee",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchEmployees.pending, (state) => {
      state.isLoading = true;
      state.error = "";
    });
    builder.addCase(fetchEmployees.fulfilled, (state, action) => {
      state.isLoading = false;
      state.error = "";
      state.Employees = action.payload.Users;
    });
    builder.addCase(fetchEmployees.rejected, (state, action) => {
      state.isLoading = false;
      state.error = action.payload?.Error || null;
    });
  },
});

export default employeeSlice.reducer;

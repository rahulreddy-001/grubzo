import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import type { PayloadAction } from "@reduxjs/toolkit";
import axios from "axios";

import type {
  User,
  SignupRequest,
  SignupResponse,
  LoginRequest,
  LoginResponse,
} from "../../types/auth";
import type { ErrorResponse } from "../../types/common";
import { handleApiError } from "../api";

const LOGIN_URL = `/auth/v1/login`;
const LOGOUT_URL = `/auth/v1/logout`;
const FETCH_USER_URL = `/auth/v1/me`;
const SIGNUP_URL = `/api/v1/user/signup`;

export interface AuthState {
  user: User | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: AuthState = {
  user: null,
  isLoading: false,
  error: null,
};

export const loginUser = createAsyncThunk<
  LoginResponse,
  LoginRequest,
  { rejectValue: ErrorResponse }
>("auth/loginUser", async (body, { rejectWithValue }) => {
  try {
    const response = await axios.post<LoginResponse>(LOGIN_URL, body);
    return response.data;
  } catch (error) {
    return handleApiError<ErrorResponse>(error, rejectWithValue);
  }
});

export const signupUser = createAsyncThunk<
  SignupResponse,
  SignupRequest,
  { rejectValue: ErrorResponse }
>("auth/signupUser", async (body, { rejectWithValue }) => {
  try {
    const response = await axios.post<SignupResponse>(SIGNUP_URL, body);
    return response.data;
  } catch (error) {
    return handleApiError<ErrorResponse>(error, rejectWithValue);
  }
});

export const fetchUser = createAsyncThunk<
  User,
  void,
  { rejectValue: ErrorResponse }
>("auth/fetchUser", async (_, { rejectWithValue }) => {
  try {
    const response = await axios.get<User>(FETCH_USER_URL);
    return response.data;
  } catch (error) {
    return handleApiError<ErrorResponse>(error, rejectWithValue);
  }
});

export const logoutUser = createAsyncThunk<
  void,
  void,
  { rejectValue: ErrorResponse }
>("auth/logoutUser", async (_, { rejectWithValue }) => {
  try {
    await axios.post(LOGOUT_URL);
  } catch (error) {
    return handleApiError<ErrorResponse>(error, rejectWithValue);
  }
});

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
      state.isLoading = false;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(loginUser.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(loginUser.fulfilled, (state, action) => {
        state.user = action.payload.User;
        state.isLoading = false;
        state.error = null;
      })
      .addCase(loginUser.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload?.Error || null;
      })
      .addCase(signupUser.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(signupUser.fulfilled, (state, action) => {
        state.user = action.payload.User;
        state.isLoading = false;
        state.error = null;
      })
      .addCase(signupUser.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload?.Error || null;
      })
      .addCase(fetchUser.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchUser.fulfilled, (state, action) => {
        state.user = action.payload;
        state.isLoading = false;
      })
      .addCase(fetchUser.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload?.Error || null;
      })
      .addCase(logoutUser.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(logoutUser.fulfilled, (state) => {
        state.user = null;
        state.isLoading = false;
      })
      .addCase(logoutUser.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload?.Error || null;
      });
  },
});

export const { setLoading, setError } = authSlice.actions;
export default authSlice.reducer;

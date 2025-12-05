import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import axios from "axios";
import type {
  ErrorResponse,
  Location,
  LocationResponse,
  FetchLocationResponse,
} from "../../types/common";
import { handleApiError } from "../api";

const LOCATION_FETCH_ALL = `/api/v1/location/all`;
const LOCATION_CREATE = `/api/v1/location/create`;
const LOCATION_UPDATE = `/api/v1/location/update`;
const SET_USER_LOCATION = `/auth/v1/me/location`;

export const setUserLocation = async function (LocationID: number) {
  await axios.put(SET_USER_LOCATION, {
    LocationID,
  });
};

export const fetchLocations = createAsyncThunk<
  FetchLocationResponse,
  void,
  { rejectValue: ErrorResponse }
>("common/fetchLocations", async (_, { rejectWithValue }) => {
  try {
    const response = await axios.get<FetchLocationResponse>(LOCATION_FETCH_ALL);
    return response.data;
  } catch (error) {
    return handleApiError<ErrorResponse>(error, rejectWithValue);
  }
});

export const createLocation = createAsyncThunk<
  LocationResponse,
  Location,
  { rejectValue: ErrorResponse }
>("common/createLocation", async (body, { rejectWithValue }) => {
  try {
    const res = await axios.post<Location>(LOCATION_CREATE, body);
    return res.data;
  } catch (error) {
    return handleApiError<ErrorResponse>(error, rejectWithValue);
  }
});

export const updateLocation = createAsyncThunk<
  LocationResponse,
  Location,
  { rejectValue: ErrorResponse }
>("common/updateLocation", async (body, { rejectWithValue }) => {
  try {
    const res = await axios.put<Location>(LOCATION_UPDATE, body);
    return res.data;
  } catch (error) {
    return handleApiError<ErrorResponse>(error, rejectWithValue);
  }
});

export interface CommonState {
  Locations: Location[] | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: CommonState = {
  Locations: [],
  isLoading: false,
  error: null,
};

const commonSlice = createSlice({
  name: "common",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchLocations.pending, (state) => {
      state.isLoading = true;
      state.error = "";
    });
    builder.addCase(fetchLocations.fulfilled, (state, action) => {
      state.isLoading = false;
      state.error = "";
      state.Locations = action.payload.Locations;
    });
    builder.addCase(fetchLocations.rejected, (state, action) => {
      state.isLoading = false;
      state.error = action.payload?.Error || null;
    });
  },
});

export default commonSlice.reducer;

import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import axios from "axios";
import { handleApiError } from "../api";
import type { ErrorResponse } from "../../types/common";
import type { Item, ModifyItemPayload } from "../../types/item";

export interface ItemState {
  items: Item[];
  selectedItem: Item | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: ItemState = {
  items: [],
  selectedItem: null,
  isLoading: false,
  error: null,
};

const BASE = `/api/v1/item`;
const ITEM_GET_ALL = `${BASE}/all`;
const ITEM_CREATE = `${BASE}/create`;
const ITEM_UPDATE = `${BASE}/update`;

export const fetchAllItems = createAsyncThunk<
  Item[],
  void,
  { rejectValue: ErrorResponse }
>("items/fetchAll", async (_, { rejectWithValue }) => {
  try {
    const res = await axios.get(ITEM_GET_ALL);
    return res.data.Items as Item[];
  } catch (err) {
    return handleApiError<ErrorResponse>(err, rejectWithValue);
  }
});

export const createItem = createAsyncThunk<
  Item,
  ModifyItemPayload,
  { rejectValue: ErrorResponse }
>("items/create", async (body, { rejectWithValue }) => {
  try {
    const res = await axios.post(ITEM_CREATE, body);
    return res.data as Item;
  } catch (err) {
    return handleApiError<ErrorResponse>(err, rejectWithValue);
  }
});

export const updateItem = createAsyncThunk<
  Item,
  ModifyItemPayload,
  { rejectValue: ErrorResponse }
>("items/update", async (body, { rejectWithValue }) => {
  try {
    const res = await axios.put(ITEM_UPDATE, body);
    return res.data as Item;
  } catch (err) {
    return handleApiError<ErrorResponse>(err, rejectWithValue);
  }
});

const itemsSlice = createSlice({
  name: "items",
  initialState,
  reducers: {
    clearSelectedItem: (state) => {
      state.selectedItem = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchAllItems.pending, (s) => {
        s.isLoading = true;
        s.error = null;
      })
      .addCase(fetchAllItems.fulfilled, (s, a) => {
        s.items = a.payload;
        s.isLoading = false;
      })
      .addCase(fetchAllItems.rejected, (s, a) => {
        s.isLoading = false;
        s.error = a.payload?.Error || "Failed to load items";
      })
      .addCase(createItem.pending, (s) => {
        s.isLoading = true;
      })
      .addCase(createItem.fulfilled, (s, a) => {
        s.items.push(a.payload);
        s.isLoading = false;
      })
      .addCase(updateItem.pending, (s) => {
        s.isLoading = true;
      })
      .addCase(updateItem.fulfilled, (s, a) => {
        const idx = s.items.findIndex((i) => i.ID === a.payload.ID);
        if (idx !== -1) s.items[idx] = a.payload;
        s.isLoading = false;
      });
  },
});

export const { clearSelectedItem } = itemsSlice.actions;
export default itemsSlice.reducer;

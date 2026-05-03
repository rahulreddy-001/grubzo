import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import axios from "axios";
import { handleApiError } from "../api";
import type { ErrorResponse } from "../../types/common";

import {
  type CartResponse,
  type Item as CartItem,
  type UpdateItemQuantityPayload,
  type Bill,
  type OrderSubmitResponse,
  type OrderSubmitRequestPayload,
} from "../../types/cart.d";

import type { Item, FetchUserItemsResponse } from "../../types/item";

export interface ItemState {
  items: Item[];
  cart: CartItem[];
  bill: Bill | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: ItemState = {
  items: [],
  cart: [],
  bill: null,
  isLoading: false,
  error: null,
};

const BASE = `/api/v1/cart`;
const PUT_ITEM_QUANTITY = `${BASE}/item_quantity`;
const ITEM_GET_USER = `/api/v1/item/user`;

export const fetchUserItems = createAsyncThunk<
  Item[],
  void,
  { rejectValue: ErrorResponse }
>("user/fetchItems", async (_, { rejectWithValue }) => {
  try {
    const res = await axios.get<FetchUserItemsResponse>(ITEM_GET_USER);
    return res.data.Items as Item[];
  } catch (err) {
    return handleApiError<ErrorResponse>(err, rejectWithValue);
  }
});

export const fetchCart = createAsyncThunk<
  CartResponse,
  void,
  { rejectValue: ErrorResponse }
>("cart/get", async (_, { rejectWithValue }) => {
  try {
    const res = await axios.get<CartResponse>(BASE);
    return res.data as CartResponse;
  } catch (err) {
    return handleApiError<ErrorResponse>(err, rejectWithValue);
  }
});

export const updateCart = createAsyncThunk<
  CartResponse,
  UpdateItemQuantityPayload,
  { rejectValue: ErrorResponse }
>("cart/update", async (body, { rejectWithValue }) => {
  try {
    const res = await axios.put<CartResponse>(PUT_ITEM_QUANTITY, body);
    return res.data as CartResponse;
  } catch (err) {
    return handleApiError<ErrorResponse>(err, rejectWithValue);
  }
});

export const clearCart = createAsyncThunk<
  CartResponse,
  void,
  { rejectValue: ErrorResponse }
>("cart/delete", async (_, { rejectWithValue }) => {
  try {
    const res = await axios.delete<CartResponse>(BASE);
    return res.data as CartResponse;
  } catch (err) {
    return handleApiError<ErrorResponse>(err, rejectWithValue);
  }
});

export const placeOrder = createAsyncThunk<
  OrderSubmitResponse,
  OrderSubmitRequestPayload,
  { rejectValue: ErrorResponse }
>("cart/placeOrder", async (payload, { rejectWithValue }) => {
  try {
    const res = await axios.post<OrderSubmitResponse>(
      "/api/v1/order/create",
      payload
    );
    return res.data;
  } catch (err) {
    return handleApiError<ErrorResponse>(err, rejectWithValue);
  }
});

const itemsSlice = createSlice({
  name: "items",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchUserItems.pending, (s) => {
        s.isLoading = true;
        s.error = null;
      })
      .addCase(fetchUserItems.fulfilled, (s, a) => {
        s.items = a.payload;
        s.isLoading = false;
      })
      .addCase(fetchUserItems.rejected, (s, a) => {
        s.isLoading = false;
        s.error = a.payload?.Error || "Failed to load items";
      })
      .addCase(fetchCart.pending, (s) => {
        s.isLoading = true;
        s.error = null;
      })
      .addCase(fetchCart.fulfilled, (s, a) => {
        s.cart = a.payload.Cart.Items;
        s.bill = a.payload.Bill;
        s.isLoading = false;
      })
      .addCase(fetchCart.rejected, (s, a) => {
        s.isLoading = false;
        s.error = a.payload?.Error || "Failed to load cart";
      })
      .addCase(updateCart.pending, (s) => {
        s.isLoading = true;
      })
      .addCase(updateCart.fulfilled, (s, a) => {
        s.cart = a.payload.Cart.Items;
        s.bill = a.payload.Bill;
        s.isLoading = false;
      })
      .addCase(updateCart.rejected, (s, a) => {
        s.isLoading = false;
        s.error = a.payload?.Error || "Failed to update cart";
      })
      .addCase(clearCart.pending, (s) => {
        s.isLoading = true;
      })
      .addCase(clearCart.fulfilled, (s) => {
        s.cart = [];
        s.bill = null;
        s.isLoading = false;
      })
      .addCase(clearCart.rejected, (s, a) => {
        s.isLoading = false;
        s.error = a.payload?.Error || "Failed to clear cart";
      })
      .addCase(placeOrder.pending, (s) => {
        s.isLoading = true;
      })
      .addCase(placeOrder.rejected, (s, a) => {
        s.isLoading = false;
        s.error = a.payload?.Error || "Failed to place order";
      })
      .addCase(placeOrder.fulfilled, (s) => {
        s.cart = [];
        s.bill = null;
        s.isLoading = false;
      });
  },
});

export default itemsSlice.reducer;

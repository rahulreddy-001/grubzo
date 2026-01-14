import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import axios from "axios";
import { handleApiError } from "../api";
import type { ErrorResponse } from "../../types/common";
import type {
  Order,
  OrdersResponse,
  UpdateOrderStatusPayload,
  UpdateOrderStatusResponse,
} from "../../types/order";

export interface OrderState {
  orders: Order[];
  isLoading: boolean;
  error: string | null;
}

const initialState: OrderState = {
  orders: [],
  isLoading: false,
  error: null,
};

const BASE = `/api/v1/order`;
const ORDER_LIST = `${BASE}/list`;
const ORDER_UPDATE_STATUS = `${BASE}/update_order_status`;

export const updateOrderStatus = createAsyncThunk<
  UpdateOrderStatusResponse,
  UpdateOrderStatusPayload,
  { rejectValue: ErrorResponse }
>("orders/updateStatus", async (payload, { rejectWithValue }) => {
  try {
    const res = await axios.put<UpdateOrderStatusResponse>(
      ORDER_UPDATE_STATUS,
      payload
    );

    return res.data;
  } catch (err) {
    return handleApiError<ErrorResponse>(err, rejectWithValue);
  }
});

export const fetchOrders = createAsyncThunk<
  Order[],
  void,
  { rejectValue: ErrorResponse }
>("orders/fetchAll", async (_, { rejectWithValue }) => {
  try {
    const res = await axios.get<OrdersResponse>(ORDER_LIST);
    return res.data.Orders;
  } catch (err) {
    return handleApiError<ErrorResponse>(err, rejectWithValue);
  }
});

const orderSlice = createSlice({
  name: "orders",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchOrders.pending, (s) => {
        s.isLoading = true;
        s.error = null;
      })
      .addCase(fetchOrders.fulfilled, (s, a) => {
        s.orders = a.payload;
        s.isLoading = false;
      })
      .addCase(fetchOrders.rejected, (s, a) => {
        s.isLoading = false;
        s.error = a.payload?.Error || "Failed to load orders";
      })
      .addCase(updateOrderStatus.pending, (s) => {
        s.isLoading = true;
      })
      .addCase(updateOrderStatus.fulfilled, (s) => {
        s.isLoading = false;
      })
      .addCase(updateOrderStatus.rejected, (s, a) => {
        s.isLoading = false;
        s.error = a.payload?.Error || "Failed to update order";
      });
  },
});

export default orderSlice.reducer;

import axios, { AxiosError } from "axios";
import type { ErrorResponse } from "../types/common";

export function handleApiError<T>(
  error: unknown,
  rejectWithValue: (value: T) => any
) {
  let errorResponse: ErrorResponse = { Error: "Something went wrong" };
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError<{ error?: string }>;
    if (axiosError.response?.data?.error) {
      errorResponse.Error = axiosError.response.data.error;
    }
  }
  return rejectWithValue(errorResponse as T);
}

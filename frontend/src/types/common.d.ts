export interface ErrorResponse {
  Error: string;
}

type ApiResult<T> =
  | { success: true; data: T }
  | { success: false; error: ErrorResponse };

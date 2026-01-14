export type OrderItem = {
  ItemID: number;
  Name: string;
  Qty: number;
  Price: number;
  Total: number;
};

export type OrderBill = {
  Subtotal: number;
  Tax: number;
  PlatformFee: number;
  Discount: number;
  TotalPayable: number;
};

export type Order = {
  ID: number;
  UserID: number;
  UserName: string;
  UserEmail: string;
  Status: string; // pending, preparing, ready, delivered, cancelled
  PaymentStatus: string; // pending, paid, refunded, voided
  PaymentMode: string; // wallet, pos

  Items: OrderItem[];
  Bill: OrderBill;

  CreatedAt: string;
};

export type OrdersResponse = {
  Message: string;
  Orders: Order[];
};

export interface UpdateOrderStatusPayload {
  OrderID: number;
  OrderStatus: string;
  PaymentStatus: string;
}

export interface UpdateOrderStatusResponse {
  Message: string;
}

import { z } from "zod";

const ItemSchema = z.object({
  Item: z.number(),
  Quantity: z.number(),
});

export type UpdateItemQuantityPayload = z.infer<typeof ItemSchema>;
export type Item = z.infer<typeof ItemSchema>;

export const BillSchema = z.object({
  Subtotal: z.number(),
  TaxP: z.number(),
  Tax: z.number(),
  PlatformFeeP: z.number(),
  PlatformFee: z.number(),
  DiscountP: z.number(),
  Discount: z.number(),
  TotalPayable: z.number(),
});
export type Bill = z.infer<typeof BillSchema>;

const CartSchema = z.object({
  Key: z.string(),
  Items: z.array(ItemSchema),
});

export type Cart = z.infer<typeof CartSchema>;

const CartResponseSchema = z.object({
  Message: z.string(),
  Cart: CartSchema,
  RemovedItems: z.array(ItemSchema),
  Bill: BillSchema,
});

export type CartResponse = z.infer<typeof CartResponseSchema>;

export interface OrderSubmitRequestPayload {
  payment_mode: "wallet" | "upi" | "cod" | "card";
}

export interface OrderSubmitResponse {
  Message: string;
  OrderID: string;
}

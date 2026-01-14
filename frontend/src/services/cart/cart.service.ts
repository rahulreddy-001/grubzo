import store from "../store";
import {
  fetchCart,
  clearCart,
  fetchUserItems,
  updateCart,
  placeOrder,
} from "./cart.slice";
import type { AppDispatch } from "../store";
import type { UpdateItemQuantityPayload } from "../../types/cart";

class CartService {
  private dispatch: AppDispatch;

  constructor() {
    this.dispatch = store.dispatch;
    this.dispatch(fetchUserItems());
    this.dispatch(fetchCart());
  }

  async fetchCart() {
    return await this.dispatch(fetchCart());
  }

  async clearCart() {
    return await this.dispatch(clearCart());
  }

  async updateCart(payload: UpdateItemQuantityPayload) {
    return await this.dispatch(updateCart(payload)).unwrap();
  }

  async fetchUserItems() {
    return await this.dispatch(fetchUserItems()).unwrap();
  }

  async submitOrder(paymentMode: "wallet" | "upi" | "cod" | "card") {
    return this.dispatch(placeOrder({ payment_mode: paymentMode })).unwrap();
  }
}

export default new CartService();

import store, { type AppDispatch } from "../store";
import type {
  UpdateOrderStatusPayload,
  UpdateOrderStatusResponse,
} from "../../types/order";
import { updateOrderStatus } from "./order.slice";

class OrderService {
  private dispatch: AppDispatch;

  constructor() {
    this.dispatch = store.dispatch;
  }

  updateStatus(
    payload: UpdateOrderStatusPayload
  ): Promise<UpdateOrderStatusResponse> {
    return this.dispatch(updateOrderStatus(payload)).unwrap();
  }
}

export default new OrderService();

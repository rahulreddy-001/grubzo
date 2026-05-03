import store from "../store";
import { fetchAllItems, createItem, updateItem } from "./item.slice";
import type { AppDispatch } from "../store";
import type { ModifyItemPayload } from "../../types/item";

class ItemsService {
  private dispatch: AppDispatch;

  constructor() {
    this.dispatch = store.dispatch;
  }

  async getAll() {
    return this.dispatch(fetchAllItems());
  }

  async create(data: ModifyItemPayload) {
    return this.dispatch(createItem(data));
  }

  async update(data: ModifyItemPayload) {
    return this.dispatch(updateItem(data));
  }
}

export default new ItemsService();

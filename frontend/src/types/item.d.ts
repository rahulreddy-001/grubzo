export type FoodCategoryType = "veg" | "nonveg" | "egg";
export type ItemStatusValue = "av" | "os" | "ac";

export const FoodCategoryOptions = [
  { value: "veg", label: "Veg", color: "green" },
  { value: "nonveg", label: "Non-Veg", color: "red" },
  { value: "egg", label: "Egg", color: "yellow" },
];

export const ItemStatusOptions = [
  { value: "av", label: "Available", color: "green" },
  { value: "os", label: "Out of Stock", color: "red" },
  { value: "ac", label: "Archived", color: "gray" },
];

export interface Item {
  ID: number;
  TenantID: number;
  LocationID: number;
  Name: string;
  Description: string;
  Price: number;
  Category: string;
  FoodType: FoodCategoryType;
  ItemStatus: ItemStatusValue;
  CreatedAt: string;
  UpdatedAt: string;
  Files: FileInfo[];
}

export interface ModifyItemPayload {
  ID: number | null;
  LocationID: number | null;
  Name: string;
  Description: string;
  Price: number;
  Category: string;
  FoodType: FoodCategoryType;
  ItemStatus: ItemStatusValue;
  FileIDs: string[];
}

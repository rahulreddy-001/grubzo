import type { X } from "lucide-react";
import { z } from "zod";

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

export const ItemSchema = z.object({
  ID: z.number(),
  TenantID: z.number(),
  LocationID: z.number(),
  Name: z.string(),
  Description: z.string(),
  Price: z.number(),
  Category: z.string(),
  FoodType: z.enum(["veg", "nonveg", "egg"]),
  ItemStatus: z.enum(["av", "os", "ac"]),
  CreatedAt: z.string(),
  UpdatedAt: z.string(),
  Files: z.array(z.any()),
});
export type Item = z.infer<typeof ItemSchema>;

export const ModifyItemPayloadSchema = z.object({
  ID: z.number().nullable(),
  LocationID: z.number().nullable(),
  Name: z.string().trim().min(1, "Item name is required."),
  Description: z.string().optional().default(""),
  Price: z
    .number({ invalid_type_error: "Enter a valid price." })
    .positive("Enter a valid price."),
  Category: z.string().trim().min(1, "Category is required."),
  FoodType: z.enum(["veg", "nonveg", "egg"]),
  ItemStatus: z.enum(["av", "os", "ac"]),
  FileIDs: z.array(z.string()),
});
export type ModifyItemPayload = z.infer<typeof ModifyItemPayloadSchema>;

import React, { useEffect, useState, useCallback } from "react";

import CModel from "../../../components/common/CModel";
import CInput from "../../../components/common/CInput";
import CButton from "../../../components/common/CButton";
import CUpload from "../../../components/common/CUpload";
import CSelect from "../../../components/common/CSelect";
import { Box, Flex, Text, TextArea } from "@radix-ui/themes";
import {
  FoodCategoryOptions,
  ItemStatusOptions,
  type Item,
  type ModifyItemPayload,
  type FoodCategoryType,
  type ItemStatusValue,
} from "../../../types/item.d";
import type { FileInfo } from "../../../types/common";

import { useSelector } from "react-redux";
import type { RootState } from "../../../services/store";

import CommonService from "../../../services/common/common.service";

const ModifyItem: React.FC<{
  item?: Item;
  onSave: (data: ModifyItemPayload) => void;
  onCancel: () => void;
}> = ({ item, onSave, onCancel }) => {
  const { Locations } = useSelector((s: RootState) => s.common);

  const [form, setForm] = useState<ModifyItemPayload>({
    ID: null,
    LocationID: null,
    Name: "",
    Description: "",
    Price: 0,
    Category: "",
    FoodType: "veg",
    ItemStatus: "av",
    FileIDs: [],
  });

  const [existingFiles, setExistingFiles] = useState<FileInfo[]>([]);
  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    if (item) {
      setForm({
        ID: item.ID,
        Name: item.Name || "",
        Description: item.Description || "",
        Price: item.Price || 0,
        Category: item.Category || "",
        ItemStatus: item.ItemStatus || "av",
        FoodType: (item as any).FoodType || "veg",
        LocationID: item.LocationID,
        FileIDs: item.Files?.map((f) => f.ID) ?? [],
      });

      setExistingFiles(Array.isArray(item.Files) ? item.Files : []);
    }

    if (Locations?.length === 0) {
      CommonService.fetchLocations();
    }
  }, [item, Locations]);

  const handleChange = <K extends keyof ModifyItemPayload>(
    key: K,
    value: ModifyItemPayload[K]
  ) => {
    setForm((prev) => ({ ...prev, [key]: value }));
  };

  const validate = useCallback(() => {
    const e: Record<string, string> = {};

    if (!form.Name?.trim()) e.Name = "Item name is required.";
    if (!form.Price || form.Price <= 0) e.Price = "Enter a valid price.";
    if (!form.Category?.trim()) e.Category = "Category is required.";

    return e;
  }, [form]);

  const handleSubmit = () => {
    const v = validate();
    setErrors(v);
    if (Object.keys(v).length > 0) return;
    console.log(form);
    onSave(form);
  };

  const isEditing = !!item;

  return (
    <CModel
      open
      size="md"
      onClose={onCancel}
      title={isEditing ? "Update Item" : "Add New Item"}
      actions={
        <Flex gap="3">
          <CButton label="Cancel" variant="soft" onClick={onCancel} />
          <CButton
            label={isEditing ? "Update" : "Add"}
            variant="solid"
            onClick={handleSubmit}
            disabled={!form.Name || !form.Price}
          />
        </Flex>
      }
    >
      <Flex direction="column" gap="4">
        <Box>
          <CInput
            label="Item Name"
            value={form.Name}
            placeholder="Enter item name"
            onChange={(val) => {
              handleChange("Name", val);
              setErrors((e) => ({ ...e, Name: "" }));
            }}
            error={errors.Name}
            fullWidth
          />
        </Box>

        <Box>
          <Text size="1" weight="medium">
            Description
          </Text>
          <TextArea
            value={form.Description || ""}
            onChange={(e) => handleChange("Description", e.target.value)}
            placeholder="Enter description"
            style={{
              width: "100%",
              padding: "8px",
              fontSize: "14px",
              borderRadius: "4px",
              border: "1px solid var(--gray-6)",
              background: "var(--color-panel-solid)",
              outline: "none",
              resize: "vertical",
            }}
          />
        </Box>

        <Box>
          <CInput
            label="Price (₹)"
            type="number"
            value={form.Price ?? ""}
            placeholder="Enter price"
            onChange={(val) => handleChange("Price", parseFloat(val) || 0)}
            error={errors.Price}
            fullWidth
          />
        </Box>

        <Box>
          <CInput
            label="Category"
            value={form.Category}
            placeholder="Category"
            onChange={(val) => {
              handleChange("Category", val);
              setErrors((e) => ({ ...e, Category: "" }));
            }}
            error={errors.Category}
            fullWidth
          />
        </Box>

        <Box>
          <CSelect
            label="Food Type"
            value={form.FoodType}
            options={FoodCategoryOptions}
            onChange={(v) => handleChange("FoodType", v as FoodCategoryType)}
            placeholder="Select food type"
          />
        </Box>

        <Box>
          <CSelect
            label="Status"
            value={form.ItemStatus}
            options={ItemStatusOptions}
            onChange={(v) => handleChange("ItemStatus", v as ItemStatusValue)}
            placeholder="Select status"
          />
        </Box>

        <Box>
          <Text size="1" weight="medium">
            Upload Item Images
          </Text>

          <CUpload
            existingFiles={existingFiles}
            multiple
            accept="image/*"
            maxFiles={10}
            maxSizeMB={5}
            onFilesChange={({ order }) => {
              const ids = order.map((f) => f.ID);
              handleChange("FileIDs", ids);
            }}
          />
        </Box>
      </Flex>
    </CModel>
  );
};

export default ModifyItem;

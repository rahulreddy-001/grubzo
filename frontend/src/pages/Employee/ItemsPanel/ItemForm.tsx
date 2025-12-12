import React from "react";
import { useFormik } from "formik";
import { zodToFormik } from "../../../utils/zodToFormik";

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
  ModifyItemPayloadSchema,
} from "../../../types/item.d";

const ItemForm: React.FC<{
  item: Item | null;
  onSave: (data: ModifyItemPayload) => void;
  onCancel: () => void;
}> = ({ item, onSave, onCancel }) => {
  const isUpdate = !!item;

  const form = useFormik<ModifyItemPayload>({
    initialValues: {
      ID: item?.ID ?? null,
      LocationID: item?.LocationID ?? null,
      Name: item?.Name ?? "",
      Description: item?.Description ?? "",
      Price: item?.Price ?? 0,
      Category: item?.Category ?? "",
      FoodType: item?.FoodType ?? "veg",
      ItemStatus: item?.ItemStatus ?? "av",
      FileIDs: item?.Files?.map((f) => f.ID) ?? [],
    },
    validate: zodToFormik(ModifyItemPayloadSchema),
    onSubmit: onSave,
  });

  return (
    <CModel
      open
      size="md"
      onClose={onCancel}
      title={isUpdate ? "Update Item" : "Add New Item"}
      actions={
        <Flex gap="3">
          <CButton label="Cancel" variant="soft" onClick={onCancel} />
          <CButton
            label={isUpdate ? "Update" : "Add"}
            variant="solid"
            onClick={form.submitForm}
            disabled={!form.values.Name || !form.values.Price}
          />
        </Flex>
      }
    >
      <form onSubmit={form.handleSubmit}>
        <Flex direction="column" gap="4">
          <Box>
            <CInput
              label="Item Name"
              placeholder="Enter item name"
              value={form.values.Name}
              onChange={(v) => form.setFieldValue("Name", v)}
              error={form.errors.Name}
              fullWidth
            />
          </Box>

          <Box>
            <Text size="1" weight="medium">
              Description
            </Text>
            <TextArea
              placeholder="Enter description"
              value={form.values.Description}
              onChange={(e) =>
                form.setFieldValue("Description", e.target.value)
              }
            />
          </Box>

          <Box>
            <CInput
              label="Price"
              placeholder="Enter price"
              value={form.values.Price}
              onChange={(val) =>
                form.setFieldValue("Price", parseFloat(val) || 0)
              }
              error={form.errors.Price}
              fullWidth
            />
          </Box>

          <Box>
            <CInput
              label="Category"
              placeholder="Category"
              value={form.values.Category}
              onChange={(val) => form.setFieldValue("Category", val)}
              error={form.errors.Category}
              fullWidth
            />
          </Box>

          <Box>
            <CSelect
              label="Food Type"
              placeholder="Select food type"
              value={form.values.FoodType}
              options={FoodCategoryOptions}
              onChange={(v) =>
                form.setFieldValue("FoodType", v as FoodCategoryType)
              }
            />
          </Box>

          <Box>
            <CSelect
              label="Status"
              value={form.values.ItemStatus}
              options={ItemStatusOptions}
              onChange={(v) =>
                form.setFieldValue("ItemStatus", v as ItemStatusValue)
              }
              placeholder="Select status"
            />
          </Box>

          <Box>
            <Text size="1" weight="medium">
              {"Upload Item Images"}
            </Text>
            <CUpload
              existingFiles={item?.Files ?? []}
              multiple
              accept="image/*"
              maxFiles={10}
              maxSizeMB={5}
              onFilesChange={({ order }) => {
                const ids = order.map((f) => f.ID);
                form.setFieldValue("FileIDs", ids);
              }}
            />
          </Box>
        </Flex>
      </form>
    </CModel>
  );
};

export default ItemForm;

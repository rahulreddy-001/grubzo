import React, { useEffect, useState } from "react";
import { useSelector } from "react-redux";
import type { RootState } from "../../../services/store";

import ItemsService from "../../../services/item/item.service";
import { apiUrl } from "../../../services/api";

import CTable from "../../../components/common/CTable";
import CButton from "../../../components/common/CButton";
import ItemForm from "./ItemForm";
import { Box, Flex, Text, Avatar, Badge, IconButton } from "@radix-ui/themes";
import { Plus, Edit, Hamburger } from "lucide-react";
import {
  FoodCategoryOptions,
  ItemStatusOptions,
  type ModifyItemPayload,
} from "../../../types/item.d";

const ItemsPanel: React.FC = () => {
  const { items, isLoading } = useSelector((s: RootState) => s.item);

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editItem, setEditItem] = useState<any>(null);
  const [search, setSearch] = useState("");

  useEffect(() => {
    ItemsService.getAll();
  }, []);

  const handleAdd = () => {
    setEditItem(null);
    setDrawerOpen(true);
  };

  const handleEdit = (item: any) => {
    setEditItem(item);
    setDrawerOpen(true);
  };

  const handleSave = async (data: ModifyItemPayload) => {
    if (editItem) {
      await ItemsService.update(data);
    } else {
      await ItemsService.create(data);
    }
    setDrawerOpen(false);
    await ItemsService.getAll();
  };

  const filteredItems = items.filter((item) => {
    if (!item) return false;
    const name = item.Name?.toLowerCase() ?? "";
    const desc = item.Description?.toLowerCase() ?? "";
    const s = search.toLowerCase();
    return name.includes(s) || desc.includes(s);
  });

  return (
    <Box>
      <CTable
        title="Menu Items"
        data={filteredItems}
        rowKey="ID"
        loading={isLoading}
        searchable
        onSearch={setSearch}
        searchPlaceholder="Search items..."
        onRefresh={() => ItemsService.getAll()}
        actions={
          <CButton
            label="Add Item"
            startIcon={<Plus size={16} />}
            onClick={handleAdd}
          />
        }
        columns={[
          {
            key: "Name",
            label: "Item",
            render: (item) => (
              <Flex gap="3" align="center" style={{ marginLeft: "10px" }}>
                <Avatar
                  size="4"
                  fallback={<Hamburger size={14} />}
                  src={apiUrl(item.Files?.[0]?.URL ?? "")}
                />

                <Box>
                  <Text weight="medium">{item.Name}</Text>
                  <Text as="div" size="1" color="gray">
                    {item.Description || ""}
                  </Text>
                </Box>
              </Flex>
            ),
          },

          {
            key: "Price",
            label: "Price",
            render: (item) => (
              <Text>₹{Number(item.Price ?? 0).toFixed(2)}</Text>
            ),
          },
          {
            key: "Category",
            label: "Category",
            render: (item) => <Badge size="2">{item.Category}</Badge>,
          },
          {
            key: "FoodType",
            label: "Food Type",
            render: (item) => {
              const status = FoodCategoryOptions.find(
                (opt) => opt.value === item.FoodType
              );
              return (
                <Badge size="2" variant="soft" color={status?.color as any}>
                  {status?.label ?? ""}
                </Badge>
              );
            },
          },
          {
            key: "Status",
            label: "Status",
            render: (item) => {
              const status = ItemStatusOptions.find(
                (opt) => opt.value === item.ItemStatus
              );
              return (
                <Badge size="2" variant="soft" color={status?.color as any}>
                  {status?.label ?? ""}
                </Badge>
              );
            },
          },

          {
            key: "Actions",
            label: "Actions",
            render: (item) => (
              <Flex gap="2" style={{ paddingLeft: "10px" }}>
                <IconButton
                  style={{
                    padding: "10px",
                    cursor: "pointer",
                  }}
                  variant="ghost"
                  radius="full"
                  onClick={() => handleEdit(item)}
                >
                  <Edit size={16} />
                </IconButton>
              </Flex>
            ),
          },
        ]}
      />

      {drawerOpen && (
        <ItemForm
          item={editItem}
          onSave={handleSave}
          onCancel={() => setDrawerOpen(false)}
        />
      )}
    </Box>
  );
};

export default ItemsPanel;

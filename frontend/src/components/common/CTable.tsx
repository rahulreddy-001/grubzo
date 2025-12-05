import React, { useEffect, useState } from "react";
import {
  Box,
  Flex,
  Text,
  Popover,
  Checkbox,
  TextField,
  Table,
  Card,
  Skeleton,
} from "@radix-ui/themes";
import { RefreshCw, Settings2, Search } from "lucide-react";
import CButton from "./CButton";
import "./style.css";

export interface Column<T> {
  key: keyof T | string;
  label: string;
  render?: (row: T) => React.ReactNode;
  width?: string | number;
  visible?: boolean;
  align?: "left" | "right" | "center";
}

export interface CTableProps<T> {
  title?: string;
  data: T[];
  columns: Column<T>[];
  onRefresh?: () => void;
  loading?: boolean;
  searchable?: boolean;
  searchPlaceholder?: string;
  onSearch?: (q: string) => void;
  actions?: React.ReactNode;
  rowKey: keyof T;
  emptyMessage?: string;
}

export default function CTable<T extends object>({
  title,
  data,
  columns,
  onRefresh,
  loading = false,
  searchable = false,
  searchPlaceholder = "Search...",
  onSearch,
  actions,
  rowKey,
  emptyMessage = "No data available",
}: CTableProps<T>) {
  const [search, setSearch] = useState("");
  const [visibleColumns, setVisibleColumns] = useState<Column<T>[]>(() =>
    columns.map((c) => ({ ...c, visible: c.visible !== false }))
  );

  useEffect(() => {
    setVisibleColumns(
      columns.map((c) => ({ ...c, visible: c.visible !== false }))
    );
  }, [columns]);

  const handleToggleColumn = (key: string | keyof T) => {
    const keyStr = String(key);
    setVisibleColumns((prev) =>
      prev.map((col) =>
        String(col.key) === keyStr
          ? { ...col, visible: !(col.visible !== false) }
          : col
      )
    );
  };

  const handleSearch = (value: string) => {
    setSearch(value);
    onSearch?.(value);
  };

  const visibleCount =
    visibleColumns.filter((c) => c.visible !== false).length || 1;

  return (
    <Box style={{ borderRadius: "0" }}>
      {(title || searchable || onRefresh || actions) && (
        <Flex justify="between" align="center" style={{ margin: "10px" }}>
          {title && <Text size="3">{title}</Text>}
          <Flex gap="3" align="center">
            {searchable && (
              <TextField.Root
                placeholder={searchPlaceholder}
                value={search}
                onChange={(e) => handleSearch(e.target.value)}
              >
                <TextField.Slot>
                  <Search height="16" width="16" />
                </TextField.Slot>
              </TextField.Root>
            )}

            {onRefresh && (
              <CButton
                label=""
                variant="surface"
                startIcon={<RefreshCw size={16} strokeWidth={2} />}
                onClick={() => onRefresh()}
              />
            )}

            <Popover.Root>
              <Popover.Trigger>
                <CButton
                  label=""
                  variant="surface"
                  startIcon={<Settings2 size={16} />}
                />
              </Popover.Trigger>
              <Popover.Content
                style={{
                  padding: "8px",
                  borderRadius: "3px",
                  fontSize: "16px",
                }}
              >
                <Flex direction="column" gap="1">
                  {visibleColumns.map((col) => {
                    const keyStr = String(col.key);
                    const isVisible = col.visible !== false;
                    return (
                      <Box
                        style={{
                          cursor: "pointer",
                          padding: "5px",
                          fontSize: "14px",
                        }}
                        className="c-hover-accent"
                      >
                        <Flex
                          key={keyStr}
                          align="center"
                          gap="1"
                          onClick={() => handleToggleColumn(keyStr)}
                        >
                          <Checkbox checked={isVisible} />
                          <Text style={{ paddingLeft: "5px" }}>
                            {col.label}
                          </Text>
                        </Flex>
                      </Box>
                    );
                  })}
                </Flex>
              </Popover.Content>
            </Popover.Root>

            {actions}
          </Flex>
        </Flex>
      )}

      <Card
        style={{
          padding: "0",
          margin: "0",
          ["--card-border-radius" as any]: "none",
        }}
      >
        {loading ? (
          <Table.Root size="2">
            <Table.Body>
              {[1, 2, 3, 4, 5].map((i) => (
                <Table.Row key={i}>
                  {visibleColumns
                    .filter((c) => c.visible !== false)
                    .map((col) => (
                      <Table.Cell key={String(col.key)}>
                        <Skeleton />
                      </Table.Cell>
                    ))}
                </Table.Row>
              ))}
            </Table.Body>
          </Table.Root>
        ) : (
          <Table.Root size="1" variant="surface">
            <Table.Header
              style={{
                background: "var(--accent-1)",
                color: "var(--accent-12)",
              }}
            >
              <Table.Row>
                {visibleColumns
                  .filter((c) => c.visible !== false)
                  .map((col) => (
                    <Table.ColumnHeaderCell
                      key={String(col.key)}
                      style={{ width: col.width, paddingLeft: "10px" }}
                    >
                      {col.label}
                    </Table.ColumnHeaderCell>
                  ))}
              </Table.Row>
            </Table.Header>

            <Table.Body>
              {data.length === 0 ? (
                <Table.Row>
                  <Table.Cell colSpan={visibleCount}>
                    <Flex
                      justify="center"
                      style={{
                        width: "100%",
                        height: "50px",
                        alignItems: "center",
                      }}
                    >
                      <Text>{emptyMessage}</Text>
                    </Flex>
                  </Table.Cell>
                </Table.Row>
              ) : (
                data.map((row) => {
                  const rowId = String((row as any)[rowKey]);
                  return (
                    <Table.Row key={rowId} className="c-hover-accent-1">
                      {visibleColumns
                        .filter((c) => c.visible !== false)
                        .map((col) => (
                          <Table.Cell
                            key={String(col.key)}
                            style={{
                              textAlign: col.align ?? "left",
                              verticalAlign: "middle",
                            }}
                          >
                            {col.render
                              ? col.render(row)
                              : col.key in row
                              ? String((row as any)[col.key as keyof T])
                              : ""}
                          </Table.Cell>
                        ))}
                    </Table.Row>
                  );
                })
              )}
            </Table.Body>
          </Table.Root>
        )}
      </Card>
    </Box>
  );
}

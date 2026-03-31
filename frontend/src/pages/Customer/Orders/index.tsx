import { useEffect, useState } from "react";
import { Box, Text, Badge, Flex, Tooltip } from "@radix-ui/themes";
import CTable from "../../../components/common/CTable";
import { type Column } from "../../../components/common/CTable";
import { Info } from "lucide-react";

type OrderItem = {
  ItemID: number;
  Name: string;
  Qty: number;
  Price: number;
  Total: number;
};

type Order = {
  ID: number;
  Status: string;
  PaymentStatus: string;
  PaymentMode: string;
  Items: OrderItem[];
  Bill: {
    Subtotal: number;
    Tax: number;
    PlatformFee: number;
    Discount: number;
    TotalPayable: number;
  };
  CreatedAt: string;
};

type OrdersResponse = {
  Message: string;
  Orders: Order[];
};

const formatMoney = (p: number) => `₹${(p / 100).toFixed(2)}`;

function Orders() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [filteredOrders, setFilteredOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  const loadOrders = async () => {
    try {
      setLoading(true);
      const res = await fetch("/api/v1/order/user_orders");
      const data: OrdersResponse = await res.json();
      const list = data.Orders || [];
      setOrders(list);
      setFilteredOrders(list);
    } catch (err) {
      console.error("Failed to load orders", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadOrders();
  }, []);

  const onSearch = (q: string) => {
    const query = q.trim().toLowerCase();

    if (!query) {
      setFilteredOrders(orders);
      return;
    }

    const res = orders.filter((o) => {
      if (String(o.ID).includes(query)) return true;
      if (o.Status.toLowerCase().includes(query)) return true;
      if (o.PaymentMode.toLowerCase().includes(query)) return true;
      if (o.PaymentStatus.toLowerCase().includes(query)) return true;

      return o.Items.some((i) => i.Name.toLowerCase().includes(query));
    });

    setFilteredOrders(res);
  };

  const renderItems = (row: Order) => {
    return (
      <Flex align="start" gap="2">
        <Box style={{ maxWidth: 260 }}>
          <Text size="2">
            {row.Items.map((i) => `${i.Name} × ${i.Qty}`).join(", ")}
          </Text>
        </Box>

        <Tooltip
          content={
            <Box
              p="3"
              style={{
                borderRadius: 12,
                minWidth: 240,
              }}
            >
              {row.Items.map((i) => (
                <Flex key={i.ItemID} justify="between" mb="2" gap="3">
                  <Box>
                    <Text size="2" weight="medium">
                      {i.Name}{" "}
                    </Text>
                    <Text size="1">
                      ({formatMoney(i.Price)} × {i.Qty})
                    </Text>
                  </Box>
                  <Text size="2" weight="bold">
                    {formatMoney(i.Total)}
                  </Text>
                </Flex>
              ))}
            </Box>
          }
        >
          <Info
            size={16}
            style={{
              cursor: "pointer",
              opacity: 0.7,
              marginTop: 2,
            }}
          />
        </Tooltip>
      </Flex>
    );
  };

  const renderBill = (row: Order) => {
    const b = row.Bill;

    return (
      <Flex align="center" gap="2" justify="end">
        <Text weight="bold">{formatMoney(b.TotalPayable)}</Text>

        <Tooltip
          content={
            <Box
              p="3"
              style={{
                borderRadius: 12,
                minWidth: 220,
              }}
            >
              <Flex justify="between">
                <Text size="2">Subtotal</Text>
                <Text size="2">{formatMoney(b.Subtotal)}</Text>
              </Flex>
              <Flex justify="between">
                <Text size="2">Tax</Text>
                <Text size="2">{formatMoney(b.Tax)}</Text>
              </Flex>
              <Flex justify="between">
                <Text size="2">Platform Fee</Text>
                <Text size="2">{formatMoney(b.PlatformFee)}</Text>
              </Flex>

              {b.Discount > 0 && (
                <Flex justify="between">
                  <Text size="2" color="green">
                    Discount
                  </Text>
                  <Text size="2" color="green">
                    −{formatMoney(b.Discount)}
                  </Text>
                </Flex>
              )}

              <Box my="2" style={{ borderTop: "1px dashed var(--gray-6)" }} />

              <Flex justify="between">
                <Text weight="bold">Payable</Text>
                <Text weight="bold">{formatMoney(b.TotalPayable)}</Text>
              </Flex>
            </Box>
          }
        >
          <Info
            size={16}
            style={{
              cursor: "pointer",
              opacity: 0.7,
            }}
          />
        </Tooltip>
      </Flex>
    );
  };

  const columns: Column<Order>[] = [
    { key: "ID", label: "Order", align: "right", width: 80 },
    {
      key: "CreatedAt",
      label: "Date",
      render: (r) => new Date(r.CreatedAt).toLocaleString(),
    },
    {
      key: "Items",
      label: "Items",
      render: renderItems,
    },
    {
      key: "Status",
      label: "Status",
      render: (r) => (
        <Badge color={r.Status === "delivered" ? "green" : "orange"}>
          {r.Status}
        </Badge>
      ),
    },
    {
      key: "PaymentStatus",
      label: "Payment",
      render: (r) => (
        <Badge
          color={
            r.PaymentStatus === "paid"
              ? "green"
              : r.PaymentStatus === "refunded"
              ? "red"
              : "gray"
          }
        >
          {r.PaymentStatus}
        </Badge>
      ),
    },
    { key: "PaymentMode", label: "Mode" },
    {
      key: "Bill",
      label: "Total",
      align: "right",
      render: renderBill,
    },
  ];

  return (
    <Box>
      <CTable
        title="Your Orders"
        data={filteredOrders}
        columns={columns}
        rowKey="ID"
        loading={loading}
        searchable
        searchPlaceholder="Search orders"
        onSearch={onSearch}
        onRefresh={loadOrders}
        emptyMessage="No orders found"
      />
    </Box>
  );
}

export default Orders;

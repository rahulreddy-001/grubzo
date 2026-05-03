import React, { useEffect, useState } from "react";
import { useSelector, useDispatch } from "react-redux";
import type { RootState } from "../../../services/store";

import { fetchOrders } from "../../../services/order/order.slice";
import type { Order } from "../../../types/order.d";

import CTable from "../../../components/common/CTable";
import { Box, Flex, Text, Badge, Tooltip, IconButton } from "@radix-ui/themes";
import { Info, Edit } from "lucide-react";
import OrderActionForm from "./OrderActionForm";

const formatMoney = (p: number) => `₹${(p / 100).toFixed(2)}`;

const OrdersPanel: React.FC = () => {
  const dispatch = useDispatch<any>();
  const { orders, isLoading } = useSelector((s: RootState) => s.order);

  const [search, setSearch] = useState("");
  const [editOrder, setEditOrder] = useState<Order | null>(null);

  useEffect(() => {
    dispatch(fetchOrders());
  }, []);

  const filtered = orders.filter((o) => {
    const s = search.toLowerCase();
    return (
      String(o.ID).includes(s) ||
      o.Status.toLowerCase().includes(s) ||
      o.PaymentStatus.toLowerCase().includes(s) ||
      o.Items.some((i) => i.Name.toLowerCase().includes(s))
    );
  });

  return (
    <Box>
      <CTable
        title="Orders"
        data={filtered}
        rowKey="ID"
        loading={isLoading}
        searchable
        onSearch={setSearch}
        searchPlaceholder="Search order, item, status"
        onRefresh={() => dispatch(fetchOrders())}
        columns={[
          {
            key: "ID",
            label: "Order",
            render: (o: Order) => <Text weight="bold">#{o.ID}</Text>,
          },
          {
            key: "CreatedAt",
            label: "Created At",
            render: (o: Order) => new Date(o.CreatedAt).toLocaleString(),
          },
          {
            key: "UserName",
            label: "By",
            render: (o: Order) => (
              <Text>
                {o.UserName}{" "}
                <Tooltip content={`Email : ${o.UserEmail}`}>
                  <Info size={12} />
                </Tooltip>
              </Text>
            ),
          },
          {
            key: "Items",
            label: "Items",
            render: (o: Order) => (
              <Flex gap="2" align="center">
                <Text size="2">
                  {o.Items.map((i) => `${i.Name} × ${i.Qty}`).join(", ")}
                </Text>

                <Tooltip
                  content={
                    <Box p="3">
                      {o.Items.map((i) => (
                        <Flex key={i.ItemID} justify="between" gap="3">
                          <Text>{i.Name} </Text>
                          <Text>{formatMoney(i.Total)}</Text>
                        </Flex>
                      ))}
                    </Box>
                  }
                >
                  <Info size={14} />
                </Tooltip>
              </Flex>
            ),
          },
          {
            key: "Status",
            label: "Status",
            render: (o: Order) => (
              <Badge
                color={
                  o.Status === "delivered"
                    ? "green"
                    : o.Status === "cancelled"
                    ? "red"
                    : "orange"
                }
              >
                {o.Status}
              </Badge>
            ),
          },
          {
            key: "PaymentStatus",
            label: "Payment",
            render: (o: Order) => (
              <Badge
                color={
                  o.PaymentStatus === "paid"
                    ? "green"
                    : o.PaymentStatus === "refunded"
                    ? "red"
                    : "gray"
                }
              >
                {o.PaymentStatus}
              </Badge>
            ),
          },
          {
            key: "PaymentMode",
            label: "Payment Mode",
            render: (o: Order) => <Badge>{o.PaymentMode}</Badge>,
          },
          {
            key: "Bill",
            label: "Total",
            align: "right",
            render: (o: Order) => (
              <Flex gap="2" align="center">
                <Text weight="bold">{formatMoney(o.Bill.TotalPayable)}</Text>
                <Tooltip
                  content={
                    <Box p="3">
                      <Flex justify="between" gap="3">
                        <Text>Subtotal</Text>
                        <Text>{formatMoney(o.Bill.Subtotal)}</Text>
                      </Flex>
                      <Flex justify="between">
                        <Text>Tax</Text>
                        <Text>{formatMoney(o.Bill.Tax)}</Text>
                      </Flex>
                      <Flex justify="between">
                        <Text>Platform</Text>
                        <Text>{formatMoney(o.Bill.PlatformFee)}</Text>
                      </Flex>
                      <Flex justify="between">
                        <Text>Discount</Text>
                        <Text>{formatMoney(o.Bill.Discount)}</Text>
                      </Flex>
                      <Flex justify="between" mt="2">
                        <Text weight="bold">Total</Text>
                        <Text weight="bold">
                          {formatMoney(o.Bill.TotalPayable)}
                        </Text>
                      </Flex>
                    </Box>
                  }
                >
                  <Info size={14} style={{ cursor: "pointer" }} />
                </Tooltip>
              </Flex>
            ),
          },
          {
            key: "Actions",
            label: "Actions",
            align: "center",
            render: (o: Order) => (
              <IconButton
                variant="ghost"
                radius="full"
                onClick={() => setEditOrder(o)}
                style={{ cursor: "pointer" }}
              >
                <Edit size={16} />
              </IconButton>
            ),
          },
        ]}
      />

      {editOrder && (
        <OrderActionForm
          order={editOrder}
          onCancel={() => setEditOrder(null)}
          onSuccess={() => dispatch(fetchOrders())}
        />
      )}
    </Box>
  );
};

export default OrdersPanel;

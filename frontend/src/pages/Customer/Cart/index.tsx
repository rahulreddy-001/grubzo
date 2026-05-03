import { Box, Flex, Text } from "@radix-ui/themes";
import { useSelector } from "react-redux";
import { useMemo } from "react";
import type { RootState } from "../../../services/store";
import cartService from "../../../services/cart/cart.service";
import { type UpdateItemQuantityPayload } from "../../../types/cart.d";
import CTable from "../../../components/common/CTable";
import { type Column } from "../../../components/common/CTable";
import CButton from "../../../components/common/CButton";
import { Minus, Plus } from "lucide-react";
import { useErrorHandler } from "../../../hooks/useErrorHandler";

type CartRow = {
  ItemID: number;
  Name: string;
  Price: number;
  Qty: number;
  Total: number;
};

const formatMoney = (paise: number) => `₹${(paise / 100).toFixed(2)}`;

function Cart() {
  const items = useSelector((state: RootState) => state.cart.items);
  const cartItems = useSelector((state: RootState) => state.cart.cart);
  const bill = useSelector((state: RootState) => state.cart.bill);
  const isLoading = useSelector((state: RootState) => state.cart.isLoading);
  const { showError, showSuccess } = useErrorHandler();

  const rows: CartRow[] = useMemo(() => {
    return cartItems
      .map((ci) => {
        const item = items.find((i) => i.ID === ci.Item);
        if (!item) return null;

        const price = Math.round(item.Price * 100);

        return {
          ItemID: item.ID,
          Name: item.Name,
          Price: price,
          Qty: ci.Quantity,
          Total: price * ci.Quantity,
        };
      })
      .filter(Boolean) as CartRow[];
  }, [cartItems, items]);

  const updateQty = (itemID: number, delta: number) => {
    const row = rows.find((r) => r.ItemID === itemID);
    const next = (row?.Qty || 0) + delta;

    cartService.updateCart({
      Item: itemID,
      Quantity: next < 0 ? 0 : next,
    } as UpdateItemQuantityPayload);
  };

  const handlePlaceOrder = async () => {
    try {
      const res = await cartService.submitOrder("wallet");
      cartService.fetchCart();
      showSuccess(res.Message);
    } catch (err) {
      showError(err);
    }
  };

  const columns: Column<CartRow>[] = [
    { key: "Name", label: "Item" },
    {
      key: "Price",
      label: "Price",
      render: (row) => formatMoney(row.Price),
    },
    {
      key: "Qty",
      label: "Qty",
      align: "center",
      render: (row) => (
        <Flex
          align="center"
          style={{
            background: "#1ba672",
            color: "white",
            borderRadius: 4,
            overflow: "hidden",
            height: 24,
            width: 70,
            justifyContent: "space-between",
            opacity: isLoading ? 0.5 : 1,
            pointerEvents: isLoading ? "none" : "auto",
            filter: isLoading ? "grayscale(40%)" : "none",
          }}
        >
          <Box
            onClick={() => updateQty(row.ItemID, -1)}
            style={{ width: 24, textAlign: "center", cursor: "pointer" }}
          >
            <Minus size={12} />
          </Box>
          <Box style={{ width: 20, textAlign: "center" }}>{row.Qty}</Box>
          <Box
            onClick={() => updateQty(row.ItemID, 1)}
            style={{ width: 24, textAlign: "center", cursor: "pointer" }}
          >
            <Plus size={12} />
          </Box>
        </Flex>
      ),
    },
    {
      key: "Total",
      label: "Total",
      render: (row) => formatMoney(row.Total),
      align: "right",
    },
  ];

  if (rows.length === 0) {
    return (
      <Flex
        align="center"
        justify="center"
        direction="column"
        style={{ height: "100%", color: "#888" }}
      >
        <Text size="4">Your cart is empty</Text>
      </Flex>
    );
  }

  return (
    <Box>
      <CTable
        title="Your Cart"
        data={rows}
        columns={columns}
        rowKey="ItemID"
        emptyMessage="Cart is empty"
        actions={
          <CButton
            disabled={isLoading}
            processing={isLoading}
            label="Place Order"
            color="primary"
            styles={{
              padding: "12px 28px",
              fontSize: "15px",
            }}
            onClick={handlePlaceOrder}
          />
        }
      />

      {bill && (
        <Box
          mt="3"
          p="3"
          style={{
            background: "var(--gray-2)",
            border: "1px solid var(--gray-5)",
          }}
        >
          <Text size="2" weight="bold" mb="2">
            Bill Details
          </Text>

          <Flex direction="column" gap="2">
            <BillRow label="Subtotal" value={bill.Subtotal} />
            <BillRow label={`Tax (${bill.TaxP}%)`} value={bill.Tax} />
            <BillRow
              label={`Platform Fee (${bill.PlatformFeeP}%)`}
              value={bill.PlatformFee}
            />

            {bill.Discount > 0 && (
              <BillRow
                label={`Discount (${bill.DiscountP}%)`}
                value={-bill.Discount}
              />
            )}

            <Box
              style={{ borderTop: "1px dashed var(--gray-6)", margin: "6px 0" }}
            />

            <BillRow label="Total Payable" value={bill.TotalPayable} bold />
          </Flex>
        </Box>
      )}
    </Box>
  );
}

export default Cart;

function BillRow({
  label,
  value,
  bold,
}: {
  label: string;
  value: number;
  bold?: boolean;
}) {
  return (
    <Flex justify="between" align="center">
      <Text size="2" color="gray">
        {label}
      </Text>
      <Text size="2" weight={bold ? "bold" : "medium"}>
        {formatMoney(value)}
      </Text>
    </Flex>
  );
}

import React from "react";
import { useFormik } from "formik";
import { Box, Flex, Text, Separator } from "@radix-ui/themes";

import CModel from "../../../components/common/CModel";
import CSelect from "../../../components/common/CSelect";
import CButton from "../../../components/common/CButton";
import { useErrorHandler } from "../../../hooks/useErrorHandler";

import type { Order } from "../../../types/order";
import orderService from "../../../services/order/order.service";

const formatMoney = (p: number) => `₹${(p / 100).toFixed(2)}`;

const OrderStatusOptions = [
  { label: "Pending", value: "pending" },
  { label: "Preparing", value: "preparing" },
  { label: "Ready", value: "ready" },
  { label: "Delivered", value: "delivered" },
  { label: "Cancelled", value: "cancelled" },
];

const PaymentStatusOptions = [
  { label: "Pending", value: "pending" },
  { label: "Paid", value: "paid" },
  { label: "Refunded", value: "refunded" },
  { label: "Voided", value: "voided" },
];

const OrderActionForm: React.FC<{
  order: Order;
  onCancel: () => void;
  onSuccess: () => void;
}> = ({ order, onCancel, onSuccess }) => {
  const { showError, showSuccess } = useErrorHandler();

  const form = useFormik({
    initialValues: {
      Status: order.Status,
      PaymentStatus: order.PaymentStatus,
    },
    onSubmit: async (values) => {
      try {
        let res = await orderService.updateStatus({
          OrderID: order.ID,
          OrderStatus: values.Status,
          PaymentStatus: values.PaymentStatus,
        });

        showSuccess(res.Message || "Order updated");
        onCancel();
        onSuccess();
      } catch (err: any) {
        showError(err);
      }
    },
  });

  return (
    <CModel
      open
      size="md"
      title={`Order #${order.ID}`}
      onClose={onCancel}
      actions={
        <Flex gap="2">
          <CButton label="Close" variant="soft" onClick={onCancel} />
          <CButton
            label="Update"
            onClick={form.submitForm}
            processing={form.isSubmitting}
          />
        </Flex>
      }
    >
      <form onSubmit={form.handleSubmit}>
        <Flex direction="column" gap="3">
          <Flex direction="column" gap="1">
            <Text size="2" weight="medium">
              {order.UserName}
            </Text>
            <Text size="1" color="gray">
              {order.UserEmail}
            </Text>
          </Flex>

          <Separator size="4" style={{ width: "100%" }} />

          <Box>
            {order.Items.map((i) => (
              <Flex key={i.ItemID} justify="between">
                <Text size="2">
                  {i.Name} × {i.Qty}
                </Text>
                <Text size="2">{formatMoney(i.Total)}</Text>
              </Flex>
            ))}
          </Box>

          <Separator size="4" style={{ width: "100%" }} />

          <Box>
            <Flex justify="between">
              <Text>Subtotal</Text>
              <Text>{formatMoney(order.Bill.Subtotal)}</Text>
            </Flex>
            <Flex justify="between">
              <Text>Tax</Text>
              <Text>{formatMoney(order.Bill.Tax)}</Text>
            </Flex>
            <Flex justify="between">
              <Text>Platform</Text>
              <Text>{formatMoney(order.Bill.PlatformFee)}</Text>
            </Flex>
            {order.Bill.Discount > 0 && (
              <Flex justify="between">
                <Text color="green">Discount</Text>
                <Text color="green">−{formatMoney(order.Bill.Discount)}</Text>
              </Flex>
            )}
          </Box>

          <Separator size="4" style={{ width: "100%" }} />

          <Flex justify="between">
            <Text weight="medium">Total</Text>
            <Text weight="medium">{formatMoney(order.Bill.TotalPayable)}</Text>
          </Flex>

          <Separator size="4" style={{ width: "100%" }} />

          <Flex direction="column" gap="2">
            <CSelect
              label="Order Status"
              value={form.values.Status}
              options={OrderStatusOptions}
              onChange={(v) => form.setFieldValue("Status", v)}
            />
            {order.PaymentMode == "pos" && (
              <CSelect
                label="Payment Status"
                value={form.values.PaymentStatus}
                options={PaymentStatusOptions}
                onChange={(v) => form.setFieldValue("PaymentStatus", v)}
              />
            )}
          </Flex>
        </Flex>
      </form>
    </CModel>
  );
};

export default OrderActionForm;

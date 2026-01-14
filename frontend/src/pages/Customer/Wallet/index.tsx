import { useEffect, useMemo, useState } from "react";
import { Box, Flex, Text, Badge } from "@radix-ui/themes";
import { Plus } from "lucide-react";
import { useRazorpay } from "react-razorpay";

import CTable, { type Column } from "../../../components/common/CTable";
import CButton from "../../../components/common/CButton";
import CModel from "../../../components/common/CModel";
import CInput from "../../../components/common/CInput";
import { useErrorHandler } from "../../../hooks/useErrorHandler";

const formatMoney = (p: number) => `₹${(p / 100).toFixed(2)}`;

type WalletTransactionDTO = {
  TenantID: number;
  UserID: number;
  Amount: number;
  BalanceAfter: number;
  Type: "credit" | "debit";
  ReferenceType: string;
  OrderID: string;
  IdempotentID: string;
  CreatedAt: string;
  Message: string;
};

type PendingWalletRechargeDTO = {
  ID: number;
  Amount: number;
  PaymentGateway: string;
  PaymentReference: string;
};

type WalletResponseDTO = {
  Wallet: {
    Balance: number;
    PendingRecharges: PendingWalletRechargeDTO[];
    Transactions: WalletTransactionDTO[];
  };
  Message: string;
};

type WalletRow = WalletTransactionDTO & {
  Sign: "+" | "-";
};

function Wallet() {
  const { Razorpay, isLoading: isRazorpayLoading, error } = useRazorpay();

  const [wallet, setWallet] = useState<WalletResponseDTO | null>(null);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [amount, setAmount] = useState("");
  const [loading, setLoading] = useState(false);
  const { showError, showSuccess } = useErrorHandler();

  const loadWallet = async () => {
    try {
      const res = await fetch("/api/v1/wallet/");
      const data: WalletResponseDTO = await res.json();

      data.Wallet.PendingRecharges = data.Wallet.PendingRecharges ?? [];
      data.Wallet.Transactions = data.Wallet.Transactions ?? [];

      setWallet(data);
    } catch (err) {
      console.error("Failed to load wallet", err);
    }
  };

  useEffect(() => {
    loadWallet();
  }, []);

  const openRazorpay = (order: any) => {
    if (!Razorpay) {
      console.error("Razorpay not ready", error);
      return;
    }

    const options = {
      key: order.key,
      amount: order.amount,
      currency: order.currency,
      name: "Grubzo Wallet",
      description: "Wallet Recharge",
      order_id: order.id,

      handler: async (response: any) => {
        try {
          await fetch("/api/v1/wallet/verify_recharge", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              OrderID: response.razorpay_order_id,
              PaymentReference: response.razorpay_payment_id,
              Signature: response.razorpay_signature,
            }),
          });
          showSuccess("Transaction completed successfully");
          await loadWallet();
        } catch (err) {
          showError("Payment verified but wallet update failed");
          console.error(err);
        }
      },

      theme: { color: "#3b82f6" },
    };

    new Razorpay(options).open();
  };

  const handleAddMoney = async () => {
    const value = Number(amount);
    if (!value || value <= 0) return;

    try {
      setLoading(true);

      const res = await fetch("/api/v1/wallet/recharge_order", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ Amount: value }),
      });

      const order = await res.json();
      openRazorpay(order);
      setDrawerOpen(false);
      setAmount("");
    } finally {
      setLoading(false);
    }
  };

  const tableData: WalletRow[] = useMemo(() => {
    return (
      wallet?.Wallet.Transactions.map((t) => ({
        ...t,
        Sign: t.Type === "credit" ? "+" : "-",
      })) ?? []
    );
  }, [wallet]);

  const columns: Column<WalletRow>[] = [
    { key: "IdempotentID", label: "Ref", width: 140 },
    {
      key: "CreatedAt",
      label: "Date & Time",
      render: (r) => new Date(r.CreatedAt).toLocaleString(),
    },
    { key: "Message", label: "Message" },
    {
      key: "Type",
      label: "Type",
      render: (r) => (
        <Badge color={r.Type === "credit" ? "green" : "red"}>{r.Type}</Badge>
      ),
    },
    {
      key: "Amount",
      label: "Amount",
      align: "right",
      render: (r) => (
        <Text weight="bold" color={r.Type === "credit" ? "green" : "red"}>
          {r.Sign}
          {formatMoney(r.Amount)}
        </Text>
      ),
    },
    {
      key: "BalanceAfter",
      label: "Balance",
      align: "right",
      render: (r) => <Text>{formatMoney(r.BalanceAfter)}</Text>,
    },
  ];

  return (
    <Box>
      <CTable
        title={`Wallet Balance: ${formatMoney(wallet?.Wallet.Balance ?? 0)}`}
        data={tableData}
        columns={columns}
        rowKey="IdempotentID"
        loading={!wallet}
        onRefresh={loadWallet}
        emptyMessage="No wallet activity yet"
        actions={
          <CButton
            label="Add Money"
            startIcon={<Plus size={16} />}
            onClick={() => setDrawerOpen(true)}
          />
        }
      />

      {drawerOpen && (
        <CModel
          open
          size="sm"
          title="Add Money to Wallet"
          onClose={() => setDrawerOpen(false)}
          actions={
            <Flex gap="3">
              <CButton
                label="Cancel"
                variant="soft"
                onClick={() => setDrawerOpen(false)}
              />
              <CButton
                label="Add Money"
                onClick={handleAddMoney}
                processing={loading || isRazorpayLoading}
              />
            </Flex>
          }
        >
          <Box>
            <CInput
              label="Amount"
              placeholder="Enter amount"
              value={amount}
              onChange={(v) => setAmount(v)}
              fullWidth
            />
          </Box>
        </CModel>
      )}
    </Box>
  );
}

export default Wallet;

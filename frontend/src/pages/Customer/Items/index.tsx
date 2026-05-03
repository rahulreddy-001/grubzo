import { Box, Container, Flex, Text, TextField } from "@radix-ui/themes";
import { useMemo, useRef, useState } from "react";
import { useSelector } from "react-redux";
import type { RootState } from "../../../services/store";
import type { Item } from "../../../types/item";
import cartService from "../../../services/cart/cart.service";
import type { UpdateItemQuantityPayload } from "../../../types/cart.d";
import { Minus, Plus, Search } from "lucide-react";

function Items() {
  const items = useSelector((state: RootState) => state.cart.items);
  const cartItems = useSelector((state: RootState) => state.cart.cart);

  const [search, setSearch] = useState("");
  const [isCartBusy, setIsCartBusy] = useState(false);

  const categoryRefs = useRef<Record<string, HTMLDivElement | null>>({});

  const getQty = (itemID: number) => {
    const found = cartItems.find((i) => i.Item === itemID);
    return found ? found.Quantity : 0;
  };

  const updateQty = async (itemID: number, delta: number) => {
    if (isCartBusy) return;

    const current = getQty(itemID);
    const next = Math.max(0, current + delta);

    setIsCartBusy(true);
    try {
      await cartService.updateCart({
        Item: itemID,
        Quantity: next,
      } as UpdateItemQuantityPayload);
    } finally {
      setIsCartBusy(false);
    }
  };

  const groupedItems = useMemo(() => {
    const map: Record<string, Item[]> = {};

    for (const item of items) {
      const match =
        item.Name.toLowerCase().includes(search.toLowerCase()) ||
        item.Description?.toLowerCase().includes(search.toLowerCase());

      if (!match) continue;

      if (!map[item.Category]) map[item.Category] = [];
      map[item.Category].push(item);
    }

    return map;
  }, [items, search]);

  const categories = Object.keys(groupedItems);

  return (
    <Container size="4" style={{ height: "100%", padding: 0 }}>
      <Flex
        justify="between"
        align="center"
        style={{
          padding: "12px 20px",
          borderBottom: "1px solid #eee",
          position: "sticky",
          top: 0,
          background: "white",
          zIndex: 10,
        }}
      >
        <Text size="4" weight="bold">
          Menu
        </Text>

        <TextField.Root
          placeholder="Search dishes"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        >
          <TextField.Slot>
            <Search height="16" width="16" />
          </TextField.Slot>
        </TextField.Root>
      </Flex>

      <Box
        style={{
          height: "calc(100% - 56px)",
          overflowY: "auto",
          padding: "16px 24px",
          pointerEvents: isCartBusy ? "none" : "auto",
          opacity: isCartBusy ? 0.6 : 1,
        }}
      >
        {categories.map((category) => (
          <Box
            key={category}
            ref={(el) => {
              categoryRefs.current[category] = el;
            }}
            style={{ marginBottom: 40 }}
          >
            <Text size="5" weight="bold">
              {category}
            </Text>

            {groupedItems[category].map((item) => {
              const qty = getQty(item.ID);
              const isOutOfStock = item.ItemStatus === "os";

              return (
                <Flex
                  key={item.ID}
                  justify="between"
                  style={{
                    padding: "20px 0",
                    borderBottom: "1px solid #f1f1f1",
                  }}
                >
                  <Box style={{ flex: 1 }}>
                    <Text as="div" weight="bold">
                      {item.Name}
                    </Text>
                    <Text as="div" size="2" color="gray">
                      {item.Description}
                    </Text>
                    <Text as="div" size="2">
                      ₹{item.Price}
                    </Text>
                  </Box>

                  <Box style={{ width: 140, position: "relative" }}>
                    <img
                      src={item.Files?.[0]?.URL}
                      style={{
                        width: "100%",
                        height: 110,
                        borderRadius: 4,
                        objectFit: "cover",
                      }}
                    />

                    {isOutOfStock ? (
                      <Box
                        style={{
                          position: "absolute",
                          bottom: -12,
                          left: "50%",
                          transform: "translateX(-50%)",
                          background: "#e5e7eb",
                          color: "#6b7280",
                          borderRadius: 4,
                          padding: "6px 14px",
                          fontSize: 12,
                          fontWeight: 600,
                          cursor: "not-allowed",
                          whiteSpace: "nowrap",
                        }}
                      >
                        OUT OF STOCK
                      </Box>
                    ) : qty === 0 ? (
                      <Box
                        onClick={() => updateQty(item.ID, 1)}
                        style={{
                          position: "absolute",
                          bottom: -12,
                          left: "50%",
                          transform: "translateX(-50%)",
                          background: "#fff",
                          borderRadius: 4,
                          padding: "6px 18px",
                          cursor: isCartBusy ? "not-allowed" : "pointer",
                          color: "#1ba672",
                          backgroundColor: "#ecfdf5",
                          fontSize: 12,
                          fontWeight: 600,
                          pointerEvents: isCartBusy ? "none" : "auto",
                          opacity: isCartBusy ? 0.5 : 1,
                        }}
                      >
                        ADD
                      </Box>
                    ) : (
                      <Flex
                        align="center"
                        style={{
                          position: "absolute",
                          bottom: -8,
                          left: "50%",
                          transform: "translateX(-50%)",
                          background: "#1ba672",
                          color: "white",
                          borderRadius: 4,
                          overflow: "hidden",
                          height: 28,
                          width: 80,
                          fontSize: 12,
                          justifyContent: "space-between",
                          opacity: isCartBusy ? 0.5 : 1,
                          pointerEvents: isCartBusy ? "none" : "auto",
                          filter: isCartBusy ? "grayscale(40%)" : "none",
                        }}
                      >
                        <Box
                          onClick={() => updateQty(item.ID, -1)}
                          style={{
                            width: 24,
                            textAlign: "center",
                            cursor: "pointer",
                          }}
                        >
                          <Minus size={12} />
                        </Box>

                        <Box style={{ width: 8, textAlign: "center" }}>
                          {qty}
                        </Box>

                        <Box
                          onClick={() => updateQty(item.ID, 1)}
                          style={{
                            width: 24,
                            textAlign: "center",
                            cursor: "pointer",
                          }}
                        >
                          <Plus size={12} />
                        </Box>
                      </Flex>
                    )}
                  </Box>
                </Flex>
              );
            })}
          </Box>
        ))}
      </Box>
    </Container>
  );
}

export default Items;

import { useState, type FC } from "react";
import { Box, Flex, TabNav } from "@radix-ui/themes";
import {
  Menu,
  ShoppingCart,
  ShoppingBag,
  Wallet as WalletIcon,
  type LucideIcon,
} from "lucide-react";

import Items from "./Items";
import Cart from "./Cart";
import Orders from "./Orders";
import Wallet from "./Wallet";

type TabConfig = {
  label: string;
  icon: LucideIcon;
  component: FC;
};

const TAB_CONFIG: TabConfig[] = [
  {
    label: "Menu",
    icon: Menu,
    component: Items,
  },
  {
    label: "Cart",
    icon: ShoppingCart,
    component: Cart,
  },
  {
    label: "Orders",
    icon: ShoppingBag,
    component: Orders,
  },
  {
    label: "Wallet",
    icon: WalletIcon,
    component: Wallet,
  },
];

function UserHome() {
  const [activeTab, setActiveTab] = useState(0);

  const Empty: FC = () => null;
  const ActiveComponent =
    TAB_CONFIG.length > 0 ? TAB_CONFIG[activeTab].component : Empty;

  return (
    <Box>
      <TabNav.Root>
        {TAB_CONFIG.map((tab, index) => (
          <TabNav.Link
            key={tab.label}
            onClick={() => setActiveTab(index)}
            active={activeTab === index}
            style={{ cursor: "pointer" }}
          >
            <Flex gap="3" style={{ margin: 10 }} align="center">
              <tab.icon size={18} strokeWidth={2} />
              {tab.label}
            </Flex>
          </TabNav.Link>
        ))}
      </TabNav.Root>

      <Box mt="4">
        <ActiveComponent />
      </Box>
    </Box>
  );
}

export default UserHome;

import React, { useState } from "react";

import {
  Home,
  ClipboardList,
  UtensilsCrossed,
  Users,
  UserLock,
  MapPin,
} from "lucide-react";
import { Box, Flex, TabNav } from "@radix-ui/themes";

import DashboardPanel from "./DashboardPanel";
import OrdersPanel from "./OrdersPanel";
import ItemsPanel from "./ItemsPanel";
import EmployeesPanel from "./EmployeesPanel";
import RBACPanel from "./RBACPanel";
import LocationPanel from "./LocationPanel";
import { useSelector } from "react-redux";
import type { RootState } from "../../services/store";

const TAB_CONFIG = [
  {
    label: "Dashboard",
    icon: Home,
    component: DashboardPanel,
    permissions: ["dashboard"],
  },
  {
    label: "Orders",
    icon: ClipboardList,
    component: OrdersPanel,
    permissions: ["orders"],
  },
  {
    label: "Items",
    icon: UtensilsCrossed,
    component: ItemsPanel,
    permissions: ["items"],
  },
  {
    label: "Employees",
    icon: Users,
    component: EmployeesPanel,
    permissions: ["employee"],
  },
  {
    label: "Location",
    icon: MapPin,
    component: LocationPanel,
    permissions: ["location"],
  },
  {
    label: "Access Control",
    icon: UserLock,
    component: RBACPanel,
    permissions: ["rbac"],
  },
];

function EmployeeHome() {
  const [activeTab, setActiveTab] = useState(0);
  let { user } = useSelector((s: RootState) => s.auth);
  const visibleTabs = TAB_CONFIG.filter((t) =>
    t.permissions.some((perm) => user?.Permisssions?.includes(perm))
  );
  const Empty: React.FC = () => null;
  const ActiveComponent =
    visibleTabs.length != 0 ? visibleTabs[activeTab].component : Empty;

  return (
    <Box>
      <TabNav.Root>
        {visibleTabs.map((tab, index) => (
          <TabNav.Link
            key={tab.label}
            onClick={() => setActiveTab(index)}
            active={activeTab == index}
            style={{ cursor: "pointer" }}
          >
            <Flex gap="3" style={{ margin: 10 }} align="center">
              <tab.icon size={18} strokeWidth={2} />
              {tab.label}
            </Flex>
          </TabNav.Link>
        ))}
      </TabNav.Root>
      <ActiveComponent />
    </Box>
  );
}

export default EmployeeHome;

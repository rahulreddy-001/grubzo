import React, { useEffect, useState } from "react";
import { useSelector } from "react-redux";
import type { RootState } from "../../../services/store";

import CTable from "../../../components/common/CTable";
import CButton from "../../../components/common/CButton";
import { Box, Flex, Text, IconButton, Badge } from "@radix-ui/themes";
import { Plus, Edit } from "lucide-react";
import EmployeeForm from "./EmployeeForm";

import CommonService from "../../../services/common/common.service";

const ItemsPanel: React.FC = () => {
  const { Employees: employees, isLoading } = useSelector(
    (s: RootState) => s.emp
  );

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editItem, setEditItem] = useState<any>(null);
  const [search, setSearch] = useState("");

  useEffect(() => {
    CommonService.fetchEmployees();
    setTimeout(() => {
      console.log(employees);
    }, 2000);
  }, []);

  const handleAdd = () => {
    setEditItem(null);
    setDrawerOpen(true);
  };

  const handleEdit = (item: any) => {
    setEditItem(item);
    setDrawerOpen(true);
  };

  const filteredItems = employees?.filter((emp) =>
    emp.Name.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <Box>
      <CTable
        title="Employees"
        data={filteredItems ?? []}
        rowKey="ID"
        loading={isLoading}
        searchable
        onSearch={setSearch}
        searchPlaceholder="Search employees"
        onRefresh={() => CommonService.fetchEmployees()}
        actions={
          <CButton
            label="Create Employee"
            startIcon={<Plus size={16} />}
            onClick={handleAdd}
          />
        }
        columns={[
          {
            key: "Name",
            label: "Name",
            render: (row) => (
              <Box>
                <Text weight="medium">{row.Name}</Text>
              </Box>
            ),
          },
          {
            key: "Email",
            label: "Email",
            render: (row) => (
              <Box>
                <Text weight="medium">{row.Email}</Text>
              </Box>
            ),
          },
          {
            key: "Roles",
            label: "Roles",
            render: (row) => (
              <Flex gap="2">
                {row.Roles.map((role) => (
                  <Badge
                    key={role}
                    color="blue"
                    size="2"
                    style={{ cursor: "pointer" }}
                  >
                    {role}
                  </Badge>
                ))}
              </Flex>
            ),
          },
          {
            key: "Actions",
            label: "Actions",
            render: (item) => (
              <Flex gap="2" style={{ marginLeft: "10px", cursor: "pointer" }}>
                <IconButton
                  variant="ghost"
                  radius="full"
                  style={{ cursor: "pointer", padding: "5px", margin: "1px" }}
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
        <EmployeeForm item={editItem} cancel={() => setDrawerOpen(false)} />
      )}
    </Box>
  );
};

export default ItemsPanel;

import React, { useState } from "react";
import { useSelector } from "react-redux";
import type { RootState } from "../../../services/store";
import { Plus, Save } from "lucide-react";
import CTable from "../../../components/common/CTable";
import CButton from "../../../components/common/CButton";
import CMultiSelect from "../../../components/common/CMultiSelect";
import { useErrorHandler } from "../../../hooks/useErrorHandler";
import { Box, Flex, Text } from "@radix-ui/themes";
import CreateRoleForm from "./CreateRoleForm";
import CommonService from "../../../services/common/common.service";

const RBACPanel: React.FC = () => {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const { showSuccess, showError } = useErrorHandler();

  const { Grid, Permissions, isLoading } = useSelector(
    (s: RootState) => s.rbac,
  );

  React.useEffect(() => {
    CommonService.fetchRBACInfo();
  }, []);

  const [grid, setGrid] = useState<
    { ID: number; Role: string; Permissions: string[] }[]
  >([]);

  React.useEffect(() => {
    const rows = Object.entries(Grid ?? []).map(([role, perms], idx) => ({
      ID: idx + 1,
      Role: role,
      Permissions: perms,
    }));
    setGrid(rows);
  }, [Grid]);

  const updatePermissions = (roleID: number, newPermissions: string[]) => {
    setGrid((prev) =>
      prev.map((r) =>
        r.ID === roleID ? { ...r, Permissions: newPermissions } : r,
      ),
    );
  };

  const handleSave = async () => {
    try {
      const invalid = grid.find((r) => r.Permissions.length === 0);
      if (invalid) {
        return showError(`Role ${invalid.Role} must have permissions.`);
      }

      const original = Grid || {};
      const diff: any[] = [];

      const map: Record<string, Set<string>> = {};
      grid.forEach((r) => (map[r.Role] = new Set(r.Permissions)));

      const roles = new Set([
        ...Object.keys(original),
        ...grid.map((r) => r.Role),
      ]);

      roles.forEach((role) => {
        const origSet = new Set(
          Array.isArray(original[role]) ? original[role] : [],
        );
        const currSet = map[role] || new Set();

        const add = [...currSet].filter((p) => !origSet.has(p));
        const rem = [...origSet].filter((p) => !currSet.has(p));

        if (add.length) diff.push({ Name: role, Permissions: add, Action: 1 });
        if (rem.length) diff.push({ Name: role, Permissions: rem, Action: 0 });
      });

      if (!diff.length) return showError("No changes detected.");

      const res = await CommonService.syncRBACInfo({ Data: diff });
      showSuccess(res.Message);
      await CommonService.fetchRBACInfo();
    } catch (err) {
      showError(err);
    }
  };

  return (
    <Box py="1">
      <CTable
        title="Access Control"
        data={grid}
        rowKey="ID"
        loading={isLoading}
        actions={
          <Flex gap="1">
            <CButton
              label="Create Role"
              startIcon={<Plus size={16} />}
              onClick={() => setDrawerOpen(true)}
            />
            <CButton
              label="Save"
              startIcon={<Save size={16} />}
              onClick={handleSave}
            />
          </Flex>
        }
        columns={[
          {
            key: "Role",
            label: "Role",
            render: (row) => (
              <Text as="div" style={{ marginLeft: "5px" }}>
                {row.Role}
              </Text>
            ),
          },
          {
            key: "Permissions",
            label: "Permissions",
            render: (row) => (
              <Flex>
                <CMultiSelect
                  placeholder="Select Permissions"
                  selected={row.Permissions}
                  options={Permissions}
                  onChange={(change) => updatePermissions(row.ID, change)}
                />
              </Flex>
            ),
          },
        ]}
      />

      {drawerOpen && <CreateRoleForm close={() => setDrawerOpen(false)} />}
    </Box>
  );
};

export default RBACPanel;

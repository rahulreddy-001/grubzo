import React, { useState, useCallback, useMemo } from "react";
import { useSelector } from "react-redux";
import type { RootState } from "../../../services/store";

import CModel from "../../../components/common/CModel";
import CInput from "../../../components/common/CInput";
import CButton from "../../../components/common/CButton";
import CMultiSelect from "../../../components/common/CMultiSelect";

import { Flex, Box, Text } from "@radix-ui/themes";

import { useErrorHandler } from "../../../hooks/useErrorHandler";
import CommonService from "../../../services/common/common.service";

interface CreateRoleProps {
  close: () => void;
}

const CreateRole: React.FC<CreateRoleProps> = ({ close }) => {
  const { Permissions: AllPermissions } = useSelector((s: RootState) => s.rbac);

  const permissionList = useMemo(
    () => (Array.isArray(AllPermissions) ? AllPermissions : []),
    [AllPermissions]
  );

  const [form, setForm] = useState({
    Name: "",
    Permissions: [] as string[],
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const { showError, showSuccess } = useErrorHandler();

  const validate = useCallback(() => {
    const e: Record<string, string> = {};

    if (!form.Name.trim()) e.Name = "Role name is required.";
    if (form.Permissions.length === 0)
      e.Permissions = "Select at least one permission.";

    return e;
  }, [form]);

  const handleSubmit = useCallback(async () => {
    const validationErrors = validate();
    setErrors(validationErrors);

    if (Object.keys(validationErrors).length > 0) return;

    try {
      const { Message } = await CommonService.createRole(form);
      showSuccess(Message);
      await CommonService.fetchRBACInfo();
      close();
    } catch (err) {
      showError(err);
    }
  }, [form, validate, close, showError, showSuccess]);

  return (
    <CModel
      open={true}
      title="Create Role"
      onClose={close}
      size="md"
      actions={
        <Flex gap="3">
          <CButton label="Cancel" variant="soft" onClick={close} />
          <CButton
            label="Create"
            variant="solid"
            onClick={handleSubmit}
            disabled={!form.Name.trim()}
          />
        </Flex>
      }
    >
      <Box>
        <Flex direction="column" gap="4">
          <CInput
            label="Role Name"
            value={form.Name}
            placeholder="Enter role name"
            onChange={(val) => {
              setForm((p) => ({ ...p, Name: val }));
              setErrors((e) => ({ ...e, Name: "" }));
            }}
            error={errors.Name}
            fullWidth
          />

          <Box>
            <Text size="1" weight="medium">
              Permissions
            </Text>

            <CMultiSelect
              selected={form.Permissions}
              options={permissionList}
              placeholder="Select Permissisons"
              onChange={(vals) => {
                setForm((p) => ({ ...p, Permissions: vals }));
                setErrors((e) => ({ ...e, Permissions: "" }));
              }}
            />

            {errors.Permissions && (
              <Text size="1" color="red" mt="1">
                {errors.Permissions}
              </Text>
            )}
          </Box>
        </Flex>
      </Box>
    </CModel>
  );
};

export default CreateRole;

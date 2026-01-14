import React from "react";
import { useSelector } from "react-redux";
import type { RootState } from "../../../services/store";

import CModel from "../../../components/common/CModel";
import CInput from "../../../components/common/CInput";
import CButton from "../../../components/common/CButton";
import CMultiSelect from "../../../components/common/CMultiSelect";

import { Flex, Box } from "@radix-ui/themes";

import { useErrorHandler } from "../../../hooks/useErrorHandler";
import CommonService from "../../../services/common/common.service";

import { useFormik } from "formik";
import { zodToFormik } from "../../../utils/zodToFormik";
import {
  CreateRolePayloadSchema,
  type CreateRolePayload,
} from "../../../types/common.d";

const CreateRoleForm: React.FC<{ close: () => void }> = ({ close }) => {
  const { showError, showSuccess } = useErrorHandler();
  const { Permissions } = useSelector((s: RootState) => s.rbac);

  const form = useFormik<CreateRolePayload>({
    initialValues: {
      Name: "",
      Permissions: [],
    },
    validate: zodToFormik(CreateRolePayloadSchema),
    onSubmit: async (values) => {
      try {
        const response = await CommonService.createRole(values);
        showSuccess(response.Message);
        await CommonService.fetchRBACInfo();
        close();
      } catch (err) {
        showError(err);
      }
    },
  });

  return (
    <CModel
      open={true}
      title="Create Role"
      onClose={close}
      size="md"
      actions={
        <Flex gap="3">
          <CButton label="Cancel" onClick={close} />
          <CButton
            label="Create"
            onClick={form.submitForm}
            disabled={!form.isValid || form.isSubmitting}
          />
        </Flex>
      }
    >
      <Box>
        <Flex direction="column" gap="4">
          <CInput
            label="Role Name"
            value={form.values.Name}
            placeholder="Enter role name"
            onChange={(v) => form.setFieldValue("Name", v)}
            error={form.errors.Name}
            fullWidth
          />

          <Box>
            <CMultiSelect
              label="Permissions"
              placeholder="Select Permissisons"
              options={Permissions ?? []}
              selected={form.values.Permissions}
              onChange={(vals) => form.setFieldValue("Permissions", vals)}
              error={form.errors.Permissions as string}
            />
          </Box>
        </Flex>
      </Box>
    </CModel>
  );
};

export default CreateRoleForm;

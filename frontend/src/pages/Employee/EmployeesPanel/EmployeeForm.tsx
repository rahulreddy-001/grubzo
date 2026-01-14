import React, { useEffect } from "react";
import CModel from "../../../components/common/CModel";
import CInput from "../../../components/common/CInput";
import CButton from "../../../components/common/CButton";
import CMultiSelect from "../../../components/common/CMultiSelect";
import { Box, Flex } from "@radix-ui/themes";

import { useFormik } from "formik";
import { zodToFormik } from "../../../utils/zodToFormik";
import { useSelector } from "react-redux";
import type { RootState } from "../../../services/store";
import { useErrorHandler } from "../../../hooks/useErrorHandler";

import { EmployeeSchema, type Employee } from "../../../types/common.d";
import CommonService from "../../../services/common/common.service";

const EmployeeForm: React.FC<{
  emp?: Employee;
  cancel: () => void;
}> = ({ emp, cancel }) => {
  const isUpdate = !!emp;
  const { showError, showSuccess } = useErrorHandler();
  const { Roles } = useSelector((s: RootState) => s.rbac);

  useEffect(() => {
    if (Roles.length == 0) {
      CommonService.fetchRBACInfo();
    }
  }, []);

  const form = useFormik<Employee>({
    initialValues: {
      ID: emp?.ID ?? -1,
      LocationID: emp?.LocationID ?? -1,
      Name: emp?.Name ?? "",
      Email: emp?.Email ?? "",
      Roles: emp?.Roles || [],
    },
    validate: zodToFormik(EmployeeSchema),
    onSubmit: async (employee: Employee) => {
      try {
        let response;
        if (isUpdate) {
          response = await CommonService.updateEmployee(employee);
        } else {
          response = await CommonService.createEmployee(employee);
        }
        showSuccess(response.Message);
        cancel();
        CommonService.fetchEmployees();
      } catch (err) {
        showError(err);
      }
    },
  });

  return (
    <CModel
      open
      size="md"
      onClose={cancel}
      title={isUpdate ? "Update Employee" : "Create Employee"}
      actions={
        <Flex gap="3">
          <CButton
            label="Cancel"
            variant="soft"
            onClick={cancel}
            disabled={form.isSubmitting}
          />
          <CButton
            label={isUpdate ? "Update" : "Add"}
            variant="solid"
            onClick={form.submitForm}
            disabled={form.isSubmitting || !form.isValid || !!form.values.Name}
            processing={form.isSubmitting}
          />
        </Flex>
      }
    >
      <form onSubmit={form.handleSubmit}>
        <Flex direction="column" gap="4">
          <Box>
            <CInput
              label="Employee Name"
              placeholder="Enter employee name"
              value={form.values.Name}
              onChange={(val) => form.setFieldValue("Name", val)}
              error={form.errors.Name}
              fullWidth
            />
          </Box>

          <Box>
            <CInput
              label="Employee Email"
              placeholder="Enter employee email"
              value={form.values.Email}
              onChange={(val) => form.setFieldValue("Email", val)}
              error={form.errors.Email}
              fullWidth
            />
          </Box>

          <Box>
            <CMultiSelect
              label="Employee Roles"
              placeholder="Select Roles"
              selected={form.values.Roles ?? []}
              options={Roles}
              onChange={(vals) => form.setFieldValue("Roles", vals)}
              error={form.errors.Roles as string}
            />
          </Box>
        </Flex>
      </form>
    </CModel>
  );
};

export default EmployeeForm;

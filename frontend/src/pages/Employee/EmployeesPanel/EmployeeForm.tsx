import React, { useEffect, useState, useCallback } from "react";
import CModel from "../../../components/common/CModel";
import CInput from "../../../components/common/CInput";
import CButton from "../../../components/common/CButton";
import CSelect from "../../../components/common/CSelect";
import CMultiSelect from "../../../components/common/CMultiSelect";
import { Box, Flex, Text } from "@radix-ui/themes";

import { useSelector } from "react-redux";
import type { RootState } from "../../../services/store";

import type { Employee } from "../../../types/common";
import CommonService from "../../../services/common/common.service";
import { useErrorHandler } from "../../../hooks/useErrorHandler";

const ModifyEmployee: React.FC<{
  item?: Employee;
  cancel: () => void;
}> = ({ item, cancel }) => {
  const { showError, showSuccess } = useErrorHandler();
  const [isProcessing, setIsProcessing] = useState(false);
  const { Locations } = useSelector((s: RootState) => s.common);
  const { Roles } = useSelector((s: RootState) => s.rbac);

  const [form, setForm] = useState<Partial<Employee>>({
    Name: "",
    Email: "",
    Roles: [],
    LocationID: 0,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    if (item) {
      setForm({
        Name: item.Name || "",
        Email: item.Email || "",
        Roles: item.Roles || [],
        LocationID: item.LocationID || 0,
      });
    }

    if (Locations?.length == 0) {
      CommonService.fetchLocations();
    }
    if (Roles.length == 0) {
      CommonService.fetchRBACInfo();
    }
  }, []);

  const handleChange = useCallback(
    <K extends keyof Employee>(key: K, value: Employee[K]) => {
      setForm((prev) => ({ ...prev, [key]: value }));
    },
    []
  );

  const validate = useCallback(() => {
    const e: Record<string, string> = {};

    if (!form.Name?.trim()) e.Name = "Employee name is required.";
    if (!form.Email?.trim()) e.Email = "Email is required.";
    if (!form.Roles || form.Roles.length === 0) e.Roles = "Roles are required.";
    if (!form.LocationID || form.LocationID === 0)
      e.LocationID = "Location is required.";

    return e;
  }, [form]);

  const isEditing = !!item;

  const handleSubmit = async () => {
    const v = validate();
    setErrors(v);

    if (Object.keys(v).length > 0) return;

    try {
      const payload: Employee = {
        ID: isEditing ? item.ID : 0,
        Name: form.Name!,
        Email: form.Email!,
        Roles: form.Roles!,
        LocationID: form.LocationID!,
      };

      let res;

      setIsProcessing(true);
      if (isEditing) {
        res = await CommonService.updateEmployee(payload);
        showSuccess(res.Message || "Employee updated successfully!");
      } else {
        res = await CommonService.createEmployee(payload);
        showSuccess(res.Message || "Employee created successfully!");
      }
      setIsProcessing(false);
      cancel();
      await CommonService.fetchEmployees();
    } catch (err) {
      showError(err);
    }
  };

  return (
    <CModel
      open
      size="md"
      onClose={cancel}
      title={isEditing ? "Update Employee" : "Create Employee"}
      actions={
        <Flex gap="3">
          <CButton
            label="Cancel"
            variant="soft"
            onClick={cancel}
            disabled={isProcessing}
          />
          <CButton
            label={isEditing ? "Update" : "Add"}
            variant="solid"
            onClick={handleSubmit}
            disabled={!form.Name}
            processing={isProcessing}
          />
        </Flex>
      }
    >
      <Flex direction="column" gap="4">
        <Box>
          <CInput
            label="Employee Name"
            value={form.Name}
            placeholder="Enter employee name"
            onChange={(val) => {
              handleChange("Name", val);
              setErrors((e) => ({ ...e, Name: "" }));
            }}
            error={errors.Name}
            fullWidth
          />
        </Box>

        <Box>
          <CInput
            label="Employee Email"
            value={form.Email}
            placeholder="Enter employee email"
            onChange={(val) => {
              handleChange("Email", val);
              setErrors((e) => ({ ...e, Email: "" }));
            }}
            error={errors.Email}
            fullWidth
          />
        </Box>

        <Box>
          <CSelect
            label="Employee Location"
            value={form.LocationID?.toString() ?? ""}
            placeholder="Select employee location"
            onChange={(val) => {
              handleChange("LocationID", Number(val));
              setErrors((e) => ({ ...e, LocationID: "" }));
            }}
            options={
              Locations?.map((loc) => ({
                value: loc.ID.toString(),
                label: `${loc.Address}, ${loc.City}, ${loc.Country} (${loc.Code})`,
              })) ?? []
            }
            error={errors.LocationID}
            fullWidth
          />
        </Box>

        <Box>
          <Text size="1" weight="medium">
            Employee Roles
          </Text>

          <CMultiSelect
            selected={form.Roles ?? []}
            options={Roles}
            placeholder="Select Roles"
            onChange={(vals) => {
              handleChange("Roles", vals);
              setErrors((e) => ({ ...e, Roles: "" }));
            }}
          />

          {errors.Roles && (
            <Text size="1" color="red" mt="1">
              {errors.Roles}
            </Text>
          )}
        </Box>
      </Flex>
    </CModel>
  );
};

export default ModifyEmployee;

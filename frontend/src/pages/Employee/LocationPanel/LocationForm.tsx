import React, { useState, useCallback, useEffect, useMemo } from "react";

import CModel from "../../../components/common/CModel";
import CInput from "../../../components/common/CInput";
import CButton from "../../../components/common/CButton";

import { Checkbox, Flex, Box, Text } from "@radix-ui/themes";
import { useErrorHandler } from "../../../hooks/useErrorHandler";
import CommonService from "../../../services/common/common.service";
import type { Location } from "../../../types/common";

interface LocationFormProps {
  close: () => void;
  location?: Location | null;
}

const LocationForm: React.FC<LocationFormProps> = ({ close, location }) => {
  const { showError, showSuccess } = useErrorHandler();
  const isEditMode = useMemo(() => !!location, [location]);

  const [form, setForm] = useState<Location>({
    ID: location?.ID || 0,
    Code: location?.Code || "",
    Address: location?.Address || "",
    City: location?.City || "",
    State: location?.State || "",
    Country: location?.Country || "",
    ZipCode: location?.ZipCode || "",
    IsPrimary: location?.IsPrimary || false,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    if (location) setForm(location);
  }, [location]);

  const handleChange = (key: keyof Location, value: any) => {
    setForm((prev) => ({ ...prev, [key]: value }));
  };

  const validate = useCallback(() => {
    console.log(form);
    const newErrors: Record<string, string> = {};
    if (!form.Code.trim()) newErrors.Code = "Location code is required.";
    if (!form.Address.trim()) newErrors.Address = "Address is required.";
    if (!form.City.trim()) newErrors.City = "City is required.";
    if (!form.State.trim()) newErrors.State = "State is required.";
    if (!form.Country.trim()) newErrors.Country = "Country is required.";
    if (!form.ZipCode.trim()) newErrors.ZipCode = "Zip code is required.";
    return newErrors;
  }, [form]);

  const handleSubmit = useCallback(async () => {
    const validationErrors = validate();
    setErrors(validationErrors);
    if (Object.keys(validationErrors).length > 0) return;

    try {
      if (isEditMode) {
        const res = await CommonService.updateLocation(form);
        showSuccess(res.Message || "Location updated successfully!");
      } else {
        const res = await CommonService.createLocation(form);
        showSuccess(res.Message || "Location created successfully!");
      }

      await CommonService.fetchLocations();
      close();
    } catch (err) {
      showError(err);
    }
  }, [form, validate, isEditMode, close, showError, showSuccess]);

  return (
    <CModel
      open={true}
      title={isEditMode ? "Update Location" : "Create Location"}
      onClose={close}
      size="md"
      actions={
        <Flex gap="3">
          <CButton label="Cancel" variant="soft" onClick={close} />
          <CButton
            label={isEditMode ? "Update" : "Create"}
            variant="solid"
            onClick={handleSubmit}
          />
        </Flex>
      }
    >
      <Box>
        <Flex direction="column" gap="3">
          <CInput
            label="Code"
            value={form.Code}
            placeholder="Enter location code"
            onChange={(val) => handleChange("Code", val)}
            error={errors.Code}
            fullWidth
          />

          <CInput
            label="Address"
            value={form.Address}
            placeholder="Enter address"
            onChange={(val) => handleChange("Address", val)}
            error={errors.Address}
            fullWidth
          />

          <CInput
            label="City"
            value={form.City}
            placeholder="Enter city"
            onChange={(val) => handleChange("City", val)}
            error={errors.City}
            fullWidth
          />

          <CInput
            label="State"
            value={form.State}
            placeholder="Enter state"
            onChange={(val) => handleChange("State", val)}
            error={errors.State}
            fullWidth
          />

          <CInput
            label="Country"
            value={form.Country}
            placeholder="Enter country"
            onChange={(val) => handleChange("Country", val)}
            error={errors.Country}
            fullWidth
          />

          <CInput
            label="Zip Code"
            value={form.ZipCode}
            placeholder="Enter zip code"
            onChange={(val) => handleChange("ZipCode", val)}
            error={errors.ZipCode}
            fullWidth
          />

          <Flex align="center" gap="2" mt="2">
            <Checkbox
              checked={form.IsPrimary}
              onCheckedChange={(v) => handleChange("IsPrimary", v === true)}
            />
            <Text size="2">Mark as primary</Text>
          </Flex>
        </Flex>
      </Box>
    </CModel>
  );
};

export default LocationForm;

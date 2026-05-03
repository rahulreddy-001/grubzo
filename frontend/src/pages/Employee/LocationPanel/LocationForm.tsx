import React from "react";

import CModel from "../../../components/common/CModel";
import CInput from "../../../components/common/CInput";
import CButton from "../../../components/common/CButton";

import { Checkbox, Flex, Box, Text } from "@radix-ui/themes";
import { useErrorHandler } from "../../../hooks/useErrorHandler";
import CommonService from "../../../services/common/common.service";
import { LocationSchema, type Location } from "../../../types/common.d";
import { useFormik } from "formik";
import { zodToFormik } from "../../../utils/zodToFormik";

interface LocationFormProps {
  close: () => void;
  location?: Location | null;
}

const LocationForm: React.FC<LocationFormProps> = ({ close, location }) => {
  const isUpdate = !!location;
  const { showError, showSuccess } = useErrorHandler();

  const form = useFormik<Location>({
    initialValues: {
      ID: location?.ID || 0,
      Code: location?.Code || "",
      Address: location?.Address || "",
      City: location?.City || "",
      State: location?.State || "",
      Country: location?.Country || "",
      ZipCode: location?.ZipCode || "",
      IsPrimary: location?.IsPrimary || false,
    },
    validate: zodToFormik(LocationSchema),
    onSubmit: async (loc: Location) => {
      try {
        let response;
        if (isUpdate) {
          response = await CommonService.updateLocation(loc);
        } else {
          response = await CommonService.createLocation(loc);
        }
        showSuccess(response.Message);
        close();
        CommonService.fetchLocations();
      } catch (e) {
        showError(e);
      }
    },
  });

  return (
    <CModel
      open={true}
      title={isUpdate ? "Update Location" : "Create Location"}
      onClose={close}
      size="md"
      actions={
        <Flex gap="3">
          <CButton label="Cancel" onClick={close} />
          <CButton
            label={isUpdate ? "Update" : "Create"}
            onClick={form.submitForm}
          />
        </Flex>
      }
    >
      <form onSubmit={form.handleSubmit}>
        <Box>
          <Flex direction="column" gap="3">
            <CInput
              label="Code"
              placeholder="Enter location code"
              value={form.values.Code}
              onChange={(val) => form.setFieldValue("Code", val)}
              error={form.errors.Code}
              fullWidth
            />

            <CInput
              label="Address"
              placeholder="Enter address"
              value={form.values.Address}
              onChange={(val) => form.setFieldValue("Address", val)}
              error={form.errors.Address}
              fullWidth
            />

            <CInput
              label="City"
              placeholder="Enter city"
              value={form.values.City}
              onChange={(val) => form.setFieldValue("City", val)}
              error={form.errors.City}
              fullWidth
            />

            <CInput
              label="State"
              placeholder="Enter state"
              value={form.values.State}
              onChange={(val) => form.setFieldValue("State", val)}
              error={form.errors.State}
              fullWidth
            />

            <CInput
              label="Country"
              placeholder="Enter country"
              value={form.values.Country}
              onChange={(val) => form.setFieldValue("Country", val)}
              error={form.errors.Country}
              fullWidth
            />

            <CInput
              label="Zip Code"
              placeholder="Enter zip code"
              value={form.values.ZipCode}
              onChange={(val) => form.setFieldValue("ZipCode", val)}
              error={form.errors.ZipCode}
              fullWidth
            />

            <Flex align="center" gap="2" mt="2">
              <Checkbox
                checked={form.values.IsPrimary}
                onCheckedChange={(v) => form.setFieldValue("IsPrimary", v)}
              />
              <Text size="2">Mark as primary</Text>
            </Flex>
          </Flex>
        </Box>
      </form>
    </CModel>
  );
};

export default LocationForm;

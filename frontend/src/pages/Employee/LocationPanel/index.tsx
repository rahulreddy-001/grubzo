import React, { useState, useEffect } from "react";
import { useSelector } from "react-redux";
import type { RootState } from "../../../services/store";
import CTable from "../../../components/common/CTable";
import CButton from "../../../components/common/CButton";
import { Edit, Plus } from "lucide-react";
import CommonService from "../../../services/common/common.service";
import LocationForm from "./LocationForm";
import { Box, Flex, Text, IconButton } from "@radix-ui/themes";

const LocationPanel: React.FC = () => {
  const { Locations, isLoading } = useSelector((s: RootState) => s.common);
  useEffect(() => {
    CommonService.fetchLocations();
  }, []);

  const [modelState, setModelState] = useState<{
    isOpen: boolean;
    selectedLocation: any | null;
  }>({ isOpen: false, selectedLocation: null });

  const handleOpen = (row: any) => {
    setModelState({
      isOpen: true,
      selectedLocation: row,
    });
  };

  const handleClose = () => {
    setModelState({
      isOpen: false,
      selectedLocation: null,
    });
  };

  return (
    <Box>
      <CTable
        title="Locations"
        data={Locations ?? []}
        rowKey="ID"
        loading={isLoading}
        onRefresh={() => CommonService.fetchLocations()}
        actions={
          <CButton
            label="Create Location"
            startIcon={<Plus size={16} />}
            onClick={handleOpen}
          />
        }
        columns={[
          {
            key: "Code",
            label: "Location Code",
            render: (row) => (
              <Text style={{ paddingLeft: 12 }}>{row.Code}</Text>
            ),
          },
          {
            key: "Address",
            label: "Address",
            render: (row) => (
              <Text>{`${row.Address}, ${row.City}, ${row.State}, ${row.Country}`}</Text>
            ),
          },
          {
            key: "ZipCode",
            label: "Zip Code",
            render: (row) => <Text>{row.ZipCode}</Text>,
          },
          {
            key: "Actions",
            label: "Actions",
            width: 80,
            render: (row) => (
              <Flex gap="2" style={{ marginLeft: "10px", cursor: "pointer" }}>
                <IconButton
                  variant="ghost"
                  radius="full"
                  style={{ cursor: "pointer", padding: "5px", margin: "1px" }}
                  onClick={() => handleOpen(row)}
                >
                  <Edit size={16} />
                </IconButton>
              </Flex>
            ),
          },
        ]}
      />

      {modelState.isOpen && (
        <LocationForm
          close={() => handleClose()}
          location={modelState.selectedLocation}
        />
      )}
    </Box>
  );
};

export default LocationPanel;

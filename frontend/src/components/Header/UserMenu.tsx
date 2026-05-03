import AuthService from "../../services/auth/auth.service";
import { useSelector } from "react-redux";
import type { RootState } from "../../services/store";
import {
  Avatar,
  Box,
  Button,
  DropdownMenu,
  Flex,
  Text,
} from "@radix-ui/themes";
import CommonService from "../../services/common/common.service";
import { MapPin } from "lucide-react";
import { useErrorHandler } from "../../hooks/useErrorHandler";

const UserMenu = () => {
  const user = useSelector((state: RootState) => state.auth.user);
  const { showError } = useErrorHandler();
  const locations = useSelector((state: RootState) => state.common.Locations);
  if (locations?.length == 0) {
    CommonService.fetchLocations();
  }

  const handleLogout = async () => {
    try {
      await AuthService.logout();
      window.location.href = "/";
    } catch (err) {
      console.error(err);
    }
  };

  const handleLocationChange = async (locaitonID: number) => {
    try {
      await CommonService.setUserLocation(locaitonID);
      location.reload();
    } catch (err) {
      showError(err);
    }
  };

  return (
    <Flex align={"center"}>
      <DropdownMenu.Root>
        <DropdownMenu.Trigger>
          <Button
            style={{
              cursor: "pointer",
              fontSize: "15px",
              padding: "25px",
              background: "transparent",
              color: "black",
              borderRadius: "0",
              outline: "none",
            }}
          >
            <Flex gap="3" align={"center"}>
              <Box
                style={{ outline: "none", border: "none", maxWidth: "250px" }}
              >
                <Flex
                  align="center"
                  style={{ color: "grey", fontSize: "13px", cursor: "pointer" }}
                >
                  <Flex direction={"column"} align={"center"}>
                    <Flex align={"center"} gap={"1"}>
                      <MapPin size={16} />
                      <Text size={"1"} weight={"bold"}>
                        {user?.Location.City} ({user?.Location.Code})
                      </Text>
                    </Flex>
                    <Text>{user?.Location.Address}</Text>
                  </Flex>
                </Flex>
              </Box>
              <Avatar
                color="gray"
                radius="full"
                size="3"
                fallback={
                  user?.Name
                    ? user.Name.split(" ")
                        .map((w) => w[0])
                        .join("")
                        .toUpperCase()
                    : ""
                }
              />
              <Text style={{ color: "gray", fontSize: "13px" }}>
                {user?.Name}
              </Text>
              <DropdownMenu.TriggerIcon
                style={{ color: "gray", fontWeight: "bold" }}
              />
            </Flex>
          </Button>
        </DropdownMenu.Trigger>
        <DropdownMenu.Content align="end" sideOffset={8}>
          <DropdownMenu.Item onClick={handleLogout}>Log out</DropdownMenu.Item>
          <DropdownMenu.Sub>
            {CommonService.locationChangeAccess() && (
              <DropdownMenu.SubTrigger>Change Location</DropdownMenu.SubTrigger>
            )}
            <DropdownMenu.SubContent>
              {locations?.map((location) => (
                <DropdownMenu.Item
                  key={location.ID}
                  onClick={() => handleLocationChange(location.ID)}
                >
                  {`${location.Address}, ${location.City} (${location.Code})`}
                </DropdownMenu.Item>
              ))}
            </DropdownMenu.SubContent>
          </DropdownMenu.Sub>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </Flex>
  );
};

export default UserMenu;

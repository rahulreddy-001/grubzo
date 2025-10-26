import { Dropdown, MenuButton, Menu, MenuItem } from "@mui/joy";
import Stack from "@mui/joy/Stack";
import Avatar from "@mui/joy/Avatar";
import Typography from "@mui/joy/Typography";
import ArrowDropDown from "@mui/icons-material/ArrowDropDown";
import AuthService from "../../services/auth/auth.service";
import { useNavigate } from "react-router";
import { useAuth } from "../../context/AuthProvider";

const plainButtonSX = {
  bgcolor: "transparent",
  color: "grey",
  "&:hover": { bgcolor: "transparent", color: "black" },
  "&:focus-visible": { outline: "none", boxShadow: "none" },
  "--Input-focusedThickness": "0px",
  border: "none",
};

const UserMenu = ({ name }: { name: any }) => {
  let navigate = useNavigate();
  let { refreshUser } = useAuth();
  const handleLogout = async () => {
    try {
      await AuthService.logout();
      await refreshUser();
      await navigate("/");
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <Dropdown>
      <MenuButton endDecorator={<ArrowDropDown />} sx={plainButtonSX}>
        <Stack direction="row" alignItems="center" gap="10px">
          <Avatar>{name[0]}</Avatar>
          <Typography sx={{ fontSize: "15px" }}>{name}</Typography>
        </Stack>
      </MenuButton>
      <Menu sx={{ py: 0 }}>
        <MenuItem
          sx={{ fontSize: "15px", height: "30px" }}
          onClick={handleLogout}
        >
          Log out
        </MenuItem>
      </Menu>
    </Dropdown>
  );
};

export default UserMenu;

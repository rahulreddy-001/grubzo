import grubzoIcon from "../../assets/grubzo_logo_r.png";
import Search from "./Search";
import LoginSignUp from "./LoginSignup";
import UserMenu from "./UserMenu";
import { useSelector } from "react-redux";
import type { RootState } from "../../services/store";
import { Box, Flex } from "@radix-ui/themes";
import { ListDivider } from "@mui/joy";

const Header: React.FC = () => {
  const user = useSelector((state: RootState) => state.auth.user);
  return (
    <Box>
      <Flex
        direction={"row"}
        justify={"between"}
        align={"center"}
        style={{ padding: "0 100px" }}
      >
        <img src={grubzoIcon} height="60px" alt="Logo" />
        {user?.Type != "employee" ? <Search /> : null}
        {user ? <UserMenu /> : <LoginSignUp />}
      </Flex>
      <ListDivider />
    </Box>
  );
};

export default Header;

import React, { useEffect } from "react";
import Stack from "@mui/joy/Stack";
import Divider from "@mui/joy/Divider";
import grubzoIcon from "../../assets/grubzo_logo_r.png";
import Search from "./Search";
import LoginSignUp from "./LoginSignup";
import UserMenu from "./UserMenu";
import { useSelector } from "react-redux";
import type { RootState } from "../../services/store";
import AuthService from "../../services/auth/auth.service";

const Header: React.FC = () => {
  const user = useSelector((state: RootState) => state.auth.user);
  useEffect(() => {
    const initUser = async (): Promise<void> => {
      const authenticated = await AuthService.isAuthenticated();
      if (authenticated && !user) {
        try {
          await AuthService.fetchUser();
        } catch (err) {
          console.error("Failed to fetch user:", err);
        }
      }
    };

    void initUser();
  }, [user]);

  return (
    <>
      <Stack
        direction="row"
        alignItems="center"
        justifyContent="space-between"
        paddingX="10vw"
      >
        <img src={grubzoIcon} height="60px" alt="Logo" />
        {user ? <Search /> : null}
        {user ? <UserMenu name={user.Name || "[User]"} /> : <LoginSignUp />}
      </Stack>
      <Divider />
    </>
  );
};

export default Header;

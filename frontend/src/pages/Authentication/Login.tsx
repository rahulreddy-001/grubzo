import React from "react";
import type { FormEvent } from "react";
import { Box, Typography, Input, Button, Link, Sheet, Stack } from "@mui/joy";
import { useNavigate, useSearchParams } from "react-router";
import grubzoIcon from "../../assets/grubzo_logo_t.png";
import bgImage from "../../assets/login_signup_bg.jpg";

import EmailIcon from "@mui/icons-material/Email";
import LockIcon from "@mui/icons-material/Lock";
import GoogleIcon from "@mui/icons-material/Google";
import GitHubIcon from "@mui/icons-material/GitHub";

import { UserTypes } from "../../types/constants.ts";
import authService from "../../services/auth/auth.service.ts";
import { useErrorHandler } from "../../hooks/useErrorHandler.tsx";

const Login: React.FC = () => {
  const navigate = useNavigate();
  const { showError } = useErrorHandler();
  const [searchParams] = useSearchParams();

  const userType = UserTypes[searchParams.get("t") ?? "usr"] ?? UserTypes.usr;
  authService.isAuthenticated().then((auth) => {
    if (auth) {
      navigate("/");
    }
  });
  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const email = (
      e.currentTarget.elements.namedItem("email") as HTMLInputElement
    )?.value.trim();
    const password = (
      e.currentTarget.elements.namedItem("password") as HTMLInputElement
    )?.value;

    if (!email || !password) {
      showError("Please fill in all fields.");
      return;
    }

    authService
      .login({
        Type: userType,
        Email: email,
        Password: password,
      })
      .then(() => {
        window.location.pathname = "/";
      })
      .catch((err) => {
        if (err?.Error) showError(err);
        else showError("Login failed. Please try again.");
        console.log(err);
      });
  };

  const handleOAuth = (provider: string) => {
    window.location.pathname = `auth/v1/oauth/login/${provider.toLowerCase()}`;
  };

  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        alignItems: "flex-end",
        justifyContent: "center",
        minHeight: "100vh",
        width: "100%",
        margin: 0,
        padding: 0,
        backgroundColor: "#f9f9f9",
        fontFamily: "'Inter', sans-serif",
        backgroundImage: `url(${bgImage})`,
        backgroundSize: "cover",
        backgroundPosition: "center",
        backgroundRepeat: "no-repeat",
      }}
    >
      <Sheet
        variant="soft"
        sx={{
          width: "100%",
          maxWidth: 360,
          height: 500,
          borderRadius: 3,
          p: 4,
          margin: 10,
          boxShadow: "sm",
          backgroundColor: "rgba(255, 255, 255, 0.8)",
        }}
      >
        <Box sx={{ mb: 3, mt: 8, display: "flex", justifyContent: "center" }}>
          <img
            src={grubzoIcon}
            alt="Grubzo Logo"
            style={{ width: 100, height: "auto" }}
          />
        </Box>

        <form onSubmit={handleSubmit}>
          <Stack spacing={1.5}>
            <Input
              name="email"
              placeholder="Email"
              type="email"
              required
              size="md"
              startDecorator={<EmailIcon />}
              variant="outlined"
              sx={{ borderRadius: 3, padding: 0.7 }}
            />
            <Input
              name="password"
              placeholder="Password"
              type="password"
              required
              size="md"
              startDecorator={<LockIcon />}
              variant="outlined"
              sx={{ borderRadius: 3, padding: 0.7 }}
            />

            <Button
              type="submit"
              variant="solid"
              color="primary"
              size="md"
              sx={{ borderRadius: 3, padding: 0.7 }}
            >
              Login
            </Button>
          </Stack>
        </form>

        <Typography
          sx={{
            mt: 2,
            mb: 1,
            textAlign: "center",
            color: "text.tertiary",
            fontSize: "0.875rem",
          }}
        >
          Or continue with
        </Typography>

        <Stack spacing={1} direction="row" sx={{ justifyContent: "center" }}>
          <Button
            variant="outlined"
            color="neutral"
            startDecorator={<GoogleIcon />}
            onClick={() => handleOAuth("Google")}
            sx={{ borderRadius: 3, flex: 1, padding: 0.7 }}
          >
            Google
          </Button>
          <Button
            variant="outlined"
            color="neutral"
            startDecorator={<GitHubIcon />}
            onClick={() => handleOAuth("GitHub")}
            sx={{ borderRadius: 3, flex: 1, padding: 0.7 }}
          >
            GitHub
          </Button>
        </Stack>

        <Typography
          sx={{
            mt: 3,
            textAlign: "center",
            color: "text.tertiary",
            fontSize: "0.875rem",
          }}
        >
          Don't have an account?{" "}
          <Link
            component="button"
            onClick={() => navigate(`/signup`)}
            underline="hover"
            sx={{ fontSize: "0.875rem" }}
          >
            Sign Up
          </Link>
        </Typography>
      </Sheet>
    </Box>
  );
};

export default Login;

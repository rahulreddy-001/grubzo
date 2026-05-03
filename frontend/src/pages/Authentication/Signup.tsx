import React, { useState, type FormEvent } from "react";
import { Box, Typography, Input, Button, Link, Sheet, Stack } from "@mui/joy";
import { useNavigate } from "react-router";
import grubzoIcon from "../../assets/grubzo_logo_t.png";
import bgImage from "../../assets/login_signup_bg.jpg";

import PersonIcon from "@mui/icons-material/Person";
import EmailIcon from "@mui/icons-material/Email";
import LockIcon from "@mui/icons-material/Lock";

import GoogleIcon from "@mui/icons-material/Google";
import GitHubIcon from "@mui/icons-material/GitHub";
import authService from "../../services/auth/auth.service";
import { useErrorHandler } from "../../hooks/useErrorHandler";

const Signup: React.FC = () => {
  const { showError } = useErrorHandler();
  const navigate = useNavigate();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!name.trim()) {
      showError("Name is required");
      return;
    }
    if (!email.trim()) {
      showError("Email is required");
      return;
    }
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      showError("Invalid email format");
      return;
    }
    if (!password || password.length < 6) {
      showError("Password must be at least 6 characters");
      return;
    }
    if (password !== confirmPassword) {
      showError("Passwords do not match");
      return;
    }

    authService
      .signup({ email, password, name })
      .then(() => {
        navigate("/login");
      })
      .catch((err) => {
        if (err.Error) showError(err);
        else showError("Signup failed. Please try again.");
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
        <Box sx={{ mb: 3, mt: 3, display: "flex", justifyContent: "center" }}>
          <img
            src={grubzoIcon}
            alt="Grubzo Logo"
            style={{ width: 100, height: "auto" }}
          />
        </Box>

        <form onSubmit={handleSubmit}>
          <Stack spacing={1.5}>
            <Input
              placeholder="Name"
              required
              size="md"
              startDecorator={<PersonIcon />}
              variant="outlined"
              value={name}
              onChange={(e) => setName(e.target.value)}
              sx={{ borderRadius: 3, padding: 0.7 }}
            />
            <Input
              placeholder="Email"
              type="email"
              required
              size="md"
              startDecorator={<EmailIcon />}
              variant="outlined"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              sx={{ borderRadius: 3, padding: 0.7 }}
            />
            <Input
              placeholder="Password"
              type="password"
              required
              size="md"
              startDecorator={<LockIcon />}
              variant="outlined"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              sx={{ borderRadius: 3, padding: 0.7 }}
            />
            <Input
              placeholder="Confirm Password"
              type="password"
              required
              size="md"
              startDecorator={<LockIcon />}
              variant="outlined"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              sx={{ borderRadius: 3, padding: 0.7 }}
            />

            <Button
              type="submit"
              variant="solid"
              color="primary"
              size="md"
              sx={{ borderRadius: 3, padding: 0.7 }}
            >
              Sign Up
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
          Already have an account?{" "}
          <Link
            component="button"
            onClick={() => navigate("/login")}
            underline="hover"
            sx={{ fontSize: "0.875rem" }}
          >
            Login
          </Link>
        </Typography>
      </Sheet>
    </Box>
  );
};

export default Signup;

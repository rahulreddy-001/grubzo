import React from "react";
import { Box, Typography, Button, Stack } from "@mui/joy";
import { useNavigate } from "react-router";
import grubzoIcon from "../assets/grubzo_logo_t.png";
import bgImage from "../assets/login_signup_bg.jpg";

const NotFound: React.FC = () => {
  const navigate = useNavigate();

  return (
    <Box
      sx={{
        height: "100vh",
        overflow: "hidden",
        width: "100%",
        backgroundImage: `url(${bgImage})`,
        backgroundSize: "cover",
        backgroundPosition: "center",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <Stack
        spacing={2}
        alignItems="center"
        justifyContent="center"
        sx={{
          background: "rgba(255,255,255,0.8)",
          backdropFilter: "blur(6px)",
          borderRadius: "2xl",
          p: 6,
          boxShadow: "lg",
          textAlign: "center",
          maxWidth: 420,
        }}
      >
        <img
          src={grubzoIcon}
          alt="Grubzo Logo"
          style={{ height: 70, objectFit: "contain", marginBottom: 8 }}
        />

        <Typography
          level="h1"
          sx={{ fontSize: "5rem", fontWeight: "bold", color: "#E53935" }}
        >
          404
        </Typography>

        <Typography level="title-lg" sx={{ color: "neutral.800" }}>
          Page Not Found
        </Typography>

        <Typography level="body-md" sx={{ color: "neutral.600" }}>
          Oops! The page you're looking for doesn’t exist or was moved.
        </Typography>

        <Stack direction="row" spacing={2} sx={{ mt: 2 }}>
          <Button
            color="danger"
            variant="solid"
            onClick={() => navigate("/")}
            sx={{ px: 3, borderRadius: 3 }}
          >
            Go Home
          </Button>
          <Button
            variant="outlined"
            color="neutral"
            onClick={() => navigate(-1)}
            sx={{ px: 3, borderRadius: 3 }}
          >
            Go Back
          </Button>
        </Stack>
      </Stack>
    </Box>
  );
};

export default NotFound;

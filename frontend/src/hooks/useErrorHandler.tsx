import { useNotification } from "../context/NotificationProvider";
import type { ErrorResponse } from "../types/common";
import { Snackbar, IconButton, Typography } from "@mui/joy";
import React, { useEffect } from "react";
import CloseIcon from "@mui/icons-material/Close";

function isErrorResponse(error: unknown): error is ErrorResponse {
  return (
    typeof error === "object" &&
    error !== null &&
    "Error" in error &&
    typeof (error as any).Error === "string"
  );
}

interface JoySnackbarProps {
  message: string;
  color?: "success" | "danger" | "neutral";
  visible?: boolean;
  setVisible?: (open: boolean) => void;
  duration?: number;
}

export const JoySnackbar: React.FC<JoySnackbarProps> = ({
  message,
  color = "success",
  visible = true,
  setVisible,
  duration = 3000,
}) => {
  useEffect(() => {
    if (visible) {
      const timer = setTimeout(() => setVisible?.(false), duration);
      return () => clearTimeout(timer);
    }
  }, [visible, duration, setVisible]);

  return (
    <Snackbar
      open={visible}
      onClose={() => setVisible?.(false)}
      endDecorator={
        <IconButton
          size="sm"
          variant="plain"
          color="neutral"
          onClick={() => setVisible?.(false)}
          sx={{
            borderRadius: "50%",
            transition: "background 0.2s",
            "&:hover": { backgroundColor: "rgba(0,0,0,0.08)" },
          }}
        >
          <CloseIcon fontSize="small" />
        </IconButton>
      }
      anchorOrigin={{ vertical: "top", horizontal: "right" }}
      variant="soft"
      color={color}
      sx={{ minWidth: 200, maxWidth: 320, py: 0.5, px: 1 }}
    >
      <Typography level="body-sm">{message}</Typography>
    </Snackbar>
  );
};

export const useErrorHandler = () => {
  const { showNotification } = useNotification();

  const showError = (error: unknown) => {
    let message = "Something went wrong";
    if (typeof error === "string") message = error;
    if (isErrorResponse(error)) message = error.Error;
    if (message) {
      message = message.charAt(0).toLocaleUpperCase() + message.slice(1);
    }
    showNotification(<JoySnackbar message={message} color="danger" />);
  };

  const showSuccess = (message: string) => {
    showNotification(<JoySnackbar message={message} color="success" />);
  };

  return { showError, showSuccess };
};

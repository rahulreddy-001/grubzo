import { Button, Spinner, type ButtonProps } from "@radix-ui/themes";
import React from "react";

type RadixColor = ButtonProps["color"];
type CButtonColor = RadixColor | "primary";

interface CButtonProps extends Omit<ButtonProps, "color"> {
  label?: string;
  color?: CButtonColor;
  fullWidth?: boolean;
  startIcon?: React.ReactNode;
  endIcon?: React.ReactNode;
  styles?: React.CSSProperties;
  processing?: boolean;
}

const CButton = React.forwardRef<HTMLButtonElement, CButtonProps>(
  (
    {
      label,
      onClick,
      type = "button",
      variant = "solid",
      color = "primary",
      fullWidth = false,
      startIcon,
      endIcon,
      disabled = false,
      processing = false,
      styles = {},
      ...rest
    },
    ref
  ) => {
    return (
      <Button
        ref={ref}
        onClick={onClick}
        type={type}
        variant={variant}
        color={color as RadixColor}
        disabled={disabled || processing}
        style={{
          width: fullWidth ? "100%" : undefined,
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          gap: "6px",
          borderRadius: "3px",
          cursor: "pointer",
          ...styles,
        }}
        {...rest}
      >
        <Spinner loading={processing} />
        {startIcon}
        {label}
        {endIcon}
      </Button>
    );
  }
);
export default CButton;

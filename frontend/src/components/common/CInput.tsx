import React from "react";
import { TextField, Flex, Text } from "@radix-ui/themes";

interface CInputProps {
  value?: string | number;
  onChange: (value: string) => void;
  placeholder?: string;
  startIcon?: React.ReactNode;
  endIcon?: React.ReactNode;
  fullWidth?: boolean;
  disabled?: boolean;
  label?: string;
  error?: string;
  name?: string;
  type?:
    | "date"
    | "datetime-local"
    | "email"
    | "hidden"
    | "month"
    | "number"
    | "password"
    | "search"
    | "tel"
    | "text"
    | "time"
    | "url"
    | "week";
}

const CInput: React.FC<CInputProps> = ({
  value = "",
  onChange,
  placeholder,
  startIcon,
  endIcon,
  fullWidth = false,
  disabled = false,
  label,
  error,
  name,
  type = "text",
}) => {
  return (
    <Flex direction="column" gap="1" width={fullWidth ? "100%" : "auto"}>
      {label && (
        <Text size="1" weight="medium">
          {label}
        </Text>
      )}

      <TextField.Root
        name={name}
        type={type}
        size="2"
        radius="medium"
        disabled={disabled}
        placeholder={placeholder}
        style={{
          minHeight: "33px",
          width: fullWidth ? "100%" : "auto",
          border: error ? "1px solid var(--red-9)" : "1px solid var(--gray-6)",
          borderRadius: "3px",
          boxShadow: error ? "0 0 0 1px var(--red-4)" : "none",
          paddingLeft: startIcon ? "8px" : undefined,
          paddingRight: endIcon ? "8px" : undefined,
        }}
        onChange={(e: any) => onChange(e.target.value)}
        value={value}
      >
        {startIcon && <TextField.Slot>{startIcon}</TextField.Slot>}

        {endIcon && <TextField.Slot>{endIcon}</TextField.Slot>}
      </TextField.Root>
      {error && (
        <Text size="1" color="red">
          {error}
        </Text>
      )}
    </Flex>
  );
};

export default CInput;

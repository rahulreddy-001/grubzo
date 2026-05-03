import React from "react";
import { Select, Text, Flex } from "@radix-ui/themes";

interface Option {
  value: any;
  label: string;
}

interface CSelectProps {
  value: string;
  options: Option[];
  placeholder?: string;
  onChange: (value: string) => void;
  label?: string;
  error?: string;
  fullWidth?: boolean;
}

const CSelect: React.FC<CSelectProps> = ({
  value,
  options,
  placeholder = "Select...",
  onChange,
  label,
  error,
  fullWidth = true,
}) => {
  return (
    <Flex
      direction="column"
      gap="1"
      style={{ width: fullWidth ? "100%" : "auto" }}
    >
      {label && (
        <Text size="1" weight="medium">
          {label}
        </Text>
      )}

      <Select.Root value={value} onValueChange={onChange} size="1">
        <Select.Trigger
          radius="small"
          placeholder={placeholder}
          style={{
            height: "34px",
            width: "100%",
            border: error
              ? "1px solid var(--red-9)"
              : "1px solid var(--gray-6)",
            borderRadius: "4px",
            cursor: "pointer",
          }}
        />

        <Select.Content position="popper">
          {options.map((opt) => (
            <Select.Item
              key={opt.value}
              value={String(opt.value)}
              style={{
                padding: "5px",
                paddingLeft: "24px",
                cursor: "pointer",
                height: "28px",
              }}
            >
              {opt.label}
            </Select.Item>
          ))}
        </Select.Content>
      </Select.Root>

      {error && (
        <Text size="1" color="red">
          {error}
        </Text>
      )}
    </Flex>
  );
};

export default CSelect;

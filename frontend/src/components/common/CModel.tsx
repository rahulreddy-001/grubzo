import React from "react";
import * as Dialog from "@radix-ui/react-dialog";
import { Flex, Box, Text, IconButton, VisuallyHidden } from "@radix-ui/themes";
import { X } from "lucide-react";
import "./style.css";

interface CModelProps {
  open: boolean;
  onClose: () => void;
  title?: string;
  size?: "sm" | "md" | "lg" | "full";
  anchor?: "right" | "left" | "bottom" | "top";
  children?: React.ReactNode;
  actions?: React.ReactNode;
  closeOnBackdrop?: boolean;
  stickyFooter?: boolean;
}

const CModel: React.FC<CModelProps> = ({
  open,
  onClose,
  title,
  size = "md",
  anchor = "right",
  children,
  actions,
  closeOnBackdrop = true,
  stickyFooter = false,
}) => {
  return (
    <Dialog.Root
      open={open}
      onOpenChange={(v) => !v && onClose()}
      modal={false}
    >
      <Dialog.Overlay
        className="cmodel-overlay"
        onClick={() => closeOnBackdrop && onClose()}
      />

      <Dialog.Content
        asChild
        onClick={(e) => e.stopPropagation()}
        className="no-default"
        style={{ pointerEvents: "auto" }}
      >
        <div className={`cmodel-content anchor-${anchor} size-${size}`}>
          <Flex justify="between" align="center" className="cmodel-header">
            {title ? (
              <Dialog.Title asChild>
                <Text size="4" weight="bold">
                  {title}
                </Text>
              </Dialog.Title>
            ) : (
              <Dialog.Title>
                <VisuallyHidden>Modal</VisuallyHidden>
              </Dialog.Title>
            )}

            <Dialog.Close asChild>
              <IconButton
                variant="ghost"
                radius="full"
                onClick={onClose}
                className="cmodel-close"
              >
                <X size={20} />
              </IconButton>
            </Dialog.Close>
          </Flex>

          <Box className="cmodel-body">{children}</Box>
          {actions && (
            <Box className={`cmodel-footer ${stickyFooter ? "sticky" : ""}`}>
              <Flex justify="end" gap="3">
                {actions}
              </Flex>
            </Box>
          )}
        </div>
      </Dialog.Content>
    </Dialog.Root>
  );
};

export default CModel;

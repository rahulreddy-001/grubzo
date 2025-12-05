// LocationStyles.ts — combined styled-components for LocationPanel & LocationForm
import styled from "styled-components";

// ========================= LOCATION PANEL STYLES =========================

export const PanelWrapper = styled.div`
  padding: 0;
`;

export const Row = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

export const ActionCell = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;

export const IconButtonStyled = styled.button`
  padding: 6px;
  border-radius: 8px;
  background: #f4f4f4;
  border: 1px solid #ddd;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: 0.2s ease;

  &:hover {
    background: #e8e8e8;
  }
`;

export const TextBold = styled.p`
  font-weight: 600;
  font-size: 0.9rem;
  margin: 0;
  text-transform: capitalize;
`;

// ========================= LOCATION FORM STYLES =========================

export const FormStack = styled.div`
  display: flex;
  flex-direction: column;
  gap: 22px;
`;

export const FieldWrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 6px;
`;

export const Label = styled.label`
  font-weight: 600;
  font-size: 0.9rem;
  color: #222;
`;

export const TextInput = styled.input`
  padding: 10px 12px;
  border: 1px solid #ccc;
  border-radius: 8px;
  font-size: 0.9rem;

  &:focus {
    border-color: #1976d2;
    outline: none;
  }
`;

export const ErrorText = styled.span`
  color: #d32f2f;
  font-size: 0.78rem;
`;

export const CheckboxRow = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;

export const FlexRow = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 12px;
`;

export const StyledButton = styled.button<{ $variant?: "outline" | "solid" }>`
  padding: 8px 18px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;

  ${(p) =>
    p.$variant === "outline"
      ? `background: transparent; border: 1px solid #aaa; color: #444;`
      : `background: #1976d2; border: none; color: #fff;`};

  opacity: ${(p) => (p.disabled ? 0.6 : 1)};
`;

export default {
  PanelWrapper,
  Row,
  ActionCell,
  IconButtonStyled,
  TextBold,
  FormStack,
  FieldWrapper,
  Label,
  TextInput,
  ErrorText,
  CheckboxRow,
  FlexRow,
  StyledButton,
};

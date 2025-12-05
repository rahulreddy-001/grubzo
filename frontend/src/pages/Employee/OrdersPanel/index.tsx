import React from "react";
import { PanelContainer, Title, BodyText, InfoBox } from "./Style";

const OrdersPanel: React.FC = () => {
  return (
    <PanelContainer>
      <Title>Orders Management</Title>
      <BodyText>
        Track order flow: Ordered → Confirmed → Cooking → Done → Payment.
      </BodyText>

      <InfoBox>(Orders list and workflow view goes here)</InfoBox>
    </PanelContainer>
  );
};

export default OrdersPanel;

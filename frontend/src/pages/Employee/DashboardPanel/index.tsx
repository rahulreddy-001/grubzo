import React from "react";
import { PanelContainer, Title, BodyText, StatsBox } from "./Style.ts";

const DashboardPanel: React.FC = () => {
  return (
    <PanelContainer>
      <Title>Dashboard</Title>
      <BodyText>
        Overview of cafeteria performance, daily summary, and quick stats.
      </BodyText>

      <StatsBox>(Dashboard stats and summary go here)</StatsBox>
    </PanelContainer>
  );
};

export default DashboardPanel;

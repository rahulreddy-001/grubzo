import styled from "styled-components";

export const PanelContainer = styled.div`
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

export const Title = styled.h2`
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #222;
`;

export const BodyText = styled.p`
  margin: 0;
  font-size: 0.9rem;
  color: #666;
`;

export const StatsBox = styled.div`
  margin-top: 12px;
  padding: 16px;
  border-radius: 8px;
  background: #f2f2f2;
  border: 1px solid #ddd;
  color: #555;
  text-align: center;
  font-size: 0.9rem;
`;

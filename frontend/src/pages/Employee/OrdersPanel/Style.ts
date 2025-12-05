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
  color: #555;
`;

export const InfoBox = styled.div`
  margin-top: 12px;
  padding: 16px;
  border-radius: 8px;
  background: #f4f4f4;
  border: 1px solid #ddd;
  text-align: center;
  color: #666;
  font-size: 0.9rem;
`;

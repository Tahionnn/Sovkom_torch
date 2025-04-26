import { Button, Card, Container } from '@mui/material';
import styled from 'styled-components';
const AuthContainer = styled(Container)`
  height: 100vh;
  width: 100%;
  display: flex;
  align-items: center;
`;

const AuthCard = styled(Card)`
  width: 100vw;
  display: flex;
  flex-direction: column;
  gap: 1rem;
`;

const AuthButton = styled(Button)``;
const AuthTitle = styled(Button)``;

const AuthLogoContainer = styled.div`
  display: flex;
  width: 100%;
`;

const AuthLogoImg = styled.img`
  width: 200px;
  margin: auto;
`;
export {
  AuthCard,
  AuthContainer,
  AuthButton,
  AuthTitle,
  AuthLogoContainer,
  AuthLogoImg,
};

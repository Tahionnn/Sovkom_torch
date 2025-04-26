import { Button, Card, Container, styled } from "@mui/material";

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

export {AuthCard, AuthContainer, AuthButton, AuthTitle}
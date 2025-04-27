import { Container, Input, Paper, Typography } from '@mui/material';
import styled from 'styled-components';

const ScanContainer = styled(Container)`
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
`;
const ScanPaper = styled(Paper)`
  border-radius: 20px;
  padding: 20px;
`;
const LoadingSpinner = styled.div`
  width: 100px;
  height: 100px;
  border: 5px solid lightgreen;
  border-bottom-color: transparent;
  border-radius: 50%;
  display: inline-block;
  box-sizing: border-box;
  animation: rotation 1s linear infinite;
`;
const WaitingContainer = styled(Container)`
  display: flex;
  align-items: center;
  flex-direction: column;
  gap: 20px;
`;
const WaitingTitle = styled(Typography)``;

const ScanFormContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 40px;
`;

const ScanFormContainerItems = styled.div`
  display: flex;
  flex-direction: column;
  gap: 40px;
  height: 250px; /* Высота блока */
  overflow-y: auto; /* Вертикальная прокрутка */
  border: 1px solid #ccc; /* Граница блока */
  padding: 10px; /* Отступы внутри блока */
`;
const ScanFormRow = styled.div`
  display: flex;
  gap: 20px;
`;

const ScanFormText = styled(Typography)`
  font-weight: 800 !important;
`;
const ScanFormInput = styled(Input)``;

const ScanButton = styled.div`
  background: #ff4e50;
  color: white;
  padding: 10px 30px;
  border-radius: 10px;
  text-align: center;
  margin-top: 50px;
`;
export {
  ScanContainer,
  LoadingSpinner,
  ScanPaper,
  WaitingContainer,
  WaitingTitle,
  ScanFormContainer,
  ScanFormText,
  ScanFormInput,
  ScanFormRow,
  ScanButton,
  ScanFormContainerItems,
};

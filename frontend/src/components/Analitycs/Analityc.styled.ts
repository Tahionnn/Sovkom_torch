import { Container, Paper, Typography } from '@mui/material';
import { PieChart } from '@mui/x-charts';
import { styled } from 'styled-components';

const AnalitycContainer = styled(Container)`
  width: 100%;
  padding: 10px 0;
  box-sizing: border-box;
`;
const AnalitycBaseChart = styled(PieChart)`
  height: 5rem;
`;
const AnalitycBaseInfo = styled(Container)`
  margin-top: 20px;
  height: 200px;
  position: relative;
`;
const AnalitycBaseCategories = styled(Container)`
  width: 100%;
  display: flex;
  flex-direction: row;
  gap: 5px;
  flex-wrap: wrap;
  margin-top: 50px;
`;
const AnalitycBaseCategory = styled.div`
  padding: 5px;
  border-radius: 10px;
  width: max-content;
  display: flex;
  gap: 5px;
`;

const AnalitycOverAll = styled(Typography)`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 60px;
  word-wrap: wrap;
  z-index: 200;
`;

const AnaliticCheckList = styled(Paper)`
  padding: 10px;
  margin-top: 50px;
  display: flex;
  flex-direction: column;
  gap: 15px;
`;
const AnalitycCheckListTitle = styled(Typography)`
  text-align: center;
  font-weight: 700 !important;
  font-size: 1.2rem !important;
`;
const AnalitycCheckListElement = styled.div`
  display: flex;
`;
const AnalitycCheckListElementTitle = styled(Typography)``;
const AnalitycCheckListElementUnderTitle = styled(Typography)`
  font-size: 12px !important;
  color: grey;
`;
const AnalitycButton = styled.div`
  background: #ff4e50;
  color: white;
  padding: 10px 30px;
  border-radius: 10px;
  text-align: center;
  margin-top: 50px;
`;
export {
  AnalitycContainer,
  AnalitycCheckListElementUnderTitle,
  AnalitycBaseCategories,
  AnalitycBaseChart,
  AnalitycButton,
  AnalitycBaseCategory,
  AnalitycBaseInfo,
  AnalitycOverAll,
  AnaliticCheckList,
  AnalitycCheckListTitle,
  AnalitycCheckListElement,
  AnalitycCheckListElementTitle,
};

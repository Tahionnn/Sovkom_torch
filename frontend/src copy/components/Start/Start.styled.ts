import { Container, Input, Typography } from '@mui/material';
import { PieChart } from '@mui/x-charts';
import styled from 'styled-components';

const StartContainer = styled(Container)`
  width: 100%;
  padding: 10px 0;
  box-sizing: border-box;
`;
const StartBaseChart = styled(PieChart)`
  height: 20rem;
`;
const StartBaseInfo = styled(Container)`
  margin-top: 20px;
  height: 200px;
`;
const StartBaseCategories = styled(Container)`
  width: 100%;
  display: flex;
  flex-direction: row;
  gap: 5px;
  flex-wrap: wrap;
  margin-top: 50px;
`;
const StartBaseCategory = styled.div`
  padding: 5px;
  border-radius: 10px;
  width: max-content;
  display: flex;
  gap: 5px;
`;

const StartBaseList = styled(Container)`
  display: flex;
  margin-top: 20px;
  flex-direction: column;
  gap: 20px;
`;

const StartBaseElement = styled.div`
  display: flex;
  align-items: center;
  justigy-content: space-between;
  gap: 10px;
  padding: 10px;
  background: rgb(15, 58, 140);
  color: white;
  border-radius: 10px;
`;
const StartBaseElementSpan = styled.span`
  height: max-content;
  color: white;
`;

const StartIconCheck = styled.img`
  width: 20px;
  height: 20px;
  color: black;
`;
const StartScanContainer = styled(Container)`
  display: flex;
  position: relative;
  margin-top: 20px;
  justify-content: center;
`;
const StartScanInput = styled(Input)`
  opacity: 0;
  position: absolute !important;
  right: 0;
  left: 0;
  top: 0;
  bottom: 0;
  z-index: 100;
  border-radius: 10px;
`;
const StartScanText = styled(Typography)`
  background: #ff4e50;
  color: white;
  padding: 5px 30px;
  border-radius: 10px;
  text-align: center;
  width: 100%;
`;
const StyledText = styled('text')(({ theme }) => ({
  textAnchor: 'middle',
  dominantBaseline: 'central',
  fontSize: 20,
}));

const StartFilterContainer = styled.div`
  margin-top: 50px;
  display: flex;
  gap: 10px;
`;
const StartFilterCategory = styled.span`
  padding: 5px 10px;
  border-radius: 20px;
  background: #ff4e50;
  color: white;
  font-size: 0.9rem;
`;
export {
  StyledText,
  StartScanText,
  StartScanInput,
  StartContainer,
  StartBaseInfo,
  StartBaseCategories,
  StartBaseCategory,
  StartBaseChart,
  StartBaseList,
  StartBaseElement,
  StartIconCheck,
  StartBaseElementSpan,
  StartScanContainer,
  StartFilterCategory,
  StartFilterContainer,
};

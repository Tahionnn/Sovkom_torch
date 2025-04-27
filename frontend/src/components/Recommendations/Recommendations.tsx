import { Container } from '@mui/material';
import { observer } from 'mobx-react-lite';
import styled from 'styled-components';
import Header from '../Header';
import DataStore from '../../stores/Data.store';
const RecsContainer = styled(Container)`
  display: flex;
  flex-direction: column;
`;
const RecsElement = styled.div``;

export const RecsButtonExport = styled.span`
  background: #ff4e50;
  position: absolute;
  right: 20px;
  top: 20px;
  color: white;
  padding: 10px 30px;
  border-radius: 10px;
  text-align: center;
`;

const Recommendations: React.FC = observer(() => {
  return (
    <>
      <RecsContainer>
        <Header></Header>
        <RecsContainer>
          {DataStore.getRecs()
            ? DataStore.getRecs().map((v: any) => {
                return RecsElement;
              })
            : ''}
        </RecsContainer>
      </RecsContainer>
    </>
  );
});

export default Recommendations;

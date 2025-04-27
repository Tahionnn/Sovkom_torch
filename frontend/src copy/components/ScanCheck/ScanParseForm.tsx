import { observer } from 'mobx-react-lite';
import DataStore from '../../stores/Data.store';
import {
  ScanButton,
  ScanFormContainer,
  ScanFormContainerItems,
  ScanFormInput,
  ScanFormRow,
  ScanFormText,
} from './ScanCheck.styled';
import { Dispatch, SetStateAction, useEffect } from 'react';
import CheckApi from '../../utils/requests/data/CheckApi';
import { routes } from '../../lists/routes';
import { useNavigate } from 'react-router';
type Props = {
  setOpen: Dispatch<SetStateAction<boolean>>;
};
const ScanParseForm: React.FC<Props> = observer(({ setOpen }) => {
  const navigate = useNavigate();
  return (
    <ScanFormContainer>
      <ScanFormText>
        Вы можете проверить правильность распознанного чека
      </ScanFormText>
      <ScanFormRow>
        <ScanFormInput
          value={DataStore.getScanParse()['shop']}
          onChange={(e) => {
            DataStore.updateParseTitle(e.target.value);
          }}
        />
      </ScanFormRow>
      <ScanFormContainerItems>
        {DataStore.getScanParse()['items'].map((v: any, i: number) => {
          return (
            <ScanFormRow>
              <ScanFormInput
                value={v['name']}
                onChange={(e) => {
                  DataStore.updateParseItems(e.target.value, 'name', i);
                }}
              />
              <ScanFormInput
                value={v['price']}
                onChange={(e) => {
                  DataStore.updateParseItems(e.target.value, 'price', i);
                }}
              />
              <ScanFormInput
                value={v['count']}
                onChange={(e) => {
                  DataStore.updateParseItems(e.target.value, 'count', i);
                }}
              />
              <ScanFormInput
                value={v['measurement']}
                onChange={(e) => {
                  DataStore.updateParseItems(e.target.value, 'measurement', i);
                }}
              />
              <ScanFormInput
                value={v['overall']}
                onChange={(e) => {
                  DataStore.updateParseItems(e.target.value, 'overall', i);
                }}
              />
            </ScanFormRow>
          );
        })}
      </ScanFormContainerItems>
      <ScanButton
        onClick={() => {
          CheckApi.sendTrueScan().then((e) => {
            setOpen(false);
            DataStore.setScan(null);
            DataStore.setScanParse(null);
            DataStore.setParseStatus(false);
            navigate(routes.start);
          });
        }}
      >
        Отправить
      </ScanButton>
    </ScanFormContainer>
  );
});

export default ScanParseForm;

import Modal from '@mui/material/Modal';
import { observer } from 'mobx-react-lite';
import { Dispatch, SetStateAction, useEffect, useState } from 'react';
import {
  LoadingSpinner,
  ScanContainer,
  ScanPaper,
  WaitingContainer,
  WaitingTitle,
} from './ScanCheck.styled';
import DataStore from '../../stores/Data.store';
import ScanParseForm from './ScanParseForm';
type Props = {
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
};
const ScanCheck: React.FC<Props> = observer(({ open, setOpen }) => {
  return (
    <Modal
      open={open}
      onClose={() => {
        setOpen(false);
      }}
    >
      <ScanContainer>
        <ScanPaper>
          <ScanParseForm setOpen={setOpen}></ScanParseForm>
        </ScanPaper>
      </ScanContainer>
    </Modal>
  );
});

export default ScanCheck;

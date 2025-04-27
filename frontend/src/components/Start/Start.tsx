import React, { useEffect, useState } from 'react';
import {
  StartBaseCategories,
  StartBaseCategory,
  StartBaseChart,
  StartBaseElement,
  StartBaseElementSpan,
  StartBaseInfo,
  StartBaseList,
  StartContainer,
  StartFilterCategory,
  StartFilterContainer,
  StartIconCheck,
  StartScanContainer,
  StartScanInput,
  StartScanText,
  StyledText,
} from './Start.styled';
import {} from '@mui/x-charts/PieChart';
import BaseApi from '../../utils/requests/data/BaseApi';
import DataStore from '../../stores/Data.store';
import { observer } from 'mobx-react-lite';
import checkIcon from '../../assets/imgs/bookmark_white.svg';
import Header from '../Header';
import { useNavigate } from 'react-router';
import { routes } from '../../lists/routes';
import ScanCheck from '../ScanCheck/ScanCheck';
import CheckApi from '../../utils/requests/data/CheckApi';
import { RecsButtonExport } from '../Recommendations/Recommendations';
type Props = {};

const Start: React.FC<Props> = observer(() => {
  const navigate = useNavigate();

  useEffect(() => {
    BaseApi.getBaseInfo();

  }, []);
  // useEffect(() => {
  //   BaseApi.getBaseList();
  // }, []);
  return (
    <>
      <StartContainer>
        <Header />
        <RecsButtonExport
          onClick={() => {
            navigate(routes.recs);
          }}
        >
          Рекоммендации
        </RecsButtonExport>
        <StartBaseInfo>
          {DataStore.getBaseData() ? (
            <StartBaseChart
              series={[
                {
                  data: DataStore.getBaseData(),
                  innerRadius: 100,
                  outerRadius: 70,
                  paddingAngle: 5,
                  cornerRadius: 5,
                  startAngle: -45,
                  cy: 95,
                },
              ]}
              hideLegend={true}
            ></StartBaseChart>
          ) : (
            ''
          )}
        </StartBaseInfo>
        <StartBaseCategories>
          {(DataStore.getBaseData() ?? []).map((v: any, i: number) => {
            return (
              <StartBaseCategory style={{ background: v.color }}>
                <span style={{ color: 'white' }}>{v.label}</span>
                <span style={{ color: 'white' }}>{v.value}</span>
              </StartBaseCategory>
            );
          })}
        </StartBaseCategories>
        <StartScanContainer>
          <StartScanInput
            type="file"
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
              CheckApi.getParseCheck();
            }}
          />
          <StartScanText>Загрузить чек</StartScanText>
        </StartScanContainer>
        <StartBaseList>
          <StartFilterContainer>
            <StartFilterCategory
              onClick={() => {
                BaseApi.getBaseListWeek();
              }}
            >
              Нед
            </StartFilterCategory>
            <StartFilterCategory
              onClick={() => {
                BaseApi.getBaseListMonth();
              }}
            >
              Месяц
            </StartFilterCategory>
            <StartFilterCategory
              onClick={() => {
                BaseApi.getBaseListYear();
              }}
            >
              Год
            </StartFilterCategory>
          </StartFilterContainer>
          {console.log(DataStore.getList())}
          {DataStore.getList()? (DataStore.getList() ?? []).map((v: any) => {
            return v['is_json_ready'] ? (
              <StartBaseElement
                onClick={() => {
                  DataStore.setCheckList(v['json_check']['gt_parse']);
                  if (
                    DataStore.getCheckList() &&
                    DataStore.getCheckList() != undefined
                  ) {
                    navigate(routes.check + '/' + v['uuid']);
                  }
                }}
              >
                <StartIconCheck src={checkIcon} />
                <StartBaseElementSpan>
                  {v['json_check']['gt_parse']['shop']}
                </StartBaseElementSpan>
                <StartBaseElementSpan>
                  {v['json_check']['gt_parse']['overall']} р.
                </StartBaseElementSpan>
              </StartBaseElement>
            ) : (
              
              <StartBaseElement
                style={{ background: 'grey', justifyContent: 'center' }}
              >
                <StartBaseElementSpan>В обработке</StartBaseElementSpan>
              </StartBaseElement>
            );
          }):""}
        </StartBaseList>
      </StartContainer>
    </>
  );
});

export default Start;

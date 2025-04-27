import { useEffect, useState } from 'react';
import Header from '../Header';
import {
  AnalitycContainer,
  AnalitycBaseChart,
  AnalitycBaseInfo,
  AnalitycBaseCategories,
  AnalitycBaseCategory,
  AnalitycOverAll,
  AnaliticCheckList,
  AnalitycCheckListTitle,
  AnalitycCheckListElement,
  AnalitycCheckListElementTitle,
  AnalitycCheckListElementUnderTitle,
  AnalitycButton,
} from './Analityc.styled';
import { observer } from 'mobx-react-lite';
import { useNavigate, useParams } from 'react-router';
import CheckApi from '../../utils/requests/data/CheckApi';
import DataStore from '../../stores/Data.store';
import { routes } from '../../lists/routes';
import ScanCheck from '../ScanCheck/ScanCheck';
type Props = {};
const Analityc: React.FC<Props> = observer(() => {
  const navigate = useNavigate();
  const params = useParams();
  useEffect(() => {
    if (params['id']) {
      CheckApi.getInfo(params['id']);
    }
  }, []);
  console.log(DataStore.getCheckList())
  const [open, setOpen] = useState<boolean>(false);

  return (
    <AnalitycContainer>
      <ScanCheck open={open} setOpen={setOpen} />
      <Header />
      {DataStore.getCheckData() && DataStore.getCheckData() != undefined ? (
        <AnalitycBaseInfo>
          <AnalitycBaseChart
            series={[
              {
                data: DataStore.getCheckData(),
                innerRadius: 100,
                outerRadius: 70,
                paddingAngle: 5,
                cornerRadius: 5,
                startAngle: -45,
                cy: 95,
              },
            ]}
            hideLegend={true}
          ></AnalitycBaseChart>
          <AnalitycOverAll>
            {DataStore.getCheckList() &&
              DataStore.getCheckList()['overall'] + 'р.'}
          </AnalitycOverAll>
        </AnalitycBaseInfo>
      ) : (
        ''
      )}
      <AnalitycBaseCategories>
        {(DataStore.getCheckData() ?? []).map((v: any, i: number) => {
          return (
            <AnalitycBaseCategory style={{ background: v.color }}>
              <span style={{ color: 'white' }}>{v.label}</span>
              <span style={{ color: 'white' }}>{v.value}</span>
            </AnalitycBaseCategory>
          );
        })}
      </AnalitycBaseCategories>
      <AnalitycButton
        onClick={() => {
          DataStore.setScanParse(DataStore.getCheckList());
          setOpen(true);
        }}
      >
        Редактировать
      </AnalitycButton>
      {DataStore.getCheckList() ? (
        <AnaliticCheckList>
          <AnalitycCheckListTitle>
            {DataStore.getCheckList()['shop']}
          </AnalitycCheckListTitle>
          {DataStore.getCheckList()['items'].map((v: any) => {
            return (
              <AnalitycCheckListElement>
                <div style={{ width: '70%', overflow: 'hidden' }}>
                  <AnalitycCheckListElementTitle>
                    {v['name']}
                  </AnalitycCheckListElementTitle>
                  <AnalitycCheckListElementUnderTitle>
                    {v['count'] + ' x ' + v['price']}
                  </AnalitycCheckListElementUnderTitle>
                </div>
                <div style={{ width: '30%' }}>
                  <AnalitycCheckListElementTitle>
                    {v['overall'] + ' р.'}
                  </AnalitycCheckListElementTitle>
                </div>
              </AnalitycCheckListElement>
            );
          })}
        </AnaliticCheckList>
      ) : (
        ''
      )}
      <AnalitycButton
        onClick={() => {
          navigate(routes.start);
        }}
      >
        На главную
      </AnalitycButton>
    </AnalitycContainer>
  );
});

export default Analityc;

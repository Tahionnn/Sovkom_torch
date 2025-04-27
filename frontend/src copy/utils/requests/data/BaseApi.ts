import { apies } from '../../../lists/apies';
import DataStore from '../../../stores/Data.store';
import { createColors, createUrl } from '../../other';
import axios from 'axios';

class BaseApi {
  async getBaseInfo() {
    try {
      let res = await axios.get(
        createUrl(apies.data.baseHost, [apies.data.baseData]),
        {
          headers: {
            Authorization: 'Bearer ' + localStorage.getItem('token'),
          },
        }
      );
      const colors = createColors();
      DataStore.setList(res.data["receipts"]);
      const data = (res.data['analytics']).map((v: any, i: number) => {
        v['color'] = colors[i];
        return v;
      });

      DataStore.setBaseData(data);
      

    } catch (error) {
      throw error;
    }
  }
  async getBaseList() {
    try {
      const res = await axios.get(
        createUrl(apies.data.baseHost, [apies.data.list]),
        {
          headers: {
            Authorization: 'Bearer ' + localStorage.getItem('token'),
          },
        }
      );
      DataStore.setList(res.data);
    } catch (error) {
      throw error;
    }
  }
  async getBaseListWeek() {
    try {
      const res = await axios.get(
        createUrl(apies.data.baseHost, [apies.data.week]),
        {
          headers: {
            Authorization: 'Bearer ' + localStorage.getItem('token'),
          },
        }
      );
      const colors = createColors();
      DataStore.setList(res.data["receipts"]);
      const data = (res.data['analytics']).map((v: any, i: number) => {
        v['color'] = colors[i];
        return v;
      });

      DataStore.setBaseData(data);
    } catch (error) {
      throw error;
    }
  }
  async getBaseListMonth() {
    try {
      const res = await axios.get(
        createUrl(apies.data.baseHost, [apies.data.month]),
        {
          headers: {
            Authorization: 'Bearer ' + localStorage.getItem('token'),
          },
        }
      );
      const colors = createColors();
      DataStore.setList(res.data["receipts"]);
      const data = (res.data['analytics']).map((v: any, i: number) => {
        v['color'] = colors[i];
        return v;
      });

      DataStore.setBaseData(data);
    } catch (error) {
      throw error;
    }
  }
  async getBaseListYear() {
    try {
      const res = await axios.get(
        createUrl(apies.data.baseHost, [apies.data.year]),
        {
          headers: {
            Authorization: 'Bearer ' + localStorage.getItem('token'),
          },
        }
      );
      const colors = createColors();
      DataStore.setList(res.data["receipts"]);
      const data = (res.data['analytics']).map((v: any, i: number) => {
        v['color'] = colors[i];
        return v;
      });

      DataStore.setBaseData(data);
    } catch (error) {
      throw error;
    }
  }

  async getRecs() {
    try {
      const res = await axios.get(
        createUrl(apies.data.baseHost, [apies.data.recs]),
        {
          headers: {
            Authorization: 'Bearer ' + localStorage.getItem('token'),
          },
        }
      );
      
    } catch (error) {
      throw error;
    }
  }
}

export default new BaseApi();

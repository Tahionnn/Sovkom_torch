import DataStore from '../../../stores/Data.store';
import { apies } from '../../../lists/apies';
import { createColors, createUrl } from '../../other';
import axios from 'axios';

class CheckApi {
  async getInfo(id: string) {
    try {
      const res = await axios.get(
        createUrl(apies.data.demoHost, [apies.data.checkInfo]),
        {
          headers: {
            Authorization: 'Bearer ' + localStorage.getItem('token'),
          },
        }
      );
      const colors = createColors();
      DataStore.setCheckList(res.data[0]['list']['gt_parse'])
      const data = res.data[0]["analytics"].map((v: any, i: number) => {
        v['color'] = colors[i];
        return v;
      });
      DataStore.setCheckData(data);
    } catch (error) {
      throw error;
    }
  }
  async getParseCheck() {
    try {
      const formdata = new FormData();
      formdata.append('file', DataStore.getScan());
      const res = await axios.get(
        createUrl(apies.data.demoHost, [apies.data.check]),
        {
          headers: {
            Authorization: 'Bearer ' + localStorage.getItem('token'),
          },
        }
      );

      if (res) {
        DataStore.setParseStatus(true);
      }
      DataStore.setScanParse(res.data[0]);
      return true;
    } catch (error) {
      throw error;
    }
  }

  async sendTrueScan() {
    try {
      //const res = await axios.post(createUrl(apies.data.demoHost,[]))
      return true;
    } catch (error) {
      throw error;
    }
  }
}

export default new CheckApi();

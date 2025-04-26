import { makeAutoObservable } from 'mobx';
type BaseData = {
  label: string;
  value: number;
};
class DataStore {
  private _baseData: any;
  private _list: any;
  private _checkData: any;
  private _checkList: any;

  private _scan: any;
  private _scanParse: any;

  private _parseStatus: boolean = false;

  private _recs: any = [];

  constructor() {
    makeAutoObservable(this);
  }

  getBaseData() {
    return this._baseData;
  }

  setBaseData(baseData: any) {
    this._baseData = baseData;
  }

  getCheckData() {
    return this._checkData;
  }

  setCheckData(checkData: any) {
    this._checkData = checkData;
  }

  getCheckList() {
    return this._checkList;
  }

  setCheckList(checkList: any) {
    this._checkList = checkList;
  }

  getList() {
    return this._list;
  }

  setList(list: any) {
    this._list = list;
  }

  setScan(scan: any) {
    this._scan = scan;
  }

  getScan() {
    return this._scan;
  }

  setScanParse(scanParse: any) {
    this._scanParse = scanParse;
  }

  getScanParse() {
    return this._scanParse;
  }

  setParseStatus(parseStatus: any) {
    this._parseStatus = parseStatus;
  }

  getParseStatus() {
    return this._parseStatus;
  }
  updateParseTitle(value: any) {
    let copy = Object.assign(this._scanParse);
    copy['shop'] = value;
    this._scanParse = copy;
  }
  updateParseItems(value: any, key: any, i: number) {
    let copy = Object.assign(this._scanParse);
    copy['items'][i][key] = value;
    this._scanParse = copy;
  }

  getRecs() {
    return this._recs;
  }

  setRecs(recs: any) {
    this._recs = recs;
  }
}

export default new DataStore();

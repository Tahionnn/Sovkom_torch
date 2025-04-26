import { makeAutoObservable} from "mobx"
class MainStore{
    constructor(){
        makeAutoObservable(this);
    }

    private _checkAuthVariant = true;

    setCheckAuthVariant(checkAuthVariant: boolean){
        this._checkAuthVariant=checkAuthVariant;
    }

    getCheckAuthVariant(){
        return this._checkAuthVariant;
    }
}

export default new MainStore();
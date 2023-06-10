import {Injectable} from "@angular/core";

@Injectable({
  providedIn: 'root'
})
export class StoreService {

  private storeMap: Map<String, any> = new Map();

  constructor() {
  }

  public getItem(key: string) {
    return this.storeMap.get(key);
  }

  public setItem(key:string ,value: any) {
    this.storeMap.set(key, value);
  }

  public removeItem(key: string) {
    this.storeMap.delete(key);
  }

}

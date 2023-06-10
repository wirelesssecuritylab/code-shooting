import {Injectable} from "@angular/core";
import {Observable} from "rxjs";
import {HttpClient, HttpHeaders} from "@angular/common/http";

@Injectable({
  providedIn: 'root'
})
export class DashboardService {

  constructor(private http: HttpClient) {
  }

  public getDashboardConfig(): Observable<any> {
    const header = new HttpHeaders().set("Content-type","application/json; charset=UTF-8");
    return this.http.get('assets/config/dashboard.json',{headers: header,responseType: 'json'});
  }
}

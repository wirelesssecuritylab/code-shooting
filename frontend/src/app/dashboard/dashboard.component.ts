import { Component, OnInit } from '@angular/core';
import {DashboardService} from "./dashboard.service";
import {DomSanitizer} from "@angular/platform-browser";
import {Router} from "@angular/router";
import {IOption} from "paletx";

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  dashboards: any[] = [];
  dashboardUrls: string[] = [];
  type: string;
  userId : string;
  department: string;
  center: string;
  institute: string;
  periodOptions:Array<IOption> = [
    {value: 'day', label: '最近一天'},
    {value: 'week', label: '最近一周'},
    {value: 'month', label: '最近一月'},
    {value: 'quarter', label: '最近一季度'}];
  periodParams = {'day': '&from=now-24h&to=now',
    'week': '&from=now-7d&to=now',
    'month': '&from=now-30d&to=now',
    'quarter': '&from=now-90d&to=now'
  };
  periodStr = {'day': '一天',
    'week': '一周',
    'month': '一月',
    'quarter': '一季度'
  };
  selectPeriod = 'month';

  constructor(private dasboardService: DashboardService,
              private sanitizer: DomSanitizer,
              private router: Router) { }

  ngOnInit(): void {
    if(this.router.routerState.snapshot.url.indexOf('manage') > -1) {
      this.type = 'manage';
      this.department = localStorage.getItem('department');
      this.institute = localStorage.getItem('institute');
    } else {
      this.type = 'user';
      this.userId = localStorage.getItem('user');
      if(!this.userId || this.userId == 'null') {
        localStorage.clear();
        this.router.navigateByUrl('/login');
      }
    }
    this.initDashboard();
  }

  initDashboard() {
    this.dasboardService.getDashboardConfig().subscribe((config) => {
      this.dashboardUrls = [...config[this.type]].map((dashboard) => dashboard.url);
      this.dashboards = [...config[this.type]].map((dashboard) => {
        return {
          name: dashboard.name,
          url: this.type == 'manage' ? this.getManagerUrl(dashboard) : this.getUserUrl(dashboard),
          widthParams: dashboard.widthParams
        };
      });
    });
  }

  getManagerUrl(dashboard: any) {
    return this.sanitizer.bypassSecurityTrustResourceUrl(`${dashboard.url}${this.periodParams[this.selectPeriod]}&var-department=${this.department}&var-institute=${this.institute}&var-cycleTxt=${this.periodStr[this.selectPeriod]}`);
  }

  getUserUrl(dashboard: any) {
    return this.sanitizer.bypassSecurityTrustResourceUrl(`${dashboard.url}${this.periodParams[this.selectPeriod]}&var-userId=${this.userId}&var-cycleTxt=${this.periodStr[this.selectPeriod]}`);
  }

  onSelectPeriod() {
    this.dashboards.forEach((dashboard,index) => {
      if(this.type == 'manage') {
        dashboard.url = this.sanitizer.bypassSecurityTrustResourceUrl(`${this.dashboardUrls[index]}${this.periodParams[this.selectPeriod]}&var-department=${this.department}&var-institute=${this.institute}&var-cycleTxt=${this.periodStr[this.selectPeriod]}`);
      }else {
        dashboard.url = this.sanitizer.bypassSecurityTrustResourceUrl(`${this.dashboardUrls[index]}${this.periodParams[this.selectPeriod]}&var-userId=${this.userId}&var-cycleTxt=${this.periodStr[this.selectPeriod]}`);
      }
    });
  }
}

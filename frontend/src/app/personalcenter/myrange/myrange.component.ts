import { Component, OnInit } from '@angular/core';
import {TypeComponent} from "../../list/type/type.component";
import {NameComponent} from "../../list/name/name.component";
import {OperationComponent} from "../../list/operation/operation.component";
import {ListService} from "../../list/list.service";
import {typeAdapter} from "../../shared/class/row";
import {Observable} from "rxjs";
import {map} from "rxjs/operators";
import {HttpClient} from "@angular/common/http";

declare let $: any;

@Component({
  selector: 'app-myrange',
  templateUrl: './myrange.component.html',
  styleUrls: ['./myrange.component.css']
})
export class MyrangeComponent implements OnInit {

  userId: string
  customBtns: any;
  columns = [];
  pageSizeSelections: number[] = [10, 30, 50];
  showLoading = false;
  data: any[] = [];
  typeAdapter = typeAdapter;
  scroll = {y: ($(document).height() - 320 + "px")}

  constructor(private listService: ListService, private http: HttpClient) { }

  ngOnInit(): void {
    this.userId = localStorage.getItem('user') ;
    this.setColumns();
    this.getShootedList();
  }

  setColumns(): void {
    this.columns = [
      {
        key: 'id',
        show: false,
      },
      {
        key: 'type',
        title: '打靶类型',
        show: true,
        width: '140px',
        contentType: 'component',
        component: TypeComponent
      },
      {
        key: 'name',
        title: '靶场',
        filter: true,
        show: true,
        width: '280px',
        contentType: 'component',
        component: NameComponent
      },
      {
        key: 'languages',
        show: false,
      },
      {
        key: 'start_at',
        show: false,
      },
      {
        key: 'end_at',
        show: false,
      },
      {
        key: 'operation',
        title: '操作',
        show: true,
        class: 'plx-table-operation',
        width: '100px',
        contentType: 'component',
        component: OperationComponent,
        inputs: {
          isMyRange: true
        }
      },
    ];

    this.customBtns = {
      iconBtns: [
        {
          tooltipInfo: '刷新',
          placement: '',
          class: 'plx-ico-refresh-16',
          callback: this.refresh.bind(this)
        },
      ]
    }
  }

  getShootedList() {
    const url = `/api/code-shooting/v1/range/shooted/user/${this.userId}/list`;
    if(!this.userId) {
      return;
    }
    this.showLoading = true;
    this.getRangeList(url).subscribe(res => {
      this.showLoading = false;
      this.data = res;
      if(this.data) {
        this.data = this.listService.typeMap(this.data, this.typeAdapter);
      }
    });
  }

  refresh(): void {
    this.getShootedList();
  }

  public getRangeList(url: string): Observable<any> {
    return this.http.get(url).pipe(
      map((res: any) => {
        if (res && res.result == 'success' && res.detail && Array.isArray(res.detail)) {
          return res.detail.map((range: any) => {
            range.targetId = range.id;
            range.start_at = new Date(range.startTime * 1000);
            range.end_at = new Date(range.endTime * 1000);
            range.languages = range.language;
            const languages = [];
            range.targets.forEach(function (value) {
              languages.push(value.language);
            });
            range.languages = [...new Set(languages)].sort();
            return range;
          });
        }
      })
    );
  }
}

import {ChangeDetectorRef, Component, Input, OnInit, ViewContainerRef, ViewChild, TemplateRef} from '@angular/core';
import {Router} from "@angular/router";
import {PlxMessage} from "paletx";
import {ListService} from "./list.service";
import {TypeComponent} from "./type/type.component";
import {NameComponent} from "./name/name.component";
import {OperationComponent} from "./operation/operation.component";
import {typeAdapter} from "../shared/class/row";

@Component({
  selector: 'app-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.css']
})
export class ListComponent implements OnInit {
  customBtns: any;
  requestBody: any;
  scroll={y: '500px'}
  public columns = [];
  public typeAdapter = typeAdapter;
  showLoading=false;
  data: any[] = [];
  cache: any[] = [];
  public pageSizeSelections: number[] = [10, 30, 50];
  public userId: string = '';
  @Input() size:any;
  constructor(
    public router: Router,
    private listService: ListService,
    private vRef: ViewContainerRef,
    private plxMessageService: PlxMessage,
    private cdr: ChangeDetectorRef
  ) { }

  private setColumns(): void {
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
        component: OperationComponent
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

  public options = [
    {
      key: 'type',
      label: '打靶类型',
      isSelected: true,
      values: [
        {key: '练习', name: '练习'},
        {key: '比赛', name: '比赛'},
      ]
    },
    {
      key: 'languages',
      label: '打靶语言',
      isSelected: true,
      values: [
      ]
    },
  ]

  ngOnInit(): void {
    this.userId = localStorage.getItem('user') ;
    this.setColumns();
    // this.data = this.listService.typeMap(this.targets, this.typeAdapter);
    this.getList();
    this.cdr.detectChanges();
  }

  getList() {
    if(!this.userId) {
      return;
    }
    this.listService.getRangeList(this.userId).subscribe(res => {
      this.data = res;
      this.initLangTagOptions(res);
      this.data = this.listService.typeMap(this.data, this.typeAdapter);
      this.cache = this.data;
      this.data = this.listService.queryLocalData(this.data, this.requestBody);
      this.cdr.detectChanges();
    });
  }

  /**
   * 直接从后端返回的靶场数据中初始化语言筛选项
   * @param data
   */
   initLangTagOptions(data) {
    let langTag: string[] = [];
    (data || []).forEach(item => {
      langTag = langTag.concat(item.languages);
    });
    let newLangTag = [...new Set(langTag)].sort();
    this.options[1].values = newLangTag.map(item => {
      return {key: item, name: item};
    });
  }

  optionChange(event: any) {
    this.data = this.cache;
    this.requestBody = this.listService.getFilterReq(event?.parent);
    console.log(this.requestBody);
    this.data = this.listService.queryLocalData(this.data, this.requestBody);
  }

  refresh(): void {
    this.showLoading = true;
    setTimeout(() => {
      this.showLoading = false;
      this.getList()
    }, 500);
  }
}

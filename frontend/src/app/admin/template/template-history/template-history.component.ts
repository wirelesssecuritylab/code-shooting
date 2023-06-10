import {Component, OnInit, TemplateRef, ViewChild} from '@angular/core';
import {PlxBreadcrumbItem, PlxMessage} from "paletx";
import {ManageService} from "../../manage.service";
import {CommonService} from "../../../shared/service/common.service";
import {errCodeAdapter} from "../../../shared/class/row";

@Component({
  selector: 'app-template-history',
  templateUrl: './template-history.component.html',
  styleUrls: ['./template-history.component.css']
})
export class TemplateHistoryComponent implements OnInit {

  @ViewChild('statusTemplate', {static: true}) public statusTemplate: TemplateRef<any>;

  breadModel: PlxBreadcrumbItem[] = [];
  templateHistoryList: any[];
  columns: any[];
  customBtns: any;
  pageSizeSelections: number[] = [10, 30, 50];
  showLoading = false;
  wsMap: Map<string, string> = new Map();

  constructor(private manageService: ManageService,
              private plxMessageService: PlxMessage,
              private commonService: CommonService,) { }

  ngOnInit(): void {
    this.initBreadModel();
    this.initColumns();
    this.initCustomBtns();
    this.loadWorkSpacesAndTemplateHis();
  }

  initBreadModel() {
    this.breadModel = [
      {label: '规范管理', routerLink: '../list', name: 'templateManage'},
      {label: '操作记录', name: 'operationRecord'},
    ];
  }

  initColumns() {
    this.columns = [
      {
        key: 'id',
        show: false,
      },
      {
        key: 'action',
        title: '操作',
        width: '100px',
        filter: true,
        format: (value: string) => {
          let actionMap = {'add': '上传', 'delete':'删除', 'enable':'启用', 'disable': '停用'};
          return actionMap[value];
        },
      },
      {
        key: 'workspace',
        title: '工作空间',
        width: '120px',
        fixed: true,
        filter: true,
        format: (value: string) => {
          return this.wsMap.has(value) ? this.wsMap.get(value) : value;
        }
      },
      {
        key: 'currentVersion',
        title: '版本',
        width: '100px',
        filter: true,
      },
      {
        key: 'changlog',
        title: '操作日志',
        width: '150px',
        filter: true,
      },
      {
        key: 'operator',
        title: '操作人',
        width: '150px',
        filter: true,
      },
      {
        key: 'opTime',
        title: '操作时间',
        width: '150px',
        filter: true,
        sort: 'desc',
        format: (value: string) => {
          return this.commonService.formatDateTime(+value);
        },
      },
      {
        key: 'opStatus',
        title: '操作结果',
        width: '100px',
        filter: true,
        contentType: 'template',
        template: this.statusTemplate,
      },
    ];
  }

  initCustomBtns() {
    this.customBtns = {
      iconBtns: [
        {
          tooltipInfo: '刷新',
          placement: '',
          class: 'plx-ico-refresh-16',
          callback: this.refresh.bind(this)
        }
      ]
    };
  }

  loadWorkSpacesAndTemplateHis() {
    this.manageService.getWorkspace().subscribe(res => {
      const workspaceInfos = [...res];
      workspaceInfos.forEach(ws => {
        this.wsMap.set(ws.id,ws.name);
      });
      this.loadTemplateHistory();
    }, err => {
      this.plxMessageService.error('获取工作空间失败！', errCodeAdapter[err.error?.errCode]);
    });
  }

  refresh() {
    this.loadTemplateHistory();
  }

  loadTemplateHistory() {
    this.showLoading = true;
    let reqBody: any = {
      name: 'queryTemplateOpHistory',
      parameters: {}
    };
    this.manageService.manageTemplateApi(reqBody).subscribe(res => {
      this.showLoading = false;
      if (res && res.result == 'success' && res.detail && Array.isArray(res.detail)) {
        this.templateHistoryList = [...res.detail];
      }
    }, err => {
      this.showLoading = false;
      this.plxMessageService.error('获取规范操作历史记录失败！', errCodeAdapter[err.error?.errCode]);
    });
  }

}

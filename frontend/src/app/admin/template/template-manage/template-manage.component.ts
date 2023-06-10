import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { ManageService } from "../../manage.service";
import { PlxMessage, PlxModal, PlxModalOptions } from "paletx";
import { CommonService } from "../../../shared/service/common.service";
import { ActivatedRoute, Router } from "@angular/router";
import { errCodeAdapter } from "../../../shared/class/row";
import { TemplateAddComponent } from "../template-add/template-add.component";


@Component({
  selector: 'app-template-manage',
  templateUrl: './template-manage.component.html',
  styleUrls: ['./template-manage.component.css']
})
export class TemplateManageComponent implements OnInit {

  @ViewChild('operTemplate', { static: true }) public operTemplate: TemplateRef<any>;
  @ViewChild('statusTemplate', { static: true }) public statusTemplate: TemplateRef<any>;
  @ViewChild('deleteConfirm', { static: true }) public deleteConfirm: TemplateRef<any>;
  @ViewChild('operConfirm', { static: true }) public operConfirm: TemplateRef<any>;
  user: string;
  templateList: any[];
  columns: any[];
  customBtns: any;
  pageSizeSelections: number[] = [10, 30, 50];
  showLoading = false;
  deleteModal: any;
  curRowData: any;
  operModal: any;
  operation: string;
  workspaceInfos: any[];
  wsMap: Map<string, string> = new Map();
  userRole: string;
  ADMINROLE: string = 'admin';

  constructor(private manageService: ManageService,
    private plxMessageService: PlxMessage,
    private commonService: CommonService,
    private router: Router,
    private activatedRoute: ActivatedRoute,
    private modalService: PlxModal, ) {
    this.userRole = localStorage.getItem('role')
  }

  ngOnInit(): void {
    this.user = localStorage.getItem('name') + localStorage.getItem('user');
    this.initColumns();
    this.initCustomBtns();
    this.loadWorkSpaceAndTemplates();
  }

  initColumns() {
    this.columns = [
      {
        key: 'id',
        show: false,
      },
      {
        key: 'workspace',
        title: '工作空间',
        width: '120px',
        fixed: true,
        filter: true,
        sort: 'desc',
        sortFunction: (direction: number, a: any, b: any) => {
          if (direction == 1) {
            return this.compareWs(a, b);
          }
          return this.compareWs(b, a);
        },
        contentType: 'html',
        format: (value: string) => {
          return this.findWsName(value);
        }
      },
      {
        key: 'version',
        title: '版本号',
        width: '120px',
        fixed: true,
        filter: true,
        sort: 'desc',
        sortFunction: (direction: number, a: any, b: any) => {
          if (direction == 1) {
            return this.compareVer(a, b);
          }
          return this.compareVer(b, a);
        },
      },
      {
        key: 'active',
        title: '使用状态',
        width: '80px',
        filter: true,
        sort: 'none',
        contentType: 'template',
        template: this.statusTemplate,
      },
      {
        key: 'uploadBy',
        title: '上传者',
        width: '150px',
        filter: true,
        sort: 'none',
      },
      {
        key: 'uploadAt',
        title: '上传时间',
        width: '150px',
        filter: true,
        sort: 'none',
        format: (value: string) => {
          return this.commonService.formatDateTime(+value);
        },
      },
      {
        key: 'operation',
        title: '操作',
        fixed: true,
        class: 'plx-table-operation',
        width: '150px',
        contentType: 'template',
        template: this.operTemplate,
      },
    ];
  }

  initCustomBtns() {
    let btns = [{
      tooltipInfo: '刷新',
      placement: '',
      class: 'plx-ico-refresh-16',
      callback: this.refresh.bind(this)
    }]
    if (this.userRole == this.ADMINROLE) {
      btns.push({
        tooltipInfo: '查看操作记录',
        placement: '',
        class: 'plx-ico-history-record-16',
        callback: this.showOperationRecord.bind(this)
      })
    }
    this.customBtns = {
      iconBtns: btns
    };
  }

  refresh() {
    this.loadTemplates();
  }

  showOperationRecord() {
    let routeStr: string = '../history';
    this.router.navigate([routeStr], {
      relativeTo: this.activatedRoute
    });
  }

  loadWorkSpaceAndTemplates() {
    this.manageService.getWorkspace().subscribe(res => {
      this.workspaceInfos = [...res];
      this.workspaceInfos.forEach(ws => {
        this.wsMap.set(ws.id, ws.name);
      });
      this.loadTemplates();
    }, err => {
      this.plxMessageService.error('获取工作空间失败！', errCodeAdapter[err.error?.errCode]);
    });
  }

  loadTemplates() {
    this.showLoading = true;
    let reqBody: any = {
      name: 'queryTemplate',
      parameters: {}
    };
    this.manageService.manageTemplateApi(reqBody).subscribe(res => {
      this.showLoading = false;
      if (res && res.result == 'success' && res.detail && Array.isArray(res.detail)) {
        const tmpList = [];
        res.detail.forEach(tmp => {
          if (tmp.workspace == '') {
            tmp.workspace = 'public';
          }
          tmpList.push(tmp);
        });
        this.templateList = [...tmpList];
      }
    }, err => {
      this.showLoading = false;
      this.plxMessageService.error('获取规范列表失败！', errCodeAdapter[err.error?.errCode]);
    });
  }

  uploadTemplate() {
    /*let routeStr: string = '../add';
    this.router.navigate([routeStr], {
      relativeTo: this.activatedRoute
    });*/
    const plxModalOptions: PlxModalOptions = { backdrop: true, keyboard: true, size: 'sm' };
    let wsAddModal = this.modalService.open(TemplateAddComponent, plxModalOptions);
    wsAddModal.componentInstance.workspaces = this.workspaceInfos;
    wsAddModal.componentInstance.callbackFunc = () => {
      this.loadTemplates();
    };
  }

  enableTemplate(rowData) {
    this.operation = 'enable';
    this.curRowData = rowData;
    this.operModal = this.modalService.open(this.operConfirm, { size: 'xs', enterEventFunc: this.operTemplateInfo.bind(this) });
  }

  disableTemplate(rowData) {
    this.operation = 'disable';
    this.curRowData = rowData;
    this.operModal = this.modalService.open(this.operConfirm, { size: 'xs', enterEventFunc: this.operTemplateInfo.bind(this) });
  }

  operTemplateInfo() {
    let reqBody: any = {
      name: this.operation,
      parameters: this.buildOper(this.operation, this.curRowData)
    };
    this.manageService.manageTemplateApi(reqBody).subscribe(() => {
      this.operModal.close();
      this.plxMessageService.success(this.operation == 'enable' ? '启用规范成功！' : '停用规范成功！', '');
      this.loadTemplates();
    }, err => {
      this.plxMessageService.error(this.operation == 'enable' ? '启用规范失败！' : '停用规范失败！', errCodeAdapter[err.error?.errCode]);
    });
  }

  buildOper(action, rowData) {
    const changeLog = {
      delete: '删除规范',
      enable: '启用规范',
      disable: '停用规范',
      download: '下载规范',
    };
    return {
      templateId: rowData.id,
      action: action,
      workspace: rowData.workspace,
      currentVersion: rowData.version,
      nextVersion: '',
      changlog: changeLog[action],
      operator: this.user,
    }
  }

  deleteTemplate(rowData) {
    this.curRowData = rowData;
    this.deleteModal = this.modalService.open(this.deleteConfirm, { size: 'xs', enterEventFunc: this.deleteTemplateInfo.bind(this) });
  }
  //规范下载
  downloadTemplate(rowData) {
    let reqBody: any = {
      name: 'download',
      parameters: this.buildOper('download', rowData)
    };
    this.manageService.downloadTemplateApi(reqBody).subscribe(
      res => {
        let fileName: string = "代码打靶落地模板 -" + rowData.version;
        const link = document.createElement('a');
        const blob = new Blob([res], { type: 'application/vnd.ms-excel' });
        link.setAttribute('href', window.URL.createObjectURL(blob));
        link.setAttribute('download', fileName + '.xlsm');
        link.style.visibility = 'hidden';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      }, err => {
        this.plxMessageService.error('导出模板文件失败！', err.cause);
      }
    );
  }

  deleteTemplateInfo() {
    let reqBody: any = {
      name: 'delete',
      parameters: this.buildOper('delete', this.curRowData)
    };
    this.manageService.manageTemplateApi(reqBody).subscribe(() => {
      this.deleteModal.close();
      this.plxMessageService.success('删除规范成功！', '');
      this.loadTemplates();
    }, err => {
      this.plxMessageService.error('删除规范失败！', errCodeAdapter[err.error?.errCode]);
    });
  }

  convert2Link(templateUrl, value) {
    if (templateUrl == "") {
      return value;
    }
    return "<a class='plx-link' target='_blank' style='cursor:pointer' href='" + templateUrl + "'>" + value + "</a>";
  }

  convert2Value(link) {
    const matches = link.match(/<a\s+[^>]*>([^<\/]+)<\/a>/);
    if (matches) {
      return matches[1];
    }
    return link;
  }

  compareVer(ver1, ver2) {
    const vers1 = this.splitVersion(ver1);
    const vers2 = this.splitVersion(ver2);
    const largeVerCompare = vers1[0] - vers2[0];
    if (largeVerCompare == 0) {
      return vers1[1] - vers2[1];
    }
    return largeVerCompare;
  }

  compareWs(ws1, ws2) {
    const wsStr1 = this.convert2Value(ws1);
    const wsStr2 = this.convert2Value(ws2);
    return wsStr1.localeCompare(wsStr2);
  }

  private splitVersion(version) {
    const verStr = version.substring(1);
    const vers = verStr.split('.');
    return [+vers[0], +vers[1]];
  }

  private findWsName(workspaceId) {
    const ws = this.workspaceInfos.filter(ws => ws.id == workspaceId);
    if (ws.length > 0) {
      return this.convert2Link(ws[0].url, ws[0].name);
    }
    return this.convert2Link("", workspaceId);
  }

}

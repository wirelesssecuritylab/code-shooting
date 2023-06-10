import { Component, OnInit, ViewChild, TemplateRef } from '@angular/core';
import { Router, ActivatedRoute } from "@angular/router";
import { TargetOperationComponent } from '../target-operation/target-operation.component';
import { ManageService } from '../../manage.service';
import { languageAdapter, langReverseAdapter, errCodeAdapter } from "../../../shared/class/row";
import { IOption, PlxMessage, PlxModal } from 'paletx';

const targetType = { true: '是', false: '否' };
const _DELETE = 'remove';

@Component({
  selector: 'app-target-list',
  templateUrl: './target-list.component.html',
  styleUrls: ['./target-list.component.css']
})
export class TargetListComponent implements OnInit {
  @ViewChild('tableTagColumTemplate', { static: true }) public tableTagColumTemplate: TemplateRef<any>;
  @ViewChild('targetIdTemplate', { static: true }) public targetIdTemplate: TemplateRef<any>;
  @ViewChild('operTemplate', { static: true }) public operTemplate: TemplateRef<any>;
  public userId: string;
  public customBtns: any;
  public columns = [];
  public sonColumns = [];
  public data: any[] = [];
  public tagDetailData: any[] = [];
  public showLoading: boolean = false;
  public pageSizeSelections: number[] = [10, 30, 50];
  public showAdvanceQuery: boolean = false;
  public languagesOption: Array<IOption> = [];
  public selectedLanguage: string;
  public modal: any;
  public isOpen = true;
  public curRowData: any;
  public wsMap: Map<string, string> = new Map();
  public role: string;
  public privilege: string;

  constructor(private router: Router, private manageService: ManageService,
    private plxMessageService: PlxMessage,
    private activatedRoute: ActivatedRoute,
    private modalService: PlxModal) { }

  ngOnInit(): void {
    this.userId = localStorage.getItem('user');
    this.role = localStorage.getItem('role');
    this.privilege = localStorage.getItem('expertPrivilege');
    this.initFilterOptions();
    this.setColumns();
    this.setSonColumns();
    this.loadWorkSpacesAndTargets();
  }

  checkedArray: Array<any> = [];
  checkboxInfoChange(item) {
    if (item.action === 'check') {
      item.actionData.forEach((value) => {
        this.checkedArray.push(value);
      });
    } else {
      item.actionData.forEach((uncheckValue) => {
        this.checkedArray = this.checkedArray.filter((item) => {
          // 此处用数据中的id属性来做过滤
          return uncheckValue.id !== item.id;
        });
      });
    }
  }
  exportBatchTarget() {
    this.manageService.exportBatchTargetApi(this.checkedArray, this.role, this.userId).subscribe(
      res => {
        if (res.size == 67) {
          this.plxMessageService.error('只有管理员才可以导出全部靶子,普通用户只允许导出自己创建的靶子!', "");
        } else {
          const link = document.createElement('a');
          const blob = new Blob([res], { type: 'application/zip' });
          link.setAttribute('href', window.URL.createObjectURL(blob));
          let fileName: string = "niuzhi";
          link.setAttribute('download', 'targets.zip');
          link.style.visibility = 'hidden';
          document.body.appendChild(link);
          link.click();
          document.body.removeChild(link);
        }
      }, err => {
        this.plxMessageService.error('批量导出失败！', err.cause);
      }
    );

  }
  /**
     * 初始化高级查询中打靶语言筛选项
     */
  initFilterOptions() {
    this.languagesOption = [{
      label: 'Go',
      value: languageAdapter['Go']
    }, {
      label: 'C/C++',
      value: languageAdapter['C/C++']
    }, {
      label: 'C',
      value: languageAdapter['C']
    }, {
      label: 'Python',
      value: languageAdapter['Python']
    }, {
      label: 'Java',
      value: languageAdapter['Java']
    }];
  }
  /**
     * 设置靶子列表表格项
     */
  private setColumns(): void {
    this.columns = [
      {
        key: 'targetId',
        title: 'ID',
        width: '50px',
        show: false,
        contentType: 'template',
        template: this.targetIdTemplate,
      },
      {
        key: 'name',
        title: '靶子名称',
        filter: true,
        show: true,
        width: '180px',
        fixed: true,
        sort: 'none'
      },

      {
        key: 'workspace',
        title: '工作空间',
        filter: true,
        show: true,
        width: '100px',
        fixed: true,
        sort: 'none',
        format: (value: string) => {
          return this.wsMap.has(value) ? this.wsMap.get(value) : value;
        }
      },
      {
        key: 'template',
        title: '规范版本',
        filter: true,
        show: true,
        width: '100px',
        fixed: true,
        sort: 'none',
        format: (value: string) => {
          return this.wsMap.has(value) ? this.wsMap.get(value) : value;
        }
      },
      {
        key: 'language',
        title: '所属语言',
        filter: true,
        show: true,
        width: '100px',
        sort: 'none'
      },
      {
        key: 'relatedRangesShow',
        title: '关联靶场',
        show: true,
        width: '150px',
        sort: 'none'
      },
      {
        key: 'targetLabel',
        title: '靶子标签',
        show: true,
        width: '100px',
        sort: 'none',
        showInDropdown: true,
        class: 'plx-table-operation',
        contentType: 'template',
        template: this.tableTagColumTemplate,
        showInDetail: false,
      },
      {
        key: 'extendedLabel',
        title: '扩展标签',
        show: true,
        filter: true,
        width: '100px',
        sort: 'none',
      },
      {
        key: 'customLable',
        title: '自定义标签',
        show: true,
        filter: true,
        width: '100px',
        sort: 'none',
      },
      {
        key: 'instituteLabel',
        title: '院级标签',
        show: true,
        filter: true,
        width: '100px',
        sort: 'none',
      },
      {
        key: 'ownerInfo',
        title: '上传者',
        filter: true,
        show: true,
        width: '100px',
        sort: 'none'
      },
      {
        key: 'isSharedShow',
        title: '是否共享',
        filter: true,
        show: true,
        width: '80px',
        sort: 'none'
      },
      {
        key: 'operation',
        title: '操作',
        show: true,
        fixed: true,
        class: 'plx-table-operation',
        width: '120px',
        // contentType: 'component',
        // component: TargetOperationComponent,
        contentType: 'template',
        template: this.operTemplate,
      }
    ];

    this.customBtns = {
      iconBtns: [
        {
          tooltipInfo: '批量导出',
          placement: '',
          class: 'plx-ico-export-16',
          callback: this.exportBatchTarget.bind(this)
        },
        {
          tooltipInfo: '刷新',
          placement: '',
          class: 'plx-ico-refresh-16',
          callback: this.refresh.bind(this)
        }
      ],
      commonBtns: [
        // {
        //   name: '高级查询',
        //   isShowExtend: true,
        //   callback: this.advanceQuery.bind(this)
        // }
      ]
    }
  }

  /**
     * 设置子表格列表头
     */
  private setSonColumns(): void {
    this.sonColumns = [
      {
        key: 'bugType',
        title: '缺陷大类',
        sort: 'none',
      },
      {
        key: 'subBugType',
        title: '缺陷小类',
        sort: 'none',
      },
      {
        key: 'bugDetail',
        title: '缺陷细项',
        sort: 'none',
      }
    ];
  }

  /**
     * 点击刷新按钮
     */
  refresh(): void {
    this.showLoading = true;
    this.data = [];
    this.getTargetList();
  }

  /**
     * 点击高级查询按钮
     */
  advanceQuery() {
    this.showAdvanceQuery = !this.showAdvanceQuery;
    if (!this.showAdvanceQuery) {
      this.selectedLanguage = '';
    }
  }

  /**
   * 获取靶子列表
   */
  getTargetList(): void {
    let reqBody: any = {
      name: 'query',
      parameters: {
        owner: this.userId
      }
    };
    if (this.selectedLanguage) {
      reqBody.parameters['language'] = this.selectedLanguage;
    }
    this.manageService.manageTargetApi(reqBody).subscribe(res => {
      this.showLoading = false;
      if (res && res.detail && Array.isArray(res.detail)) {
        this.data = res.detail.map(target => {
          target['isSharedShow'] = target.isShared ? '是' : '否';
          target['relatedRangesShow'] = target['relatedRanges'] || '--';
          target['ownerInfo'] = target['ownerName'] + target['owner'];
          target['tagDetail'] = [{
            bugType: target['tagName']['mainCategory'] || '',
            subBugType: target['tagName']['subCategory'] || '',
            bugDetail: target['tagName']['defectDetail'] || '',

          }];
          target['extendedLabel'] = target['extendedLabel'] ? target['extendedLabel'].join(',') : '';
          return target;
        });
      }
    }, error => {
      this.showLoading = false;
      this.plxMessageService.error('获取靶子列表失败！', error.cause);
    });
  }

  /**
     * 点击新建靶场
     */
  addTarget(): void {
    this.router.navigateByUrl('/main/manage/target/add');
  }

  /**
     * 批量删除靶子(暂不实现)
     */
  patchDeleteTarget() {

  }

  /**
     * 点击展开子表格显示标签详情
     * @param data 父表格中某一行的数据
     */
  public showTagDetail(data: any) {
    data.isDropdownOpen = !data.isDropdownOpen;
  }

  /**
     * 编辑靶子操作
     */
  modifyTarget(rowData) {
    let routeStr: string = '../edit';
    this.router.navigate([routeStr], {
      queryParams: {
        'id': rowData.id,
        'language': rowData.language,
        'owner': rowData.owner,
        'ownerName': rowData.ownerName
      },
      relativeTo: this.activatedRoute
    });
  }

  /**
     * 删除靶子时打开确认删除模态框
     * @param content
     */
  open(content, data) {
    this.curRowData = data;
    this.modal = this.modalService.open(content, { size: 'xs', enterEventFunc: this.func.bind(this) });
    this.isOpen = true;
  }

  public func(): void {
    if (this.isOpen) {
      this.modal.close();
    }
    this.isOpen = !this.isOpen;
  }

  /**
     * 确认删除靶子
     * @returns
     */
  deleteTarget() {
    if (!this.curRowData.id) {
      return;
    }
    let deleteBody: any = {
      name: _DELETE,
      parameters: {
        id: this.curRowData.id,
        owner: this.curRowData.owner
      }
    };
    this.manageService.manageTargetApi(deleteBody).subscribe(res => {
      this.plxMessageService.success("删除成功！", '');
      this.getTargetList();
      if (this.checkedArray.length > 0) {
        var newcheckedArray: Array<any> = [];
        for (var i = 0; i < this.checkedArray.length; i++) {
          if (this.checkedArray[i].id != this.curRowData.id) {
            newcheckedArray.push(this.checkedArray[i]);
          }
        }
        this.checkedArray = newcheckedArray;
      }

    }, err => {
      this.plxMessageService.error('删除失败！', err.cause);
    });
  }

  loadWorkSpacesAndTargets() {
    this.manageService.getWorkspace().subscribe(res => {
      const workspaceInfos = [...res];
      workspaceInfos.forEach(ws => {
        this.wsMap.set(ws.id, ws.name);
      });
      this.getTargetList();
    }, err => {
      this.plxMessageService.error('获取工作空间失败！', errCodeAdapter[err.error?.errCode]);
    });
  }
}


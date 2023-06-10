import { Component, OnInit, ViewChild, TemplateRef } from '@angular/core';
import { Router, ActivatedRoute } from "@angular/router";
import { IOption, PlxMessage, PlxModal } from 'paletx';
import { ManageService } from 'src/app/admin/manage.service';

const targetType = { true: '是', false: '否' };
const _DELETE = 'remove';

@Component({
  selector: 'app-mytarget',
  templateUrl: './mytarget.component.html',
  styleUrls: ['./mytarget.component.css']
})
export class MyTargetComponent implements OnInit {
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

  constructor(private router: Router, private manageService: ManageService,
    private plxMessageService: PlxMessage,
    ) { }

  ngOnInit(): void {
    this.userId = localStorage.getItem('user');
    this.role = localStorage.getItem('role');
    this.setColumns();
    this.setSonColumns();
    this.getTargetList();
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
    if (this.data.length == 0) {
      this.plxMessageService.error('导出失败', '没有可供导出的数据！');
      return
    }
    if (this.checkedArray.length == 0) {
      // 没有选中时，将自己所有的靶子全部导出
      this.checkedArray = this.data
    }
    this.manageService.exportBatchTargetApi(this.checkedArray, this.role, this.userId).subscribe(
      res => {
        if (res.size == 67) {
          this.plxMessageService.error('只有管理员才可以导出全部靶子,普通用户只允许导出自己创建的靶子!', "");
        } else {
          const link = document.createElement('a');
          const blob = new Blob([res], { type: 'application/zip' });
          link.setAttribute('href', window.URL.createObjectURL(blob));
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
        key: 'template',
        title: '版本号',
        filter: true,
        show: true,
        width: '100px',
        sort: 'none'
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
        key: 'relatedRangesShow',
        title: '关联靶场',
        show: true,
        width: '150px',
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
        key: 'fileName',
        title: '文件名',
        sort: 'none',
      },
      {
        key: 'class',
        title: '缺陷大类',
        sort: 'none',
      },
      {
        key: 'subClass',
        title: '缺陷小类',
        sort: 'none',
      },
      {
        key: 'describe',
        title: '缺陷细项',
        sort: 'none',
      },
      {
        key: 'startLineNum',
        title: '起始行号',
        sort: 'none',
      },
      {
        key: 'endLineNum',
        title: '结束行号',
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
   * 获取靶子列表
   */
  getTargetList(): void {
    let reqBody: any = {
      name: 'query',
      parameters: {
        owner: this.userId
      }
    };
    this.manageService.manageTargetApi(reqBody).subscribe(res => {
      this.showLoading = false;
      if (res && res.detail && Array.isArray(res.detail)) {
        this.data = []
        res.detail.forEach(target => {
          target['isSharedShow'] = target.isShared ? '是' : '否';
          target['relatedRangesShow'] = target['relatedRanges'] || '--';
          target['ownerInfo'] = target['ownerName'] + target['owner'];
          if (this.userId == target.owner) {
            // 只展示自己创建的靶子
            this.data.push(target)
          }
        });
      }
    }, error => {
      this.showLoading = false;
      this.plxMessageService.error('获取靶子列表失败！', error.cause);
    });
  }

  /**
     * 检查靶子
     */
  checkTarget(data: any) {
    if (data.isDropdownOpen) {
      // 打开情况下直接收起
      data.isDropdownOpen = false;
      return
    }

    this.showLoading = true;
    let reqBody: any = {
      name: 'check',
      parameters: {
        id: data.id
      }
    };
    this.manageService.manageTargetApi(reqBody).subscribe(res => {
      this.showLoading = false;
      if (res && res.detail && Array.isArray(res.detail)) {
        data.tagDetail = []
        res.detail.forEach(target => {
          console.log(target)
          data.tagDetail.push(target)
        });
        this.dropdownClose()
        data.isDropdownOpen = true;
      } else {
        this.plxMessageService.info('提示', '靶子规范正确无需更新~');
      }
    }, error => {
      this.showLoading = false;
      this.plxMessageService.error('检查靶子失败！', error.cause);
    });
  }

  dropdownClose() {
    this.data.forEach(data => data.isDropdownOpen = false)
  }
}


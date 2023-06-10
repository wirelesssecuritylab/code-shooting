import { Component, OnInit, ViewChild, TemplateRef } from '@angular/core';
import {Router, ActivatedRoute} from "@angular/router";
import { RangeOperationComponent} from '../range-operation/range-operation.component';
import { ManageService } from '../../manage.service';
import { PlxMessage, PlxModal } from 'paletx';
import { rangeType } from '../../../shared/class/row';
import {CommonService} from "../../../shared/service/common.service";

const _DELETE = 'remove';

@Component({
  selector: 'app-range-list',
  templateUrl: './range-list.component.html',
  styleUrls: ['./range-list.component.css']
})
export class RangeListComponent implements OnInit {
  @ViewChild('operTemplate', {static: true}) public operTemplate: TemplateRef<any>;
  public modal:any;
  public isOpen = true;
  public curRowData: any;
  public userId: string;
  customBtns: any;
  public columns = [];
  rangeList: any[] = [];
  showLoading: boolean = false;
  public pageSizeSelections: number[] = [10, 30, 50];
  constructor(private router: Router, private manageService: ManageService,
              private plxMessageService: PlxMessage,
              private activatedRoute: ActivatedRoute,
              private modalService: PlxModal,
              private commonService: CommonService,) { }

  ngOnInit(): void {
    this.userId = localStorage.getItem('user');
    this.setColumns();
    this.getRangeList();
  }

  /**
   * 设置表格列
   */
  private setColumns(): void {
    this.columns = [
      {
        key: 'id',
        show: false,
      },
      {
        key: 'name',
        title: '靶场名称',
        show: true,
        width: '160px',
        fixed: true,
        sort: 'none',
        filter: true,
      },
      {
        key: 'typeShow',
        title: '打靶类型',
        show: true,
        width: '80px',
        sort: 'none',
        filter: true,

      },
      {
        key: 'projectName',
        title: '所属组织',
        filter: true,
        show: true,
        width: '100px',
        sort: 'none',
      },
      {
        key: 'languages',
        title: '支持语言',
        show: true,
        width: '100px',
        sort: 'none',
        filter: true,
      },
      {
        key: 'ownerInfo',
        title: '创建者',
        show: true,
        width: '100px',
        sort: 'none',
        filter: true,
      },
      {
        key: 'start_at',
        title: '开始时间',
        show: true,
        width: '150px',
        sort: 'none',
        filter: true,
      },
      {
        key: 'end_at',
        title: '结束时间',
        show: true,
        width: '150px',
        sort: 'none',
        filter: true,
      },
      {
        key: 'operation',
        title: '操作',
        show: true,
        fixed: true,
        class: 'plx-table-operation',
        width: '140px',
        // contentType: 'component',
        // component: RangeOperationComponent,
        contentType: 'template',
        template: this.operTemplate,
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

  /**
   * 点击刷新按钮
   */
  refresh(): void {
    this.showLoading = true;
    this.rangeList= [];
    this.getRangeList();
  }

  /**
   * 获取靶场列表
   */
  getRangeList(): void {
    let reqBody: any = {
      name: 'query',
      parameters: {
      user: this.userId
      }
    };
    this.manageService.manageRangeApi(reqBody).subscribe(res => {
      this.showLoading = false;
      if (res && res.result == 'success' && res.detail && Array.isArray(res.detail)) {
        this.rangeList = res.detail.map(range => {
          range['typeShow'] = rangeType[range.type];
          range['start_at'] = (range.type == 'compete' && range.startTime) ? this.commonService.formatDateTime(range.startTime) : '--';
          range['end_at'] = (range.type == 'compete' && range.endTime) ? this.commonService.formatDateTime(range.endTime) : '--';
          range['ownerInfo'] = range['ownerName'] + range['owner'];
          let languages: string[] = [];
          (range.targets || []).forEach(item => {
            if (languages.indexOf(item.language) < 0) {
              languages.push(item.language);
            }
          });
          range['languages'] = languages.join(',');
          return range;
        });
      }
    }, error => {
      this.showLoading = false;
      this.plxMessageService.error('获取靶场列表失败！', error.cause);
    });
  }

  /**
   * 点击新建靶场
   */
  addRange(): void {
    this.router.navigateByUrl('/main/manage/range/add');
  }

  /**
   * 管理员进入靶场查询或导出成绩
   */
  goToRange(rowData) {
    let routeStr: string = '../query';
    this.router.navigate([routeStr], {
      queryParams: {
        'id': rowData.id,
        'name': rowData.name,
        'languages': rowData.languages,
        'projectId': rowData.project
      },
      relativeTo: this.activatedRoute
    });
  }

  /**
   * 编辑靶场
   */
  modifyRange(rowData) {
    let routeStr: string = '../edit';
    this.router.navigate([routeStr], {
      queryParams: {
        'id': rowData.id
      },
      relativeTo: this.activatedRoute
    });
  }

  /**
   * 删除靶子时打开确认删除模态框
   * @param content
   */
  openConfirmDlg(content, rowData) {
    this.curRowData = rowData;
    this.modal = this.modalService.open(content, {size: 'xs', enterEventFunc: this.func.bind(this)});
    this.isOpen = true;
  }

  public func(): void {
    if(this.isOpen) {
      this.modal.close();
    }
    this.isOpen = !this.isOpen;
  }

  /**
   * 删除靶场
   * @returns
   */
  deleteRange() {
    if (!this.curRowData.id) {
      return;
    }
    let deleteBody: any = {
      name: _DELETE,
      parameters: {
        id: this.curRowData.id,
      }
    };
    this.manageService.manageRangeApi(deleteBody).subscribe(res => {
      this.plxMessageService.success("删除成功！", '');
      this.getRangeList();
    }, err => {
      this.plxMessageService.error('删除失败！', err.cause);
    });
  }
}

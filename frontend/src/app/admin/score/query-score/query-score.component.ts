import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { PlxBreadcrumbItem, IOption, PlxMessage, PlxModal, FieldType } from 'paletx';
import { ManageService } from '../../manage.service';

@Component({
  selector: 'app-query-score',
  templateUrl: './query-score.component.html',
  styleUrls: ['./query-score.component.css']
})
export class QueryScoreComponent implements OnInit, OnDestroy {

  public rangeId: string;
  public rangeName: string;
  public languages: string;
  public projectId: string;
  public verbose: boolean;
  public customBtns: any;
  public columns = [];
  public scoreData: any[] = [];
  public showLoading: boolean = false;
  public pageSizeSelections: number[] = [10, 30, 50];
  public breadModel: PlxBreadcrumbItem[] = [];
  public languagesOption: Array<IOption> = [];
  public departmentOption: Array<IOption> = [];
  public selectedDepartment: string[] = [];
  public selectedLanguage: string;
  public totalDptNum: number = 0;
  public isLoading: boolean;

  public modal: any;
  public formSetting: any;

  constructor(private activatedRoute: ActivatedRoute,
    private manageService: ManageService,
    private plxMessage: PlxMessage,
    private modalService: PlxModal) {
  }

  ngOnInit(): void {
    this.activatedRoute.queryParamMap.subscribe((paramMap) => {
      this.rangeId = paramMap.get('id');
      this.rangeName = paramMap.get('name');
      this.languages = paramMap.get('languages');
      this.projectId = paramMap.get('projectId');
      this.getRangeInfoAndSetBtns(this.rangeId);
      if (this.projectId != 'public') {
        this.initDepartmentOptions();
      }
    });

    this.breadModel = [
      { label: '靶场管理', routerLink: '../list', name: 'rangeManage' },
      { label: this.rangeName, name: 'link2' }
    ];

    this.setColumns();
    this.initFilterOptions();
    this.queryRangeScore();
    this.formSetting = this.initFormSetting();
  }

  /**
   * 初始打靶语言过滤项
   */
  initFilterOptions(): void {
    let lang: string[] = this.languages.split(',');
    lang.forEach((item: any) => {
      this.languagesOption.push({
        label: item,
        value: item
      });
    });
    this.selectedLanguage = this.languagesOption.length > 0 ? this.languagesOption[0].value : '';
  }

  /**
   * 初始化部门过滤项
   */
  initDepartmentOptions() {
    this.manageService.getDepartListApi(this.projectId).subscribe(res => {
      let departList: Array<IOption> = [];
      (res || []).forEach((item: any) => {
        departList.push({
          label: item,
          value: item
        });
      });
      this.departmentOption = departList;
      this.totalDptNum = this.departmentOption.length;
    }, err => {
      this.plxMessage.error('获取部门列表失败！', err.msg);
    });
  }

  /**
   * 设置表格列
   */
  setColumns(): void {
    this.columns = [
      {
        key: 'id',
        show: false,
      },
      {
        key: 'project',
        title: '项目名称',
        show: false,
        width: '100px',
      },
      {
        key: 'department',
        title: '部门名称',
        show: true,
        width: '100px',
        filter: false,
      },
      {
        key: 'teamName',
        title: '团队名称',
        show: true,
        width: '100px',
        filter: true,
      },
      {
        key: 'userName',
        title: '姓名',
        filter: true,
        show: true,
        width: '80px',
      },
      {
        key: 'userId',
        title: '工号',
        show: true,
        width: '80px',
        filter: true,
      },
      {
        key: 'hitNum',
        title: '命中靶环数',
        show: true,
        width: '80px',
        filter: false
      },
      {
        key: 'hitScore',
        title: '命中靶环总分',
        show: true,
        width: '80px',
        filter: false
      },
      {
        key: 'hitRate',
        title: '命中率',
        show: true,
        width: '80px',
        filter: false,
        sort: 'desc',
        sortIndex: 2,
        sortFunction: (direction: number, a: any, b: any) => {
          if (direction == 1) {
            return this.compareRate(a, b);
          }
          return this.compareRate(b, a);
        }
      },
      {
        key: 'hitScoreHundredth',
        title: '命中靶环总分（百分制）',
        show: true,
        width: '80px',
        filter: false,
        sort: 'desc',
        sortIndex: 1
      }
    ];
  }

  compareRate(rate1, rate2) { // 比较时，去除百分号
    const rateNum1 = +rate1.substring(0,rate1.length -1);
    const rateNum2 = +rate2.substring(0,rate2.length -1);
    return rateNum1 - rateNum2;
  }

  getRangeInfoAndSetBtns(rangeId: string) {
    let reqBody: any = {
      name: 'query',
      parameters: {
        user: localStorage.getItem('user')
      }
    };
    this.manageService.manageRangeApi(reqBody).subscribe(res => {
      if (res && res.result == 'success' && res.detail && Array.isArray(res.detail)) {
        const ranges = res.detail.map(range => {
          return {
            id: range.id,
            type: range.type,
            endTime: range.endTime
          };
        }).filter(rangeInfo => rangeInfo.id == rangeId);
        if (ranges && ranges.length > 0) {
          this.setCustomBtns(ranges[0]);
        } else {
          this.setCustomBtns(null);
        }
      }
    }, () => {
      this.setCustomBtns(null);
    });
  }

  setCustomBtns(range) {
    this.customBtns = {
      iconBtns: []
    };
    if (range && (range.type != 'compete' || new Date().getTime() > +range.endTime * 1000)) {
      this.customBtns.iconBtns.push({
        tooltipInfo: '导出',
        placement: '',
        class: 'plx-ico-export-16',
        callback: this.openDlg.bind(this)
      });
    }
    this.customBtns.iconBtns.push({
      tooltipInfo: '刷新',
      placement: '',
      class: 'plx-ico-refresh-16',
      callback: this.refresh.bind(this)
    });
  }

  /**
   * 初始化导出弹出框中的表单项
   * @returns
   */
  initFormSetting() {
    let that = this;
    return {
      isShowHeader: false,
      isGroup: false,
      labelClass: 'col-sm-3',
      componentClass: 'col-sm-8',
      srcObj: {
        exportType: true
      },
      fieldSet: [
        {
          fields: [
            {
              name: 'exportType',
              label: '导出类型 ',
              placeholder: '请选择',
              type: FieldType.RADIO,
              callback: (values, $event, controls) => {
                console.log('select:' + values.exportType);
              },
              valueSet: [
                {
                  label: '详情',
                  value: true,
                },
                {
                  label: '摘要',
                  value: false
                }
              ],
            }
          ]
        }
      ]
    }
  }


  /**
   * 点击刷新按钮
   */
  refresh(): void {
    this.showLoading = true;
    this.scoreData = [];
    setTimeout(() => {
      this.showLoading = false;
      this.queryRangeScore();
    }, 500);
  }

  /**
   * 查询靶场成绩
   */
  queryRangeScore(): void {
    if (!this.rangeId || !this.selectedLanguage) {
      return;
    }
    let reqBody: any = this.buildRequestBody('query');
    this.manageService.queryScoreApi(this.rangeId, this.selectedLanguage, reqBody).subscribe(res => {
      this.scoreData = res.map(item => {
        return {
          department: item.department,
          teamName: item.teamName,
          userId: item.userId,
          userName: item.userName,
          hitNum: item.rangeScore.hitNum || 0, // 取到靶场数据
          hitScore: item.rangeScore.hitScore || 0,
          hitRate: item.rangeScore.hitRate || '0%',
          hitScoreHundredth: item.rangeScore.hitScoreHundredth || 0
        };
      });
    }, err => {
      this.plxMessage.error('获取打靶成绩失败！', err.msg);
    });
  }

  /**
   * 组装查询或导出的请求体
   * @param type 查询：query, 导出： export
   * @returns
   */
  buildRequestBody(type: string) {
    let departmentFilter: string;
    if (this.selectedDepartment.length == 0 || this.selectedDepartment.length == this.totalDptNum) {
      departmentFilter = 'all';
    } else {
      departmentFilter = this.selectedDepartment.join(',');
    }

    let params: any = {
      'department': departmentFilter,
      'verbose': type === 'query' ? false : this.verbose,
      'role': localStorage.getItem('role')
    };
    return params;
  }

  /**
   * 点击表格右上方的导出图标按钮
   * @returns
   */
  public openDlg(): any {
    if (this.scoreData.length == 0) {
      this.plxMessage.error('当前无成绩可导出，请更换查询条件后重试！', '');
      return;
    }
    if (!this.rangeId || !this.selectedLanguage) {
      return;
    }
    document.getElementById("exportDlgBtn").click();
  }

  /**
   * 打开对话框的事件函数
   * @param content
   */
  public open(content) {
    const size: 'sm' | 'lg' | 'xs' = 'sm';
    const options = {
      size: size,
      // contentClass: 'plx-modal-custom-content',
      destroyOnClose: true,
      modalId: 'plx-modal-1',
      openCallback: () => {
        console.log('open');
      }
    };
    this.modal = this.modalService.open(content, options);
  }

  /**
   * 点击导出对话框中的导出按钮
   * 导出当前查询条件下的靶场成绩
   */
  confirmExport(): void {
    this.verbose = this.formSetting.formObject.value.exportType;
    let reqBody: any = this.buildRequestBody('export');
    this.manageService.exportScoreApi(this.rangeId, this.selectedLanguage, reqBody).subscribe(res => {
      // var filename = "simDetail.xls";
      // var blob = new Blob([res], {type: "application/vnd.ms-excel"});
      // var objectUrl = URL.createObjectURL(blob);
      // window.open(objectUrl);

      const link = document.createElement('a');
      const blob = new Blob([res], { type: 'application/vnd.ms-excel' });
      link.setAttribute('href', window.URL.createObjectURL(blob));
      let fileName: string = this.rangeName;
      if (this.verbose) {
        fileName += '_detail';
      } else {
        fileName += '_summary'
      }
      link.setAttribute('download', fileName + '.xls');
      link.style.visibility = 'hidden';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    });
    this.modal.close();
  }

  /**
   * 点击导出对话框中的取消按钮
   */
  public cancel(): void {
    this.modal.close();
  }

  public ngOnDestroy() {
    this.modalService.destroyModalInstance();
  }

}

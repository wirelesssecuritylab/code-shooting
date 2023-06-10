import { Component, OnInit, ViewChild, TemplateRef } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { IOption, PlxMessage } from 'paletx';
import { ManageService } from 'src/app/admin/manage.service';
import { DefectService } from "../../shoot/services/defect-service";

const DEFAULT_WORKSPACE = 'public'
const DEFAULT_LANG_OPTION = { label: '所有', value: '所有' }

@Component({
  selector: 'app-free-practice',
  templateUrl: './freepractice.component.html',
  styleUrls: ['./freepractice.component.css']
})
export class FreePracticeComponent implements OnInit {
  @ViewChild('operTemplate', { static: true }) public operTemplate: TemplateRef<any>;

  public userId: string = "";
  public language: string = "所有";
  public workspace: string = DEFAULT_WORKSPACE;
  public customBtns: any;
  public columns = [];
  public targetData: any[] = [];
  public showLoading: boolean = false;
  public pageSizeSelections: number[] = [10, 30, 50];
  public targetLangOptions: Array<IOption> = [DEFAULT_LANG_OPTION,];
  public targetWorkspaceOptions: Array<IOption> = [];
  public mainCategoryOptions: Array<IOption> = [{ label: '所有', value: '所有' }];
  public subCategoryOptions: Array<IOption> = [{ label: '所有', value: '所有' }];
  public defectDetailOptions: Array<IOption> = [{ label: '所有', value: '所有' }];
  public mainCategory: string = '所有';
  public subCategory: string = '所有';
  public defectDetail: string = '所有';
  public isLoading: boolean;
  public defectResult: any;
  public isShowExtend: boolean = false;
  public targetIdFromUrl: string;

  constructor(private activatedRoute: ActivatedRoute, private manageService: ManageService,
    private plxMessageService: PlxMessage,
    private router: Router, private defectService: DefectService) { }

  ngOnInit(): void {
    this.userId = localStorage.getItem('user');
    this.setColumns();
    this.getWorkspaceList();
    this.getTargetList();
    this.getLangOptions()
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
        key: 'language',
        title: '所属语言',
        filter: false,
        show: true,
        width: '100px',
        sort: 'none'
      },
      {
        key: 'operation',
        title: '操作',
        show: true,
        fixed: true,
        class: 'plx-table-operation',
        width: '50px',
        contentType: 'template',
        template: this.operTemplate,
      }
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
   * 初始化高级查询的查询工作空间的语言列表
   * @returns
   */
  getLangOptions() {
    if (!this.workspace) {
      return;
    }
    this.manageService.getWorkSpaceLang(this.workspace).subscribe(res => {
      this.targetLangOptions = [DEFAULT_LANG_OPTION,]
      for (const key in res) {
        if (Object.prototype.hasOwnProperty.call(res, key)) {
          this.targetLangOptions.push({ label: key, value: key })
        }
      }
      // if (!this.targetLangOptions.some(item => item.value == this.language)) {
      //   // 如果新的选项中没有已选语言就置空。
      //   this.language = ''
      // }
    }, error => {
      this.mainCategoryOptions = [{ label: '所有', value: '所有' }];
      this.mainCategory = '所有';
      this.initSubBugTypeOptions(this.mainCategory);
      this.targetLangOptions = [DEFAULT_LANG_OPTION]
      // this.plxMessageService.error('获取靶子编码信息失败！', error.cause);
    })
  }
  /**
   * 初始化高级查询的缺陷大类下拉选项
   * @returns
   */
  initTagOptions() {
    let language = this.language != '所有' ? this.language : ''
    if (!language) {
      return;
    }
    this.manageService.getWorkSpaceLangDefect(this.workspace, language, "", true).subscribe(res => {
      if (res && res.result == 'success' && res.detail) {
        this.defectResult = res.detail;
        let options = [{ label: '所有', value: '所有' }];
        for (const key of Object.keys(res.detail)) {
          options.push({ label: key, value: key });
        }
        this.mainCategoryOptions = options;
      } else {
        this.mainCategoryOptions = [{ label: '所有', value: '所有' }];
      }
      this.mainCategory = '所有';
      this.initSubBugTypeOptions(this.mainCategory);
    }, error => {
      this.mainCategoryOptions = [{ label: '所有', value: '所有' }];
      this.mainCategory = '所有';
      this.initSubBugTypeOptions(this.mainCategory);
      // this.plxMessageService.error('获取靶子编码信息失败！', error.cause);
    })
  }

  /**
   * 切换缺陷大类时要重新初始化缺陷小类下拉项
   * @param bugType 缺陷大类的值
   */
  initSubBugTypeOptions(bugType: string) {
    this.subCategoryOptions = [{ label: '所有', value: '所有' }];;
    if (bugType && bugType !== '所有') {
      let bugTypeObj = this.defectResult[bugType];
      for (const key of Object.keys(bugTypeObj)) {
        this.subCategoryOptions.push({ label: key, value: key });
      }
    }
    this.subCategory = '所有';
    this.initBugDetailOptions(this.subCategory);
  }

  /**
   * 切换缺陷小类时要重新初始化缺陷细项下拉项
   * @param subBugType 缺陷小类的值
   */
  initBugDetailOptions(subBugType: any) {
    this.defectDetailOptions = [{ label: '所有', value: '所有' }];
    let bugType: string = this.mainCategory;;
    if (subBugType && subBugType !== '所有') {
      let subBugTypeObj = this.defectResult[bugType][subBugType];
      let options: any = subBugTypeObj.map(item => {
        return { label: item.description, value: item.description }
      });
      this.defectDetailOptions = this.defectDetailOptions.concat(options);
    }
    this.defectDetail = '所有';
  }
  /**
  * 查询指定工作空间的语言列表
  */
  getWorkspaceList() {
    this.manageService.getWorkspace().subscribe(
      res => {
        console.log('workspace list: ', res);
        this.targetWorkspaceOptions = []
        res.forEach(elem => {
          this.targetWorkspaceOptions.push({ label: elem.name, value: elem.id })
        });
      },
      error => {
        this.plxMessageService.error('获取工作空间列表失败！', error.cause);
      })
  }
  /**
   * 根据筛选条件查询靶子列表，默认查询所有语言的靶子，高级查询支持按标签过滤
   */
  getTargetList() {
    let language = this.language != '所有' ? this.language : ''

    let reqBody: any = {
      name: 'queryalreadytarget',
      //name: 'query',
      parameters: {
        workspace: this.workspace,
        language: language,
        user: this.userId,
        tagName: {
          mainCategory: this.mainCategory,
          subCategory: this.subCategory,
          defectDetail: this.defectDetail
        }
      }
    };
    this.manageService.manageTargetApi(reqBody).subscribe(res => {
      this.showLoading = false;
      if (res.result == 'success' && res.detail && Array.isArray(res.detail)) {
        let targetList = res.detail;
        let langOptions: any[] = [];
        let langs: string[] = [];


        this.targetData = targetList.map(item => {
          if (!language && langs.indexOf(item['language']) < 0) {
            langs.push(item['language']);
            langOptions.push({ label: item['language'], value: item.language });
          }
          return item;
        });
        // 语言下拉选项从靶子数据中提取而来
        // if (!this.language) {
        //   this.targetLangOptions = langOptions;
        // }
      } else {
        this.targetData = [];
      }
    }, error => {
      this.showLoading = false;
      if (!language) {
        this.targetLangOptions = [DEFAULT_LANG_OPTION,];
      }
      this.targetData = [];
      this.plxMessageService.error('获取靶子列表失败！', error.cause);
    });

  }

  getAllTargetList() {
    this.language = '所有'
    this.getTargetList();
  }
  /**
   * 点击刷新按钮
   */
  refresh(): void {
    this.showLoading = true;
    this.targetData = [];
    this.getTargetList();
  }

  /**
   * 点击高级查询按钮
   */
  advancedClick() {
    this.isShowExtend = !this.isShowExtend;
    if (!this.isShowExtend) {
      this.resetAdvanceFilter();
    }
  }

  /**
   * 重置标签过滤条件
   */
  resetAdvanceFilter() {
    this.mainCategory = '所有';
    this.subCategory = '所有';
    this.defectDetail = '所有';
    this.workspace = DEFAULT_WORKSPACE;
    this.language = '所有';
    this.getTargetList();
  }

  /**
   * 选择某个靶子进行打靶
   * @param rowData
   */
  gotoShoot(rowData: any) {
    this.defectService.targets = [];  // 自由练习时默认都是开始打靶，暂不支持保存草稿功能
    let routeStr: string = '../../shoot';
    this.router.navigate([routeStr], {
      queryParams: {
        'targetId': rowData.id,
        'targetName': rowData.name,
        'language': rowData.language,
        'shootType': 'startShoot',
        'from': 'test'
      },
      relativeTo: this.activatedRoute
    });
  }

  /**
   * 查看打靶成绩
   * @param rowData
   */
  queryScore(rowData: any) {
    let routeStr: string = '../../query';
    this.router.navigate([routeStr], {
      queryParams: {
        'targetId': rowData.id,
        'language': rowData.language
      },
      relativeTo: this.activatedRoute
    });
  }
}

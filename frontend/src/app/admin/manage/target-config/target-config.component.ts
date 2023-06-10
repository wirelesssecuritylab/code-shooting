import { Component, OnInit, ViewChild, Input, Output, EventEmitter } from '@angular/core';
import { languageAdapter } from "../../../shared/class/row";
import { ManageService } from '../../manage.service';


@Component({
  selector: 'app-target-config',
  templateUrl: './target-config.component.html',
  styleUrls: ['./target-config.component.css']
})
export class TargetConfigComponent implements OnInit {
  @ViewChild('langTemplate', {static: true}) private langTemplate: any;
  @ViewChild('instituteLabelTemplate', {static: true}) private instituteLabelTemplate: any;
  @ViewChild('targetTemplate', {static: true}) private targetTemplate: any;
  @ViewChild('operTemplate', {static: true}) private operTemplate: any;
  @Output() public dataChanged = new EventEmitter<any>();
  @Input() data: any[];
  public userId: string;
  public columns = [];
  public langOptions: any[] = [];
  public instituteLabelOptions:any[] = [];
  public allTargetList: any[] = [];

  constructor(private manageService: ManageService) { }

  ngOnInit(): void {
    this.userId = localStorage.getItem('user');
    this.instituteLabelOptions = [{
      label:'通用', value: '通用'
    }, {
      label: '严选', value:'严选'
    }];
    this.getAllTargets();
    this.setColumns();
  }

  /**
   * 初始化编辑表格语言列下拉选项
   */
  public initLangOptions() {
    let options: any[] = [];
    let langs: string[] = [];
    this.allTargetList.forEach(item => {
      if (langs.indexOf(item['language']) < 0) {
        langs.push(item['language']);
        options.push({label: item['language'], value: item.language});
      }
    });
    this.langOptions = options;
  }

  /**
   * 初始化某一行靶子列下拉选项
   * @param rowData
   */
  public initTargetList(rowData) {
    console.log(rowData);
    let filterList = JSON.parse(JSON.stringify(this.allTargetList));
    // 过滤某一语言下的靶子
    if (rowData.language) {
      filterList = filterList.filter(item => {
        return item.language == rowData.language;
      });
    }
    // 过滤含有院级标签的靶子
    if (rowData?.instituteLabel?.length > 0) {
      filterList = filterList.filter(item => {
        return (item.instituteLabel ||[]).join('') == rowData.instituteLabel.join('');
      });
    }
    this.data[rowData.plxTableDataIndex]['targetOptions'] = [];
    filterList.forEach(item => {
      this.data[rowData.plxTableDataIndex].targetOptions.push({
        'label': item.name,
        'value': item.id
      });
    });
    // 如果上一次选择的靶子不存在于新过滤后的靶子列表中，则把当前靶子重置为空
    let filterSameTarget = this.data[rowData.plxTableDataIndex].targetOptions.filter(item => {
      if (item.value == rowData.target) {
        return item;
      }
    });
    if (filterSameTarget.length < 1) {
      this.data[rowData.plxTableDataIndex].target = '';
    }
    this.emitTableData();
  }

  /**
   * 设置表格列
   */
  public setColumns() {
    this.columns = [
      {
        key: 'language',
        title: '语言',
        show: true,
        width: '160px',
        contentType: 'template',
        template: this.langTemplate,
      },
      {
        key: 'instituteLabel',
        title: '院级标签',
        show: true,
        width: '160px',
        contentType: 'template',
        template: this.instituteLabelTemplate,
      },
      {
        key: 'target',
        title: '靶子',
        show: true,
        width: '160px',
        contentType: 'template',
        template: this.targetTemplate,
      },
      {
        key: 'operate',
        title: '操作',
        show: true,
        width: '50px',
        contentType: 'template',
        template: this.operTemplate,
      },
    ];
  }

  /**
   * 查询当前用户可以使用的所有靶子
   * @param lang
   */
  public getAllTargets(lang?: string) {
    let reqBody:any = {
      name: 'query',
      parameters: {
        owner: this.userId
      }
    };
    if (lang) {
      reqBody.parameters['language'] = lang;
    }
    this.manageService.manageTargetApi(reqBody).subscribe(res => {
      if (res && res.detail && Array.isArray(res.detail)) {
        this.allTargetList = res.detail;
        this.initLangOptions();
      }
    });
  }

  /**
   * 新增一条靶子信息
   */
  public addData() {
    this.data = [...this.data, this.getRow()];
  }

  private getRow() {
    return {
      language: '',
      target: '',
    };
  }

  /**
   * 删除一条靶子信息
   * @param row
   */
  public deleteData(row: any) {
    this.data.splice(this.data.indexOf(row), 1);
    this.emitTableData();
  }

  public emitTableData() {
    this.dataChanged.emit(this.data);
  }
}

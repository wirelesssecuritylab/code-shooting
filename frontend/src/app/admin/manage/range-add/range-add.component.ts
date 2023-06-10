import { Component, NgModule, OnInit } from '@angular/core';
// import { Validators, AbstractControl, ValidatorFn} from '@angular/forms';
import { FieldType, PlxBreadcrumbItem, PlxMessage, PlxDateRangePickerComponent,DateTimePickerComponent } from 'paletx';
import { Router, ActivatedRoute } from '@angular/router';
import { TargetConfigComponent } from '../target-config/target-config.component';
import { ManageService } from '../../manage.service';

const _ADD = 'add';
const _EDIT = 'modify';

@Component({
  selector: 'app-range-add',
  templateUrl: './range-add.component.html',
  styleUrls: ['./range-add.component.css']
})
export class RangeAddComponent implements OnInit {
  formSetting: any;
  srcObj: any = {
    name: '',
    project: '',
    type: 'compete',
    startTime: '',
    endTime: '',
    desensitiveTime: '',
  };
  // 时间段组件的值
  setTimeRange: any = {startTime: '', endTime: ''};
  setTimeDesensitive: any = {desensitiveTime: ''}
  desensitiveTime: Date;
  // 新建靶场时，配置靶子可编辑表格必须默认有一条数据，即使为空
  defaultTarget: any[] = [{language: '', targetId: ''}];
  startTime: Date;
  endTime: Date;
  breadModel: PlxBreadcrumbItem[] = [];
  userId: string;
  userName: string;
  operateType: string = _ADD;
  rangeId: string;
  targetList: any[] = [];
  public projectOptions = [];
  public refreshDefValues: boolean = true;

  constructor(private router: Router, private manageService: ManageService,
              private plxMessageService: PlxMessage,
              private activatedRoute: ActivatedRoute) { }

  ngOnInit(): void {
    this.userId = localStorage.getItem('user');
    this.userName = localStorage.getItem('name');
    this.initProjectList();
    this.activatedRoute.queryParamMap.subscribe((paramMap) => {
      this.rangeId = paramMap.get('id');
    });
    let curRouterUrl = this.router.routerState.snapshot.url;
    if (curRouterUrl.indexOf("/main/manage/range/edit") >= 0) {
      this.operateType = _EDIT;
    }
    this.breadModel = [
      {label: '靶场管理', routerLink: '../list', name: 'rangeManage'}
    ];
    if (this.operateType == _EDIT) {
      this.breadModel.push({label: '编辑靶场', name: 'editLink'});
      this.getRangeDetail();
    } else {
      this.breadModel.push({label: '新建靶场', name: 'addLink'});
    }
    this.initFormSetting();
  }

  /**
   * 初始化项目下拉列表
   * @returns
   */
  initProjectList() {
    if (!this.userId) {
      this.projectOptions = [];
      return;
    }
    this.manageService.getProjectListApi(this.userId).subscribe(res => {
      if(res && res.projects && Array.isArray(res.projects)) {
        let options: any[] = [];
        res.projects.forEach(item => {
          options.push({label: item.name, value: item.id});
        });
        this.projectOptions = options;
        this.formSetting.fieldMap['project'].valueSet = this.projectOptions;
      }
    }, error => {
      this.plxMessageService.error('获取项目列表失败！', error.cause);
    });
  }

  /**
   * 获取单个靶子详情
   * @returns
   */
  getRangeDetail() {
    if (!this.rangeId) {
      return;
    }
    let reqBody: any = {
      name: 'query',
      parameters: {
        id: this.rangeId,
        owner: this.userId
      }
    };
    this.manageService.manageRangeApi(reqBody).subscribe(res => {
      this.initRangeFormInfo(res);
    }, error => {
      this.plxMessageService.error('获取靶场详情失败！', error.cause);
    });
  }

  /**
   * 初始化表单中非自定义组件的参数
   * @param rangeDetail 靶场详情
   * @returns
   */
  initRangeFormInfo(rangeDetail: any) {

    if (!rangeDetail || !rangeDetail.detail || !Array.isArray(rangeDetail.detail) || rangeDetail.detail.length != 1) {
      return;
    }
    let detail: any = rangeDetail.detail[0];
    Object.keys(this.srcObj).forEach((key) => {
      if (key == 'startTime') {
        if (detail.type == 'compete') {
          this.startTime = new Date(parseInt(detail['startTime'])*1000);
          this.setTimeRange.startTime = this.startTime;
          this.srcObj['timeRange'] = this.setTimeRange;
          this.formSetting.fieldMap['timeRange'].binding.startDate = this.startTime;
        }
      } else if (key == 'endTime') {
        if (detail.type == 'compete') {
          this.endTime = new Date(parseInt(detail['endTime'])*1000);
          this.setTimeRange.endTime = this.endTime;
          this.srcObj['timeRange'] = this.setTimeRange;
          this.formSetting.fieldMap['timeRange'].binding.endDate = this.endTime;
        }
      } else if (key == "desensitiveTime"){
         if(detail.desensitiveTime != 0 && detail.desensitiveTime != -62135596800){
          this.desensitiveTime = new Date(parseInt(detail['desensitiveTime'])*1000);
          this.setTimeDesensitive.desensitiveTime = this.desensitiveTime;
          this.srcObj['timeDesensitive'] = this.setTimeDesensitive;
          this.formSetting.fieldMap['timeDesensitive'].binding.dateValue = this.desensitiveTime;
         }
      }else {
        this.srcObj[key] = detail[key];
      }
    });
    this.refreshDefValues = !this.refreshDefValues;
    // 手动设置表单项targets的值
    this.formSetting.formObject.controls['targets'].setValue(detail.targets);
    this.formSetting.formObject.controls['targets'].markAsDirty();
    this.getAllTargets(detail);
  }

  /**
   * 获取当前用户能访问的所有靶子列表,主要是用来初始化靶子配置可编辑表格中靶子的下拉选项
   * @param rangeDetail 靶场详情
   * @param lang 语言
   */
  public getAllTargets(rangeDetail: any, lang?: string) {
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
        let allTargetList = res.detail;
        this.targetList = [];
        // 再设置靶子配置表单项的输入参数
        let inputTargets: any[] = (rangeDetail.targets || []).map(item => {
          let language = item.language;
          let filterList = allTargetList.filter(item => {
            return item.language == language;
          });
          item['target'] = item.targetId
          item['targetOptions'] = [];
          filterList.forEach(filterItem => {
            item['targetOptions'].push({
              'label': filterItem.name,
              'value': filterItem.id
            });
          });
          this.targetList.push({targetId: item.targetId});
          return item;
        });
        // 设置targets表单项(自定义组件绑定的输入属性的值, 以用来初始化可编辑表格中的数据)
        this.formSetting.fieldMap['targets'].binding.data = inputTargets;
        this.refreshDefValues = !this.refreshDefValues;
      }
    });
  }

  /**
   * 初始化表单项
   */
  initFormSetting(): void {
    this.formSetting = {
      // size: 'lg',
      isShowHeader: false,
      header: '新建靶场',
      isGroup: true,
      hideGroup: false,
      srcObj: this.srcObj,
      advandedFlag: false,
      // 按钮
      buttons: [
        {
          type: 'submit',
          label: '确定',
          class: 'plx-btn plx-btn-primary plx-btn-sm',
          hidden: false,
          disabled: false,
          callback: (values, $event, controls) => {
            this.addRange(values, $event, controls);
          }
        },
        {
          type: 'cancel',
          label: '取消',
          class: 'plx-btn plx-btn-sm',
          hidden: false,
          disabled: false,
          callback: (values, $event, controls) => {
            // alert('你点了取消？data=' + JSON.stringify(values));
            this.router.navigateByUrl('/main/manage/range/list');
          }
        }
      ],
      // 表单项
      fieldSet: [
        {
          group: '基本信息',
          fields: [
            {
              name: 'name',
              label: '名称',
              type: FieldType.STRING,
              // desc: '支持字母、数字、"_"、"-"的组合，4-20个字符',
              required: true,
              disabled: false,
              // validators: [Validators.pattern(/^[a-zA-Z0-9_-]{4,20}$/), Validators.minLength(4), Validators.maxLength(20)],
              // validateInfos: {
              //   pattern: '名称只能输入字母、数字、"_"、"-"的组合，4-20个字符'
              // }
            },
            {
              name: 'project',
              label: '所属组织',
              type: FieldType.SELECTOR,
              required: true,
              multiple: false,
              valueSet: this.projectOptions
            },
            {
              name: 'type',
              label: '打靶类型',
              type: FieldType.SELECTOR,
              required: true,
              multiple: false,
              valueSet: [
                {
                  label: '练习',
                  value: 'test'
                },
                {
                  label: '比赛',
                  value: 'compete'
                },
              ],
            },
            {
              name: 'timeRange',
              label: '起止时间',
              required: true,
              type: FieldType.CUST_COMPONENT,
              component: PlxDateRangePickerComponent,
              componentClass: 'col-sm-5',
              hidden: (values: any, control: any) => {
                return values.type === 'test';
              },
              binding: {
                showTime: true,
                showSeconds: true,
                dateFormat: 'YYYY-MM-DD HH:mm:ss',
                placeHolderStartDate: '请选择开始时间',
                placeHolderEndDate: '请选择结束时间',
                startDate: this.startTime,
                endDate:this.endTime
              },
              outputs: {
                onStartDateClosed: (event) => {
                  // 由于 是自定义组件，所以表单项的值要手动设置一下，然后提交时解析所用
                  this.setTimeRange.startTime = event.date;
                  this.formSetting.formObject.controls['timeRange'].setValue(this.setTimeRange);
                  this.formSetting.formObject.controls['timeRange'].markAsDirty();
                },
                onEndDateClosed: (event) => {
                  this.setTimeRange.endTime = event.date;
                  this.formSetting.formObject.controls['timeRange'].setValue(this.setTimeRange);
                  this.formSetting.formObject.controls['timeRange'].markAsDirty();
                }
              }
            },
            {
              name: 'timeDesensitive',
              label: '靶标脱敏时间',
              required: false,
              type: FieldType.CUST_COMPONENT,
              component: DateTimePickerComponent,
              componentClass: 'col-sm-5',
              binding: {
                showTime: true,
                showSeconds: true,
                dateFormat: 'YYYY-MM-DD HH:mm:ss',
                placeHolder: '请选择靶标脱敏时间',
                dateValue: this.desensitiveTime
              },
              outputs:{
                onConfirm: (event) =>{
                  this.setTimeDesensitive.desensitiveTime = event.date;
                  this.formSetting.formObject.controls['timeDesensitive'].setValue(this.setTimeDesensitive);
                  this.formSetting.formObject.controls['timeDesensitive'].markAsDirty();
                }
              }
              
            }
          ]
        },
        {
          group: '配置靶子',
          fields: [
            {
              name: 'targets',
              label: '靶子',
              required: true,
              type: FieldType.CUST_COMPONENT,
              component: TargetConfigComponent,
              binding: {
                data: this.defaultTarget,
              },
              outputs: {
                dataChanged:(event) => {
                  // event是string类型的数组
                  this.buildTargetList(event);
                  if (!event || event.length == 0) {
                    this.formSetting.formObject.controls['targets'].setValue('');
                    this.formSetting.formObject.controls['targets'].markAsDirty();
                  }
                }
              }
            }
          ]
        }
      ]
    };
  }

  /**
   * 把表格数据组装成下发时的结构体targetList
   * @param tableData 可编辑表格绑定的数据
   */
  buildTargetList(tableData: any) {
    this.targetList = [];
    (tableData || []).forEach(data => {
      if (data.target) {
        this.targetList.push({targetId: data.target});
      }
    });
    this.formSetting.formObject.controls['targets'].setValue(this.targetList);
    this.formSetting.formObject.controls['targets'].markAsDirty();
  }

  /**
   * 创建或编辑靶场
   * @param values
   * @param $event
   * @param controls
   */
  addRange(values: any, $event: any, controls): void {
    
    
    if (!this.checkFormValid(values)) {
      return;
    }
    let reqBody = this.buildReqBody(values);
    let sucMsg: string = this.operateType == _EDIT ? '修改成功' : '创建成功！';
    let errorMsg: string = this.operateType == _EDIT ? '修改失败' : '创建失败！';
    this.manageService.manageRangeApi(reqBody).subscribe(res => {
      this.plxMessageService.success(sucMsg, '');
      this.router.navigateByUrl('/main/manage/range/list');
    }, err => {
      this.plxMessageService.error(errorMsg, err.cause);
    });
  }

  /**
   * 提交时验证表单的合法性
   * @param formValues 表单数据
   * @returns
   */
  checkFormValid(formValues: any) {
    if (formValues.type == 'compete') {
      if (!formValues.timeRange.startTime) {
        this.plxMessageService.error('请设置开始时间！', '');
        return false;
      }
      if (!formValues.timeRange.endTime) {
        this.plxMessageService.error('请设置结束时间！', '');
        return false;
      }
      let startTimeStamp: number = new Date(formValues.timeRange.startTime).getTime()/1000;
      let endTimeStamp: number = new Date(formValues.timeRange.endTime).getTime()/1000;
      if (endTimeStamp <= startTimeStamp) {
        this.plxMessageService.error('结束时间必须晚于开始时间，请修改！', '');
        return false;
	}
      let desensitiveTimeStamp: number = new Date(formValues.timeDesensitive.desensitiveTime).getTime() / 1000;
      if (typeof (formValues.timeDesensitive.desensitiveTime) != "undefined") {
        if (desensitiveTimeStamp < endTimeStamp) {
          this.plxMessageService.error('靶标脱敏时间必须大于等于比赛结束时间，请修改！', '');
          return false;
        }
      }
    }
    if (this.targetList.length == 0) {
      this.plxMessageService.error('请至少配置一个靶子', '');
      return false;
    }
    // 判断是否有重复的靶子
    let targetId = [];
    for(let i = 0; i < this.targetList.length; i++) {
      if(targetId.indexOf(this.targetList[i].targetId) < 0) {
        targetId.push(this.targetList[i].targetId);
      } else {
        break;
      }
    }
    if (targetId.length < this.targetList.length) {
      this.plxMessageService.error('不能配置重复的靶子', '');
      return false;
    }
    return true;
  }


  /**
   * 构造新建或编辑靶场时的请求体
   * @param formValues 表单数据
   * @returns
   */
  buildReqBody(formValues:any) {
    let addRangeParam = {
      name: formValues.name,
      project: formValues.project,
      type: formValues.type,
      owner: this.userId,
      ownerName: this.userName,
      targets: this.targetList,
      desensitiveTime: new Date(formValues.timeDesensitive.desensitiveTime).getTime()/1000

    };
    if (formValues.type == 'compete') {
      addRangeParam['startTime'] = new Date(formValues.timeRange.startTime).getTime()/1000;
      addRangeParam['endTime'] = new Date(formValues.timeRange.endTime).getTime()/1000;
    }
    let editTargetParam: any;
    if (this.operateType == _EDIT) {
      editTargetParam = JSON.parse(JSON.stringify(addRangeParam));
      editTargetParam['id'] = this.rangeId;
    }
    let reqBody: any = {
      name: this.operateType == _EDIT ? _EDIT : _ADD,
      parameters: this.operateType == _EDIT ? editTargetParam : addRangeParam
    };
    return reqBody;
  }
}

import { Component, OnInit, ViewChild } from '@angular/core';
import {
  FieldType,
  PlxFileUpLoaderComponent,
  PlxBreadcrumbItem,
  FileItem,
} from 'paletx';
import { Router, ActivatedRoute } from '@angular/router';
import { PlxMessage } from 'paletx';
import { AddTargetBody, TargetFormField } from './target-param.class';
import { ManageService } from '../../manage.service';
import { TargetTagComponent } from './../target-tag/target-tag.component';
import { Validators } from '@angular/forms';
import { errCodeAdapter } from "../../../shared/class/row";
import { error } from '@angular/compiler/src/util';

// 编辑时用这个url上传
const _EDITFILEURL = '/api/code-shooting/v1/actions/uploadTarget/target/';
// 新建时用这个url进行上传
const _ADDFILEURL = '/api/code-shooting/v1/actions/uploadTemp/owner/';
const _TARGET = 'target';
const _ANSWER = 'answer';
const _ADD = 'add';
const _EDIT = 'modify';
const _EXPERTPRIVILEGE = 'editTargetTag';
const DEFAULTTAG: any[] = [{ label: '所有', value: '所有' }];

@Component({
  selector: 'app-target-add',
  templateUrl: './target-add.component.html',
  styleUrls: ['./target-add.component.css'],
})
export class TargetAddComponent implements OnInit {
  @ViewChild('plxDragUploader') public plxDragUploader: any;
  @ViewChild('plxSingleDragUploader') public plxSingleDragUploader: any;

  public userId: string;
  public userName: string;
  public targetId: string;
  public ownerName: string;
  public owner: string;
  public privilege: string;
  public breadModel: PlxBreadcrumbItem[] = [];
  public refreshDefValues: boolean = true;
  public formSetting: any;
  // 用以设置表单的默认参数值，单向绑定
  public srcObj: TargetFormField = new TargetFormField();
  // public defaultTagOption: string[] = [];// 默认的靶子标签（主要是编辑的时候用:树形下样面板时用）

  // 可拖拽上传组件中拖拽区域的提示文本
  public dragTipsAnswer = {
    desc: '仅支持.xlsm格式的文件。',
  };
  public dragTipsTarget = {
    desc: '支持go, py, java, c, cpp, js, ts, html, md等多种格式。',
  };
  // 上传组件允许上传的文件类型
  public allowedFileType = [
    'go',
    'py',
    'java',
    'c',
    'cpp',
    'js',
    'ts',
    'css',
    'html',
    'xlsx',
    'xls',
    'xlsm',
    'md',
    'h',
    'json',
  ];
  public allowedAnswerType = ['xlsm'];
  public targetAnswerUploaderConfig: any;
  public targetUploaderConfig: any;
  public targetAnswer: string = '';
  public targetList: string[] = [];
  public targetQueue: any = [];
  public answerQueue: any = [];
  public operateType: string = _ADD;

  public defectDetail: any;
  public bugTypeOptions: any[] = DEFAULTTAG;
  public subBugTypeOptions: any[] = DEFAULTTAG;
  public bugDetailOptions: any[] = DEFAULTTAG;

  public defaultUploadUrl: string;
  // 根据工作区间加载语言类型
  public langOptions: any[];
  // 根据语言选择规范
  public templateOptions: any[];


  wsMap: Map<string, string> = new Map();

  constructor(
    private router: Router,
    private plxMessageService: PlxMessage,
    private manageService: ManageService,
    private activatedRoute: ActivatedRoute
  ) { }

  ngOnInit(): void {
    this.userId = localStorage.getItem('user');
    this.userName = localStorage.getItem('name');
    this.privilege = localStorage.getItem('expertPrivilege');
    this.activatedRoute.queryParamMap.subscribe((paramMap) => {
      this.targetId = paramMap.get('id');
      this.ownerName = paramMap.get('ownerName');
      this.owner = paramMap.get('owner');
    });
    let curRouterUrl = this.router.routerState.snapshot.url;
    if (curRouterUrl.indexOf('/main/manage/target/edit') >= 0) {
      this.operateType = _EDIT;
    }
    this.breadModel = [
      {
        label: '靶子管理',
        routerLink: '/main/manage/target/list',
        name: 'rangeManage',
      },
    ];
    if (this.operateType == _EDIT) {
      this.defaultUploadUrl = _EDITFILEURL + this.targetId + '/type/';
      this.breadModel.push({ label: '编辑靶子', name: 'link2' });
      this.getTargetDetail();
    } else {
      this.defaultUploadUrl = _ADDFILEURL + this.userId + '/type/';
      this.breadModel.push({ label: '新建靶子', name: 'link2' });
    }
    this.initFileUploaderConfig();
    this.initFormSetting();
    this.loadWorkSpaces();
    if (this.operateType !== _EDIT && this.privilege === _EXPERTPRIVILEGE && this.formSetting.fieldSet) {
      let fieldSet = this.formSetting.fieldSet.find(item => item.group === '标签信息');
      if (fieldSet && fieldSet.fields) {
        let instituteLabel = fieldSet.fields.find(field => field.name === 'instituteLabel');
        if (instituteLabel) {
          instituteLabel.disabled = false;
        }
      }
    }
  }

  /**
   * 获取某一语言下的模板信息(即规范), 以用来初始化标签信息中的各下拉项
   * @param lang 所属语言
   * @param detail 靶子详情
   * @returns
   */
  initTagOptions(lang: string, detail?: any) {
    if (!lang) {
      return;
    }

    this.manageService.getWorkSpaceLangDefect(detail ? detail.workspace : "public", lang, detail.template, true).subscribe(
      (res) => {
        if (res && res.result == 'success' && res.detail) {
          this.bugTypeOptions = [{ label: '所有', value: '所有' }];
          this.defectDetail = res.detail;
          for (const key of Object.keys(this.defectDetail)) {
            this.bugTypeOptions.push({ label: key, value: key });
          }
          this.formSetting.fieldMap['bugType'].valueSet = this.bugTypeOptions;
          // 如果是编辑状态设置缺陷大类、缺陷小类和缺陷细项的值
          if (detail) {
            this.srcObj['bugType'] = detail['tagName']['mainCategory'];
            this.srcObj['subBugType'] = detail['tagName']['subCategory'];
            this.srcObj['bugDetail'] = detail['tagName']['defectDetail'];
            this.refreshDefValues = !this.refreshDefValues;
          }
        }
        this.setFormFieldStatus(detail);
      },
      (error) => {
        this.bugTypeOptions = [{ label: '所有', value: '所有' }];
        this.formSetting.fieldMap['bugType'].valueSet = this.bugTypeOptions;
        this.formSetting.formObject.controls['bugType'].setValue('所有');
        this.formSetting.formObject.controls['bugType'].markAsDirty();
        this.plxMessageService.error('获取靶子编码信息失败！', error.cause);
        this.setFormFieldStatus(detail);
      }
    );
  }

  /**
   * 编辑状态下设置表单中相关参数字段的可编辑状态：
   * 如果是targetchooser角色的人员，则不能修改其他人创建的靶子除院级标签外的参数； 只能修改院级标签外
   * @param detail 
   */
  setFormFieldStatus(detail) {
    let owner = detail.owner;
    // 非靶子所有者即targetchooser角色的人员，不能修改院级标签外的其他属性
    if (owner !== this.userId) {
      let fieldSet: string[] = ['name', 'workspace', 'template', 'language', 'isShared', 'bugType', 'subBugType', 'bugDetail', 'extendedLabel', 'customLable'];
      fieldSet.forEach((item) => {
        this.formSetting.fieldMap[item].disabled = true;
      });
      (<HTMLElement>document.getElementById('answer').children[0]).style.cssText = 'opacity: 0.5; pointer-events: none';
      (<HTMLElement>document.getElementById('targets').children[0]).style.cssText = 'opacity: 0.5; pointer-events: none';
    }
    // 如果是具有targetchooser角色的人员，则可以修改院级标签
    if (this.privilege === _EXPERTPRIVILEGE) {
      this.formSetting.fieldMap['instituteLabel'].disabled = false;
    }
  }
  /**
   * 初始化表单项
   */
  initFormSetting(): void {
    this.formSetting = {
      // size: 'lg',
      isShowHeader: false,
      header: '创建靶子',
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
            this.submitTarget(values, $event, controls);
          },
        },
        {
          type: 'cancel',
          label: '取消',
          class: 'plx-btn plx-btn-sm',
          hidden: false,
          disabled: false,
          callback: (values, $event, controls) => {
            this.router.navigateByUrl('/main/manage/target/list');
          },
        },
      ],
      // 表单项
      fieldSet: [
        {
          group: '基础信息',
          fields: [
            {
              name: 'name',
              label: '靶子名称',
              type: FieldType.STRING,
              // desc: '支持字母、数字、"_"、"-"的组合，4-20个字符',
              required: true,
              disabled: false,
              // validators: [Validators.pattern(/^[a-zA-Z0-9_-]{4,20}$/), Validators.minLength(4), Validators.maxLength(20)],
              // validateInfos: {
              //   pattern: '名称只能输入字母、数字、"_"、"-"的组合，4-20个字符'
              // }
              callback: (values, $event, controls) => {
                if (this.operateType == _ADD) {
                  this.checkTargetName(values);
                }
              },
            },
            {
              name: 'workspace',
              label: '工作空间',
              type: FieldType.SELECTOR,
              required: true,
              disabled: false,
              valueSet: [],
              callback: (values, $event, controls) => {
                this.onSelectWorkSpace(values)
              },
            },

            {
              name: 'template',
              label: '规范版本',
              type: FieldType.SELECTOR,
              required: true,
              multiple: false,
              // disabled: this.operateType == _EDIT,
              valueSet: this.templateOptions,
              callback: (values, $event, controls) => {
                this.UpdateLanguageOPtions(values.workspace, values.template, values.language);
              }
            },
            {
              name: 'language',
              label: '所属语言',
              type: FieldType.SELECTOR,
              required: true,
              multiple: false,
              // disabled: this.operateType == _EDIT,
              valueSet: this.langOptions,
              callback: (values, $event, controls) => {
                this.UpdateTagOptions(values.template, values.language, values.workspace);
              },
            },
            // 自定义树形节点下拉面板组件的使用案例
            // {
            //   name: 'targetTag',
            //   label: '靶子标签',
            //   type: FieldType.CUST_COMPONENT,
            //   component: TargetTagComponent,
            //   required: false,
            //   binding: {
            //     // selectedOptions: ['05', '06'],
            //     // selectedNodeKey: "'05', '06'",
            // selectedOptions: this.defaultTagOption,
            // selectedNodeKey: this.defaultTagOption.join(',')
            //   },
            //   outputs: {
            //     tagOptionChanged:(event) => {
            //       // event是string类型的数组
            //       console.log("选项变化 了！");
            //       this.formSetting.formObject.controls['targetTag'].setValue(event);
            //       this.formSetting.formObject.controls['targetTag'].markAsDirty();
            //     }
            //   }
            // },
            {
              name: 'isShared',
              label: '是否共享',
              type: FieldType.SWITCH,
              required: true,
              callback: (values, $event, controls) => {
              },
            },
          ],
        },
        {
          group: '标签信息',
          fields: [
            {
              name: 'bugType',
              label: '缺陷大类',
              type: FieldType.SELECTOR,
              required: false,
              multiple: false,
              valueSet: this.bugTypeOptions,
              callback: (values, $event, controls) => {
                this.initSubBugTypeOptions(values.bugType);
              },
              selected: (event) => { },
            },
            {
              name: 'subBugType',
              label: '缺陷小类',
              type: FieldType.SELECTOR,
              required: false,
              multiple: false,
              valueSet: this.subBugTypeOptions,
              callback: (values, $event, controls) => {
                console.log('twst');
                this.initBugDetailOptions(values);
              },
            },
            {
              name: 'bugDetail',
              label: '缺陷细项',
              type: FieldType.SELECTOR,
              required: false,
              multiple: false,
              valueSet: this.bugDetailOptions,
            },
            {
              name: 'extendedLabel',
              label: '扩展标签',
              type: FieldType.SELECTOR,
              required: false,
              multiple: true,
              valueSet: [
                {
                  label: '通用',
                  value: '通用',
                },
                {
                  label: '严选',
                  value: '严选',
                },
                {
                  label: '外场故障',
                  value: '外场故障',
                },
                {
                  label: '内部故障',
                  value: '内部故障',
                },
                {
                  label: '学习',
                  value: '学习',
                },
              ],
            },
            {
              name: 'customLable',
              label: '自定义标签',
              type: FieldType.STRING,
              required: false,
              validators: [Validators.maxLength(100)],
            },
            {
              name: 'instituteLabel',
              label: '院级标签',
              type: FieldType.SELECTOR,
              disabled: true, // 默认不能修改
              required: false,
              multiple: true,
              valueSet: [
                {
                  label: '通用',
                  value: '通用',
                },
                {
                  label: '严选',
                  value: '严选',
                },
              ],
            },
          ],
        },
        {
          group: '核心信息',
          fields: [
            {
              name: 'answer',
              label: '上传靶标',
              type: FieldType.CUST_COMPONENT,
              required: true,
              component: PlxFileUpLoaderComponent,
              binding: {
                pxFileUploaderConfig: this.targetAnswerUploaderConfig,
                viewQueue: this.answerQueue,
                isDrag: true,
                dragTips: this.dragTipsAnswer,
              },
            },
            {
              name: 'targets',
              label: '上传靶子',
              type: FieldType.CUST_COMPONENT,
              required: true,
              component: PlxFileUpLoaderComponent,
              binding: {
                pxFileUploaderConfig: this.targetUploaderConfig,
                isDrag: true,
                dragTips: this.dragTipsTarget,
                viewQueue: this.targetQueue,
              },
            },
          ],
        },
      ],
    };
  }

  /**
   * 切换缺陷大类时要重新初始化缺陷小类下拉项
   * @param bugType 缺陷大类的值
   */
  initSubBugTypeOptions(bugType: string) {
    this.subBugTypeOptions = [{ label: '所有', value: '所有' }];
    if (bugType && bugType !== '所有') {
      let bugTypeObj = this.defectDetail[bugType];
      for (const key of Object.keys(bugTypeObj)) {
        this.subBugTypeOptions.push({ label: key, value: key });
      }
    }
    this.formSetting.fieldMap['subBugType'].valueSet = this.subBugTypeOptions;
    this.formSetting.formObject.controls['subBugType'].setValue('所有');
    this.formSetting.formObject.controls['subBugType'].markAsDirty();
  }

  /**
   * 切换缺陷小类时要重新初始化缺陷细项下拉项
   * @param subBugType
   */
  initBugDetailOptions(formValue: any) {
    this.bugDetailOptions = [{ label: '所有', value: '所有' }];
    let bugType: string = formValue.bugType;
    let subBugType: string = formValue.subBugType;
    if (subBugType && subBugType !== '所有') {
      let subBugTypeObj = this.defectDetail[bugType][subBugType];
      let options: any = subBugTypeObj.map((item) => {
        return { label: item.description, value: item.description };
      });
      this.bugDetailOptions = this.bugDetailOptions.concat(options);
    }
    this.formSetting.fieldMap['bugDetail'].valueSet = this.bugDetailOptions;
    this.formSetting.formObject.controls['bugDetail'].setValue('所有');
    this.formSetting.formObject.controls['bugDetail'].markAsDirty();
  }

  /**
   * 初始化文件上传组件配置文件
   */
  initFileUploaderConfig() {
    // 上传靶标
    this.targetAnswerUploaderConfig = {
      defaultUploadUrl: this.defaultUploadUrl + _ANSWER,
      showCustomUrl: false,
      useBrowserProgress: true,
      autoUpload: true,
      showDelIcon: true,
      isCustom: true,
      queueLimit: 1,
      allowedFileType: this.allowedAnswerType,
      fileListConfig: {
        showFileName: true,
        showDetail: false,
        showTooltip: true,
      },
      // 文件验证不通过时的回调
      onFilterFile: (filter: { name: string }) => {
        let msg = '';
        if (filter.name === 'queueLimit') {
          msg = '只能上传一个靶标文件，请先删除再重新上传';
        }
        return msg;
      },
      // 当文件上传成功时的回调
      onSuccessItem: (item) => {
        this.plxMessageService.success(item.file.name + '上传成功', '');
        this.targetAnswer = item.file.name;
        this.formSetting.formObject.controls['answer'].setValue(
          this.targetAnswer
        );
        this.formSetting.formObject.controls['answer'].markAsDirty();
      },
      // 当文件上传错误时的回调
      onErrorItem: (item) => {
        this.plxMessageService.error(item.file.name + '上传失败', '');
      },
      // 上传成功的文件是否能删除
      onDelItem: (item) => {
        this.targetAnswer = '';
        this.formSetting.formObject.controls['answer'].setValue(
          this.targetAnswer
        );
        this.formSetting.formObject.controls['answer'].markAsDirty();
        return Promise.resolve(true);
      },
      // 上传失败或未上传的文件是否能删除
      onRemoveItem: (item) => {
        return Promise.resolve(true);
      },
      onDelViewItem: (item) => {
        this.targetAnswer = '';
        this.formSetting.formObject.controls['answer'].setValue(
          this.targetAnswer
        );
        this.formSetting.formObject.controls['answer'].markAsDirty();
        return Promise.resolve(true);
      },
      // 文件选择完毕后的回调
      onAfterAddingAll: (items: FileItem[]) => {
        return new Promise((resolve) =>
          resolve(
            this.isAllowFile(
              items,
              this.allowedAnswerType,
              this.plxDragUploader
            )
          )
        );
      },
    };

    // 可拖拽多个靶子文件上传配置
    this.targetUploaderConfig = {
      defaultUploadUrl: this.defaultUploadUrl + _TARGET,
      useBrowserProgress: true,
      autoUpload: true,
      showDelIcon: true,
      isCustom: true,
      // allowedFileType: this.allowedFileType,
      fileListConfig: {
        showFileName: true,
        showDetail: false,
        showTooltip: true,
      },
      // 当文件上传成功时的回调
      onSuccessItem: (item) => {
        this.plxMessageService.success(item.file.name + '上传成功', '');
        if (this.targetList.indexOf(item.file.name) < 0) {
          this.targetList.push(item.file.name);
        }
        // 上传成功后手动设置一下该表单字段的值，否则必填项校验不通过
        this.formSetting.formObject.controls['targets'].setValue(
          this.targetList.join(',')
        );
        this.formSetting.formObject.controls['targets'].markAsDirty();
      },
      // 当文件上传错误时的回调
      onErrorItem: (item) => {
        this.plxMessageService.error(item.file.name + '上传失败', '');
      },
      // 上传成功的文件是否能删除
      onDelItem: (item) => {
        let index = this.targetList.indexOf(item.file.name);
        if (index >= 0) {
          this.targetList.splice(index, 1);
        }
        this.formSetting.formObject.controls['targets'].setValue(
          this.targetList.join(',')
        );
        this.formSetting.formObject.controls['targets'].markAsDirty();
        return Promise.resolve(true);
      },
      // 上传失败或未上传的文件是否能删除
      onRemoveItem: (item) => {
        return Promise.resolve(true);
      },
      // 查看模式下删除文件时的回调，
      onDelViewItem: (item) => {
        let index = this.targetList.indexOf(item.fileName);
        if (index >= 0) {
          this.targetList.splice(index, 1);
        }
        this.formSetting.formObject.controls['targets'].setValue(
          this.targetList.join(',')
        );
        this.formSetting.formObject.controls['targets'].markAsDirty();
        return Promise.resolve(true);
      },
      // 文件选择完毕后的回调
      onAfterAddingAll: (items: FileItem[]) => {
        return Promise.resolve(true);
        // return new Promise(resolve => resolve(this.isAllowFile(items, this.allowedFileType, this.plxDragUploader)));
      },
    };
  }

  /**
   * 查询某个特定id的靶子信息
   */
  getTargetDetail() {
    if (!this.targetId) {
      return;
    }
    let reqBody: any = {
      name: 'query',
      parameters: {
        id: this.targetId,
        owner: this.userId,
      },
    };
    this.manageService.manageTargetApi(reqBody).subscribe(
      (res) => {
        this.initTargetFormInfo(res);
      },
      (error) => {
        this.plxMessageService.error('获取靶子详情失败！', error.cause);
      }
    );
  }

  /**
   * 编辑状态下初始化当前待编辑靶子的表单信息
   * @param targetDetail 靶子详情信息
   */
  initTargetFormInfo(targetDetail: any) {
    if (
      !targetDetail ||
      !targetDetail.detail ||
      !Array.isArray(targetDetail.detail) ||
      targetDetail.detail.length != 1
    ) {
      return;
    }
    let detail: any = targetDetail.detail[0];
    // 编辑状态下给表单中各元素设置初值(标签信息在defect接口响应后再赋值)
    if (detail) {
      Object.keys(this.srcObj).forEach((key) => {
        if (['bugType', 'subBugType', 'bugDetail'].indexOf(key) < 0) {
          if (key == 'targets') {
            this.srcObj[key] = detail[key].join(',');
          } else {
            this.srcObj[key] = detail[key];
          }
        }
      });
      // 设置上传组件查看模式下的 文件列表
      if (detail.answer) {
        this.answerQueue = [{ fileName: detail.answer }];
        this.targetAnswer = detail.answer;
      }
      this.formSetting.fieldMap['answer'].binding.viewQueue = this.answerQueue;
      if (detail.targets && Array.isArray(detail.targets)) {
        this.targetList = detail.targets;
        detail.targets.forEach((target) => {
          this.targetQueue.push({ fileName: target });
        });
      }
      this.formSetting.fieldMap['targets'].binding.viewQueue = this.targetQueue;
      this.refreshDefValues = !this.refreshDefValues;
    } else {
      // 以下setValue也是设置初值的一种方式
      this.formSetting.formObject.controls['bugType'].setValue('所有');
      this.formSetting.formObject.controls['subBugType'].setValue('所有');
      this.formSetting.formObject.controls['bugDetail'].setValue('所有');
      this.formSetting.formObject.controls['bugType'].markAsDirty();
      this.formSetting.formObject.controls['subBugType'].markAsDirty();
      this.formSetting.formObject.controls['bugDetail'].markAsDirty();
    }
    // 先获取当前靶子所属语言的标签信息后再赋值
    this.initTagOptions(detail['language'], detail);
  }


  /**
   * 获取靶子列表
   */
  checkTargetName(values: any): void {
    let reqBody: any = {
      name: 'query',
      parameters: {
        owner: this.userId
      }
    };
    reqBody.parameters['name'] = values.name;

    this.manageService.manageTargetApi(reqBody).subscribe(res => {

      if (res == true) {
        this.plxMessageService.error('靶子名称重复，请修改靶子名称！', "");
      }

    }, error => {
      this.plxMessageService.error('获取靶子名称失败！', error.cause);
    });
  }



  /**
   * 创建或编辑靶子
   * @param values
   * @param $event
   * @param controls
   */
  submitTarget(values: any, $event: any, controls): void {
    console.log(values)
    let addTargetParam: AddTargetBody = new AddTargetBody();
    addTargetParam = {
      workspace: values.workspace,
      name: values.name,
      language: values.language,
      template: values.template == undefined ? "" : values.template,
      owner: this.operateType == _EDIT ? this.owner: this.userId,
      ownerName: this.operateType == _EDIT ? this.ownerName: this.userName,
      isShared: values.isShared,
      tagName: {
        mainCategory: values.bugType,
        subCategory: values.subBugType,
        defectDetail: values.bugDetail,
      },
      extendedLabel: values.extendedLabel,
      customLable: values.customLable,
      instituteLabel: values.instituteLabel,
      answer: this.targetAnswer,
      targets: this.targetList,
    };
    let editTargetParam = JSON.parse(JSON.stringify(addTargetParam));
    editTargetParam['id'] = this.targetId;
    // 当编辑状态当前用户与靶子所有者不同时重写部分参数
    if (this.operateType == _EDIT && this.owner !== this.userId) {
      editTargetParam['workspace'] = this.srcObj.workspace;
      editTargetParam['name'] = this.srcObj.name;
      editTargetParam['language'] = this.srcObj.language;
      editTargetParam['template'] = this.srcObj.template;
      editTargetParam['isShared'] = this.srcObj.isShared;
      editTargetParam['extendedLabel'] = this.srcObj.extendedLabel;
      editTargetParam['customLable'] = this.srcObj.customLable;
      editTargetParam['tagName'] = {
        mainCategory: this.srcObj.bugType,
          subCategory: this.srcObj.subBugType,
          defectDetail: this.srcObj.bugDetail,
      };
    }
    let reqBody: any = {
      name: this.operateType == _EDIT ? _EDIT : _ADD,
      parameters: this.operateType == _EDIT ? editTargetParam : addTargetParam,
    };
    let sucMsg: string = this.operateType == _EDIT ? '修改成功' : '创建成功！';
    let errorMsg: string =
      this.operateType == _EDIT ? '修改失败' : '创建失败！';
    this.manageService.manageTargetApi(reqBody).subscribe(
      (res) => {
        this.plxMessageService.success(sucMsg, '');
        this.router.navigateByUrl('/main/manage/target/list');
      },
      (err) => {
        this.plxMessageService.error(errorMsg, err.cause);
      }
    );
  }

  public isAllowFile(
    items: FileItem[],
    allowedFileType: any,
    uploaderInstance: any
  ): boolean {
    if (items.length) {
      const fileName = items[0].file.name;
      const ext = fileName.slice(fileName.lastIndexOf('.') + 1).toLowerCase();
      const allowFiles = allowedFileType.filter((_type) => ext === _type);
      if (allowFiles && allowFiles.length) {
        return true;
      } else {
        // uploaderInstance.deleteFiles(items);
        this.plxMessageService.show('error', {
          title: '不支持上传该格式的文件，请删除后重新选择',
          isLightweight: true,
        });
      }
    }
    return false;
  }

  loadWorkSpaces() {
    this.manageService.getWorkspace().subscribe(res => {
      const workspaceInfos = [...res];
      workspaceInfos.forEach(ws => {
        this.wsMap.set(ws.id, ws.name);
      });
      this.setWsSelector();
    }, err => {
      this.plxMessageService.error('获取工作空间失败！', errCodeAdapter[err.error?.errCode]);
    });
  }

  setWsSelector() {
    const options = [];
    this.wsMap.forEach((value, key) => {
      options.push({ value: key, label: value });
    });
    this.formSetting.fieldMap['workspace'].valueSet = options;
  }
  onSelectWorkSpace(values) {
    this.manageService.getWorkSpaceLang(values.workspace).subscribe(
      (res) => {
        this.langOptions = []
        // 去重
        let val = []
        for (const key in res) {
          if (Object.prototype.hasOwnProperty.call(res, key)) {
            val.push(key)
          }
        }

        for (let index = 0; index < val.length; index++) {
          const element = val[index];
          let options = {
            label: element,
            value: element,
          }
          this.langOptions = this.langOptions.concat(options)
        }
        this.formSetting.fieldMap['language'].valueSet = this.langOptions;
        // 默认值置空，避免切换时没有修改之前的数据
        if (!val.includes(values.language)) {
          this.formSetting.formObject.controls['language'].setValue(null);
          this.formSetting.formObject.controls['language'].markAsDirty();
        }

      },
      (error) => {
        this.plxMessageService.error('获取工作空间所属语言失败', error.cause);
      }
    );

    this.manageService.getTemplateByWorkspace(values.workspace).subscribe(
      (res) => {
        if (res == null) {
          this.plxMessageService.error('工作空间没有启用打靶规范', null);
        } else {

          this.templateOptions = []
          for (let index = 0; index < res.length; index++) {
            const element = res[index];
            let options = {
              label: element,
              value: element,
            }
            this.templateOptions = this.templateOptions.concat(options)
          }
          this.formSetting.fieldMap['template'].valueSet = this.templateOptions;
          // 默认值置空，避免切换时没有修改之前的数据
          if (!res.includes(values.template)) {
            this.formSetting.formObject.controls['template'].setValue(null);
            this.formSetting.formObject.controls['template'].markAsDirty();
          }
        }
      },
      (error) => {
        this.plxMessageService.error('获取打靶规范失败失败', error.cause);
      }
    );
  }

  UpdateLanguageOPtions(workspace: string, template: string, language: string) {
    if (!template) {
      return
    }
    this.manageService.getLangByworkspaceAndTemplate(workspace, template).subscribe(
      (res) => {
        this.langOptions = []
        // 去重
        let val = []
        for (const key in res) {
          if (Object.prototype.hasOwnProperty.call(res, key)) {
            val.push(key)
          }
        }

        for (let index = 0; index < val.length; index++) {
          const element = val[index];
          let options = {
            label: element,
            value: element,
          }
          this.langOptions = this.langOptions.concat(options)
        }
        this.formSetting.fieldMap['language'].valueSet = this.langOptions;
        // 默认值置空，避免切换时没有修改之前的数据
        if (!val.includes(language)) {
          this.formSetting.formObject.controls['language'].setValue(null);
          this.formSetting.formObject.controls['language'].markAsDirty();
        }
      },
      (error) => {
        this.plxMessageService.error('规范下面没有语言，请修改规范！', error.cause);
      }
    );

  }


  /**
   * 获取某一语言下的模板信息(即规范), 以用来初始化标签信息中的各下拉项
   * @param lang 所属语言
   * @param workspace 工作空间
   * @param detail 靶子详情values
   * @returns
   */
  UpdateTagOptions(template: string, lang: string, workspace: string, detail?: any) {
    if (!lang) {
      return;
    }
    //this.manageService.getTemplateByWorkspaceAndLang(workspace, lang).subscribe(
    //  (res) => {
    //    this.templateOptions = []
    //    for (let index = 0; index < res.length; index++) {
    //      const element = res[index];
    //      let options = {
    //        label: element,
    //        value: element,
    //      }
    //      this.templateOptions = this.templateOptions.concat(options)
    //    }
    //    this.formSetting.fieldMap['template'].valueSet = this.templateOptions;
    //    // 默认值置空，避免切换时没有修改之前的数据
    //    if (!res.includes(template)) {
    //      this.formSetting.formObject.controls['template'].setValue(null);
    //      this.formSetting.formObject.controls['template'].markAsDirty();
    //    }
    //  },
    //  (error) => {
    //    this.plxMessageService.error('获取打靶规范失败失败', error.cause);
    //  }
    //);

    this.manageService.getWorkSpaceLangDefect(workspace, lang, template).subscribe(
      (res) => {
        if (res && res.result == 'success' && res.detail) {
          this.bugTypeOptions = [{ label: '所有', value: '所有' }];
          this.defectDetail = res.detail;
          for (const key of Object.keys(this.defectDetail)) {
            this.bugTypeOptions.push({ label: key, value: key });
          }
          this.formSetting.fieldMap['bugType'].valueSet = this.bugTypeOptions;
          // 如果是编辑状态设置缺陷大类、缺陷小类和缺陷细项的值
          if (detail) {
            this.srcObj['bugType'] = detail['tagName']['mainCategory'];
            this.srcObj['subBugType'] = detail['tagName']['subCategory'];
            this.srcObj['bugDetail'] = detail['tagName']['defectDetail'];
            this.refreshDefValues = !this.refreshDefValues;
          }
        }
      },
      (error) => {
        this.bugTypeOptions = [{ label: '所有', value: '所有' }];
        this.formSetting.fieldMap['bugType'].valueSet = this.bugTypeOptions;
        this.formSetting.formObject.controls['bugType'].setValue('所有');
        this.formSetting.formObject.controls['bugType'].markAsDirty();
        this.plxMessageService.error('获取靶子编码信息失败,请选择正确的打靶规范！', error.cause);
      }
    );


  }

}

import {Component, OnInit, ViewChild} from '@angular/core';
import {FieldType, FileItem, PlxActiveModal, PlxBreadcrumbItem, PlxFileUpLoaderComponent, PlxMessage} from "paletx";
import {Router} from "@angular/router";

@Component({
  selector: 'app-template-add',
  templateUrl: './template-add.component.html',
  styleUrls: ['./template-add.component.css']
})
export class TemplateAddComponent implements OnInit {

  @ViewChild('plxDragUploader') public plxDragUploader: any;

  user: string;
  breadModel: PlxBreadcrumbItem[] = [];
  formSetting: any;
  templateUploaderConfig: any;
  allowedTemplateType = ['xlsm'];
  workspaces: any[];
  callbackFunc: Function;

  UPLOAD_TEMPLATE_URL = '/api/code-shooting/v1/actions';

  @ViewChild('settingForm') settingForm;

  constructor(/*private router: Router,*/
              private plxMessageService: PlxMessage,
              private modal: PlxActiveModal) { }

  ngOnInit(): void {
    this.user = localStorage.getItem('name') + localStorage.getItem('user');
    //this.initBreadModel();
    this.initFileUploaderConfig();
    this.initFormSetting();
  }

  /*initBreadModel() {
    this.breadModel = [
      {label: '规范管理', routerLink: '../list', name: 'templateManage'},
      {label: '上传规范', name: 'uploadTemplate'},
    ];
  }*/

  initFormSetting() {
    this.formSetting = {
      isShowHeader: false,
      isGroup: false,
      hideGroup: true,
      srcObj: {},
      advandedFlag: false,
      lableClass: 'col-sm-3',
      componentClass: 'col-sm-8',
      buttons: [
      ],
      fieldSet: [
        {
          group: '',
          fields: [
            {
              name: 'workspace',
              label: '工作空间',
              type: FieldType.SELECTOR,
              required: true,
              multiple: false,
              valueSet: this.buildWsSelector(),
              callback: (values, $event, controls) => {
                this.templateUploaderConfig.defaultUploadUrl = this.UPLOAD_TEMPLATE_URL + "/" + values.workspace+ "/uploadTemplate"
              },
              binding: {
                dropdownContainer: 'body',
                scrollSelectors: ['.modal-body'],
              }
            },
            {
              name: 'template',
              label: '上传规范',
              notice: '仅支持《代码打靶落地模板-vX.X.xlsm》格式的文件',
              type: FieldType.CUST_COMPONENT,
              required: true,
              component: PlxFileUpLoaderComponent,
              binding: {
                pxFileUploaderConfig: this.templateUploaderConfig,
              }
            }
          ]
        }
      ]
    };
  }

  initFileUploaderConfig() {
    this.templateUploaderConfig = {
      defaultUploadUrl: this.UPLOAD_TEMPLATE_URL + "/uploadTemplate",
      showCustomUrl: false,
      useBrowserProgress: true,
      autoUpload: true,
      showDelIcon: true,
      queueLimit: 1,
      allowedFileType: this.allowedTemplateType,
      additionalParameter: {
        operator: this.user,
      },
      fileListConfig: {
        showFileName: true,
        showDetail: false,
        showTooltip: true
      },
      // 文件验证不通过时的回调
      onFilterFile: (filter: { name: string }) => {
        let msg = '';
        if (filter.name === 'queueLimit') {
          msg = '只能上传一个规范文件，请先删除再重新上传';
        }
        return msg;
      },
      onSuccessItem: (item) => {
        this.plxMessageService.success(item.file.name + '上传成功', '');
        //this.router.navigateByUrl('/main/manage/template/list');
        this.modal.close();
        this.callbackFunc();
      },
      onErrorItem: (item: any, response: any, status: any) => {
        const resp = JSON.parse(response)
        this.plxMessageService.error(item.file.name + '上传失败', resp.status);
      },
      // 上传成功的文件是否能删除
      onDelItem: (item) => {
        return Promise.resolve(true);
      },
      // 上传失败或未上传的文件是否能删除
      onRemoveItem: (item) => {
        return Promise.resolve(true);
      },
      // 文件选择完毕后的回调
      onAfterAddingAll: (items: FileItem[]) => {
        return new Promise(resolve => resolve(this.isAllow(items, this.allowedTemplateType, this.plxDragUploader)));
      }
    };
  }

  public isAllow(items: FileItem[],  allowedFileType: any, uploaderInstance: any): boolean {
    let values = this.settingForm.getFormObject().value;
    if(!values.workspace) {
      this.plxMessageService.show('error', {
        title: '请先选择工作空间，再删除文件后重新上传',
        isLightweight: true
      });
      return false;
    }
    if (items.length) {
      const fileName = items[0].file.name;
      if(this.validateTemplateFileName(fileName)) {
        return true;
      } else {
        this.plxMessageService.show('error', {
          title: '不支持上传该格式的文件，请删除后重新选择',
          isLightweight: true
        });
      }
    }
    return false;
  }

  validateTemplateFileName(fileName): boolean {
    const fileRegexp = /^代码打靶落地模板-v((\d)+)\.((\d)+)\.xlsm$/;
    if (!fileRegexp.test(fileName)) {
      return false;
    }
    return true;
  }

  buildWsSelector() {
    const options = [];
    if(this.workspaces && this.workspaces.length > 0) {
      for(const ws of this.workspaces) {
        options.push({value: ws.id, label: ws.name});
      }
    }
    return options;
  }

  onClose() {
    this.modal.close();
  }
}

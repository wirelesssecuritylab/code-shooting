import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {FileItem, PlxMessage} from 'paletx';
import {errCodeAdapter} from "../../shared/class/row";
import {QueryService} from "../query/score/score.service";
@Component({
  selector: 'app-upload',
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.css']
})
export class UploadComponent implements OnInit {
  @Input() public disabled: boolean;
  @Input() public targetId: string;
  @Input() public language: string;
  @Input() public userId: string;
  pxFileUploaderConfig:any = {};
  res: any;
  message:string;
  hoverStr: string;
  public allowedFileType = ['xlsm'];
  constructor(private plxMessageService: PlxMessage, private queryService: QueryService) {}
  @ViewChild('plxSingleDragUploader') public plxSingleDragUploader: any;

  ngOnInit(): void {
    this.hoverStr = '将答卷所在的excel文件作为一个整体提交到打靶服务';
    this.initFileSelector();
  }

  private initFileSelector() {
    let lang = this.queryService.transferCharacter(this.language);
    return this.pxFileUploaderConfig = {
      defaultUploadUrl: `/api/code-shooting/v1/answers/submit?rangeId=${this.targetId}&language=${lang}&userId=${this.userId}`,
      asyncProgressBar: true,
      showCustomUrl: false,
      //allowedFileType: this.allowedFileType,
      useBrowserProgress: true,
      autoUpload: false,
      showDelIcon: true,
      fileListConfig: {
        showFileName: true,
        showDetail: true,
        showTooltip: true
      },
      maxFileSize: 1024 * 1024 * 100, // 100M
      queueLimit: 5,
      minFileSize: 10,
      maxFileNameLen: 40,
      onAfterAddingAll: (items: FileItem[]) => {
        return new Promise(resolve => resolve(this.isAllowFile(items, this.plxSingleDragUploader)));
      },
      onSuccessItem: (item) => {
        console.log(item.file.name + ' upload success');
        console.log(item);
        console.log(this.pxFileUploaderConfig.defaultUploadUrl);
        // this.plxMessageService.show('success',
        //   {
        //     title: '上传成功',
        //     isLightweight: true
        //   });
        this.plxMessageService.success('上传成功','');
      },
      onErrorItem: (item) => {
        this.res = JSON.parse(item._xhr.response);
        console.log(item.file.name + ' upload failure');
        console.log(item);
        console.log(this.pxFileUploaderConfig.defaultUploadUrl);
        if (item._xhr.status == 429) {
          // this.plxMessageService.show('error',
          //   {
          //     title: '当前请求人数过多，请稍后重试。',
          //     isLightweight: true
          //   });
          this.plxMessageService.error('当前请求人数过多，请稍后重试。', '');
        } else {
          console.log(this.res);
          this.message = errCodeAdapter[this.res?.errCode];
          if (this.message != '') {
            // this.plxMessageService.show('error',
            //   {
            //     title: '上传失败：'+ this.message + '。',
            //     isLightweight: true
            //   });
            this.plxMessageService.error('上传失败：'+ this.message + '。', '');
          } else {
            // this.plxMessageService.show('error',
            //   {
            //     title: '上传失败，请联系管理员。',
            //     isLightweight: true
            //   });
            this.plxMessageService.error('上传失败，请联系管理员。', '');
          }
        }
      },
      onDelItem: (item) => {
        return Promise.resolve(false);
      },
      onRemoveItem: (item) => {
        return Promise.resolve(true);
      },
      onFilterFile: (filter: { name: string }, filters: Array<{ name: string, fileName?: string }>) => {
        let msg = '';
        if (filter.name === 'queueLimit') {
          msg = '文件数量超过限制';
        }
        if (filter.name === 'fileSize') {
          msg = filter['fileName'] + '大小超过限制';
        }
        if (filter.name === 'minFileSize') {
          msg = filter['fileName'] + '大小小于最小值';
        }
        if (filter.name === 'maxFileNameLen') {
          msg = filter['fileName'] + '文件名超过限制';
        }
        if (filters) {
          console.log('invalid info', filters);
        }
        return msg;
      },
      onCancelItem: (item) => {
        console.log('cancel');
      },
      // externalParameterInvoke: this.getParams.bind(this), // 重要
    };
  }

  public isAllowFile(items: FileItem[], uploaderInstance: any): boolean {
    if (items.length) {
      const fileName = items[0].file.name;
      const ext = fileName.slice(fileName.lastIndexOf('.') + 1).toLowerCase();
      const allowFiles = this.allowedFileType.filter(_type => ext === _type);
      if (allowFiles && allowFiles.length) {
        return true;
      } else {
        uploaderInstance.deleteFiles(items);
        // this.plxMessageService.show('error', {
        //   title: '只支持xlsm格式的excel文件',
        //   isLightweight: true
        // });
        this.plxMessageService.warning('只支持xlsm格式的excel文件', '');
      }
    }
    return false;
  }

}

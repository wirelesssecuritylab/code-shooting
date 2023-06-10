import {Component, OnInit} from '@angular/core';
import {BaseOperate} from "../../misc/BaseOperate";
import {DefectService} from "../../services/defect-service";
import {JigsawNotification, PopupService} from "@rdkmaster/jigsaw";
import {DefectSelectDialog} from "./defect-select/defect-select.dialog";
import {TargetInfo} from "../../misc/target-types";
import {Utils} from "../../misc/utils";
import { PlxModal } from 'paletx';

@Component({
    selector: 'tp-edit',
    template: `
        <!--打靶时显示div-->
        <div class="comment-tools" *ngIf="initData?.shootType !== 'view'" (mouseover)="showCurComment()"
              placement="right-top" [triggers]="'manual'" [plxPopover]="popTable" #p="plxPopover" (mouseover)="p.open()"
              popoverMaxHeight="200px"popoverMaxWidth="550px"
              style="top:25px;left:15px;margin-left:500px;">
            <button class="plx-icon-btn plx-icon-word-btn plx-ico-modify-16 comment-btn-modify" (click)="openCommentsDialog(commentContent, 'edit')"></button>
            <div class="comment-line"></div>
            <button class="plx-icon-btn plx-icon-word-btn plx-ico-remove-16 comment-btn-delete" (click)="openCommentsDialog(commentContent, 'delete')"></button>
        </div>

         <!--查看答卷时显示div-->
        <div class="comment-tools" *ngIf="initData?.shootType == 'view'" title="查看评审意见" style="top:25px;left:15px;margin-left:500px;">
            <button class="plx-icon-btn plx-icon-word-btn plx-ico-data-view-16 comment-btn" (mouseover)="showCurComment()"
              placement="right-top" [triggers]="'manual'" [plxPopover]="popTable" #p="plxPopover" (mouseover)="p.open()"
              popoverMaxHeight="200px"popoverMaxWidth="550px"></button>
        </div>

        <!--鼠标移到编辑或查看按钮区时显示框 -->
        <ng-template #popTable>
          <div class="plx-title-level4" style="margin-bottom:15px;">评审列表</div>
          <plx-table [data]="$defectTarget"
            [tableType]="'sm'"
            [columns]="columns"
            [resizable]="false"
            [showColFilterToggle]="false"
            [showCustomCol]="false"
            [showGlobalFilter]="false"
            [showPagination]="false">
          </plx-table>
        </ng-template>

        <!--点击编辑或删除按钮时弹出来的列表模态框-->
        <ng-template #commentContent let-c="close" let-d="dismiss">
          <div class="modal-header">
            <h4 class="modal-title">请选择一条评审</h4>
            <button type="button" class="close" (click)="d('Cross click')">
                <span class="plx-ico-close-16"></span>
            </button>
          </div>
          <div class="modal-body">
            <plx-table [data]="$defectTarget"
              [tableType]="'sm'"
              [columns]="columns"
              [resizable]="false"
              [showColFilterToggle]="false"
              [showCustomCol]="false"
              [showGlobalFilter]="false"
              [showPagination]="false"
              [showCheckBox]="true"
              (checkboxInfoChange)="checkboxInfoChange($event)"
              [isRelatedCheckboxOnRowClick]="true">
            </plx-table>
          </div>
          <div class="modal-footer">
              <div class="form-group">
                  <div class="btnGroup modal-btn"  style="margin-top: -16px;">
                      <button *ngIf="curOperateType == 'edit'" type="button" class="plx-btn plx-btn-primary" [disabled]="!hasSelectData" (click)="confirmModify()">确定</button>
                      <button *ngIf="curOperateType == 'delete'" type="button" class="plx-btn plx-btn-error" [disabled]="!hasSelectData" (click)="confirmDelete();c('Close click')">删除</button>
                      <button type="button" class="plx-btn" (click)="cancel()">取消</button>
                  </div>
              </div>
          </div>
      </ng-template>
    `,
    styles: [`
        .comment-btn-delete{
          padding: 0px;
          height:18px;
          line-height:18px;
          color: red;
        }

        .comment-btn-modify {
          padding: 0px;
          height:18px;
          line-height:18px;
          /* color:white; */
        }
        .comment-tools {
            /* background: #64A3E6; */
            border-radius: 3px;
            /* height: 32px; */
            height: 18px;
            display: flex;
            align-items: center;
            padding: 0 8px;
        }

        :host ::ng-deep .comment-tools .plx-icon-btn.plx-icon-word-btn {
          /* color: white; */
          /* color: blue; */
        }

        .comment-icon {
            margin: 0 5px;
            color: white;
            font-size: 16px;
            cursor: pointer;
        }

        .comment-icon:hover {
            color: #9FC5F0 !important;
        }

        .comment-line {
            width: 1px;
            height: 18px;
            background-color: white;
            margin: 0 5px;
        }
    `]
})
export class EditComponent extends BaseOperate implements OnInit{
  public $defectTarget: any[] = [];
  public columns: any[] = [];
  public modal: any;
  public hasSelectData: boolean = false;
  public curOperateType: string;
  public curInitData: any;

  constructor(private _defectService: DefectService, private _popupService: PopupService,
              private modalService: PlxModal) {
      super();
  }

  ngOnInit(): void {
    this.setColumns();
  }

  /**
   * 设置评审列表表格列
   */
  public setColumns() {
    this.columns = [
      {
        key: 'StartLineNum',
        title: '起始行号',
        show: true,
        width: '80px',
        fixed: true,
      },
      {
        key: 'EndLineNum',
        title: '结束行号',
        show: true,
        width: '80px',
        fixed: true,
      },
      {
        key: 'DefectClass',
        title: '缺陷大类',
        show: true,
        width: '80px',
        fixed: true,
      },
      {
        key: 'DefectSubClass',
        title: '缺陷小类',
        show: true,
        width: '80px',
        fixed: true,
      },
      {
        key: 'DefectDescribe',
        title: '缺陷细项',
        show: true,
        width: '80px',
        fixed: true,
      },
      {
        key: 'Remark',
        title: '缺陷备注',
        show: true,
        width: '80px',
        fixed: true,
      }
    ];
  }

  /**
   * 鼠标移入时显示当前开始行的所有评审列表信息
   */
  showCurComment() {
    console.log("编辑弹出");
    this.curInitData = JSON.parse(JSON.stringify(this.initData));
    this.$defectTarget = this._defectService.getSameLineTarget(this.initData?.range, this.initData?.fileName);
    this.$defectTarget.map(item => {
      item['isChecked'] = false;
      return item;
    });
  }

  /**
   * 编辑单个评审意见
   * @returns
   */
  public _$edit(): void {
    let that = this;
    const editTarget = this._defectService.getTarget(that.curInitData?.range, that.curInitData?.fileName);
    if (!editTarget) {
        JigsawNotification.showWarn('未检索到打靶记录！');
        return;
    }
    this.curInitData = {...this.curInitData, defectTarget: editTarget};
    const selectDialog = this._popupService.popup(DefectSelectDialog, Utils.getModalOptions(), this.curInitData);
    const selectDialogHandler = selectDialog.answer.subscribe((result: TargetInfo) => {
      selectDialogHandler.unsubscribe();
      this._defectService.editTarget(result);
      this.answer.emit(result);
    });
  }

  /**
   * 删除单个评审意见
   * @returns
   */
  public _$delete(): void {
    let that = this;
    const target = this._defectService.deleteTarget(that.curInitData?.range, that.curInitData?.fileName);
    if (!target) {
        JigsawNotification.showWarn('未检索到打靶记录！');
        return;
    }
    this.answer.emit(target);
  }

  /**
   * 取消按钮
   */
  public cancel(): void {
    this.modal.close();
  }

  /**
   * 确认修改按钮
   */
  confirmModify() {
    this._$edit();
    this.modal.close();
  }

  /**
   * 确认删除按钮
   */
  confirmDelete() {
    this._$delete();
    this.modal.close();
  }
  /**
   * 编辑或删除时打开评审列表信息模态框
   * @param content
   * @param operateType 修改：edit, 删除： delete
   */
  public openCommentsDialog(content: any, operateType: string) {
    this.$defectTarget = [];
    this.showCurComment();
    this.hasSelectData = false;
    this.curOperateType = operateType;
    const size: 'sm' | 'lg' |'xs' = 'sm';
    const options = {
      size: size,
      // contentClass: 'plx-modal-custom-content',
      escCallback: this.escCallback.bind(this),
      destroyOnClose: true,
      modalId: 'plx-modal-1',
      openCallback: () => {
        console.log('open');
      }
    };
    this.modal = this.modalService.open(content, options);
  }

  /**
   * 模态框Esc键调用函数
   * @returns
   */
  public escCallback(): boolean {
    console.info('escCallback');
    return true;
  }

  /**
   * 表格行选中时事件函数
   * @param $event
   */
  checkboxInfoChange($event): void {
    console.log($event);
    if ($event.action == 'check') {
      this.hasSelectData = true;
    } else {
      this.hasSelectData = false;
    }
    let range = {
      endColumn: $event.actionData[0].EndColNum,
      endLineNumber: $event.actionData[0].EndLineNum,
      startColumn: $event.actionData[0].StartColNum,
      startLineNumber: $event.actionData[0].StartLineNum
    };
    this.curInitData.range = range;
  }

  ngOnDestroy(): void {
    // 销毁模态框实例
    console.log("ngOnDestroy");
    this.modalService.destroyModalInstance();
  }
}

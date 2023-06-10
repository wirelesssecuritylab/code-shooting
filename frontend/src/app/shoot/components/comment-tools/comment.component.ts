import {Component, OnInit} from '@angular/core';
import {BaseOperate} from "../../misc/BaseOperate";
import {DefectService} from "../../services/defect-service";
import {PopupService, JigsawNotification} from "@rdkmaster/jigsaw";
import {DefectSelectDialog} from "./defect-select/defect-select.dialog";
import {TargetInfo} from "../../misc/target-types";
import {Utils} from "../../misc/utils";

// 编辑器中显示的评审意见组件:评审意见+操作按钮
@Component({
    selector: 'tp-comment',
    template: `
        <div style="width:1100px;height:18px;">
          <div class="comment"  [jigsawTooltip]="_$tipContent" [jigsawTooltipTheme]="'light'"  (click)="clickedComment($event)"
                                (mouseover)="showBgColor()" (mouseout)="clearBgColor()" (dblclick)="doubleClickedComment($event)">
              评审意见: {{_$defectTarget?.DefectDescribe}}
              <i class="plx-btn-icon plx-ico-pick-16 comment-result-right" [hidden]="initData.shootType != 'view' || _$targetScore <= 0"></i>
              <i class="plx-btn-icon plx-ico-close-16 comment-result-wrong" [hidden]="initData.shootType != 'view' || _$targetScore != 0"></i>
          </div>
          <div class="comment-btn-hide">
            <button class="plx-icon-btn plx-icon-word-btn plx-ico-modify-16 comment-btn-modify" (click)="editComment($event)"></button>
            <div class="comment-line"></div>
            <button class="plx-icon-btn plx-icon-word-btn plx-ico-remove-16 comment-btn-delete" (click)="deleteComment($event)"></button>
          </div>
        </div>
    `,
    styles: [`
        /**评审区Hover时改变字体和背景色 */
        div.comment:hover{
            background: #108EE9;
            opacity: 0.3;
            color: #fffffe;
        }
        /**点击时样式 */
        .focusComment{
            width: max-content;
            margin-left: 1px;
            float: left;
            width:1000px;
            font-size: 12px;
            background: #108EE9;
            opacity: 0.3;
            color: #fffffe;
        }
        /**正常时的样式 */
        .comment {
            width: max-content;
            margin-left: 1px;
            float: left;
            width:1000px;
            font-size: 12px;
            color: #737373;
        }

        .comment-result-right{
          padding-left: 10px;
          color: #54C91B;
          font-size: 16px;
          font-weight: bold;
          vertical-align: text-top;
          height:18px;
          line-height:18px;
        }

        .comment-result-wrong{
          padding-left: 10px;
          color: red;
          font-size: 16px;
          font-weight: bold;
          vertical-align: text-top;
          height:18px;
          line-height:18px;
        }

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
        /** 按钮区隐藏时的样式*/
        .comment-btn-hide {
            /* background: #64A3E6; */
            border-radius: 3px;
            /* height: 32px; */
            height: 18px;
            display: flex;
            align-items: center;
            padding: 0 8px;
            float: left;
            visibility:hidden;
        }
        /** 按钮区显示时的样式*/
        .comment-btn-show {
            /* background: #64A3E6; */
            border-radius: 3px;
            /* height: 32px; */
            height: 18px;
            display: flex;
            align-items: center;
            padding: 0 8px;
            float: left;
            visibility:visible;
        }

        .comment-line {
            width: 1px;
            height: 18px;
            /* background-color: #1993ff; */
            color: #1993ff;
            margin: 0 5px;
        }

    `]
})
export class CommentComponent extends BaseOperate implements OnInit {
  public _$defectTarget: TargetInfo;
  public _$tipContent: string;
  public curInitData: any;
  public _$targetScore: number;
  public timer = null;
  constructor(private _defectService: DefectService, private _popupService: PopupService) {
      super();
  }

  ngOnInit(): void {
    this._$defectTarget = this._defectService.getTarget(this.initData?.range, this.initData?.fileName);
    if (!this._$defectTarget) {
        return;
    }

    if (this.initData.shootType == 'view') {
      this._$targetScore = this._defectService.getTargetScore(this.initData?.range, this.initData?.fileName);
      if (this._$targetScore == -1) {
        return;
      }
    }
    this._$tipContent = Utils.stripBlank(`
      起始行: ${this._$defectTarget.StartLineNum}
      结束行: ${this._$defectTarget.EndLineNum}
      缺陷大类：${this._$defectTarget.DefectClass}
      缺陷小类：${this._$defectTarget.DefectSubClass}
      缺陷描述：${this._$defectTarget.DefectDescribe}
      缺陷备注：${this._$defectTarget.Remark || ''}
    `);

    this.curInitData = this.initData;
  }

  /**
   * 单击评审意见
   * @param evt
   * @returns
   */
  public clickedComment(evt: any) {
    // 查看答卷模式下不显示编辑和删除按钮
    if (this.initData.shootType == 'view') {
      return;
    }
    clearTimeout(this.timer);
    this.timer = setTimeout(() => {
      this.showActionBtn(evt);
    },200);
  }

  /**
   * 双击评审意见
   * @param e
   */
  public doubleClickedComment(e: any) {
    // 查看答卷模式下不显示编辑和删除按钮
    if (this.initData.shootType == 'view') {
      return;
    }
    clearTimeout(this.timer);
    this.editComment(e);
  }

  /**
   * 显示对应代码区的背景色
   */
  public showBgColor() {
    let that = this;
    let target = JSON.parse(JSON.stringify(this._defectService.getTarget(that.curInitData?.range, that.curInitData?.fileName)))
    target['mouseAction'] = "over";
    target['range'] = that.curInitData?.range;
    this.answer.emit(target);
  }

  /**
   * 清除对应代码区的背景色
   */
  public clearBgColor() {
    let that = this;
    let target = JSON.parse(JSON.stringify(this._defectService.getTarget(that.curInitData?.range, that.curInitData?.fileName)))
    target['mouseAction'] = "out";
    target['range'] = that.curInitData?.range;
    this.answer.emit(target);
  }
  /**
   * 点击评审意见显示操作按钮
   * @param evt
   * @returns
   */
  public showActionBtn(evt: any) {
    let commentClassName = evt.target.getAttribute("class");
    if (commentClassName == 'comment') {
      evt.target.classList.replace("comment", "focusComment");
    } else {
      evt.target.classList.replace("focusComment", "comment");
    }

    let btnClassName = evt.target.parentElement.lastElementChild.getAttribute("class");
    if (btnClassName == 'comment-btn-hide') {
      evt.target.parentElement.lastElementChild.classList.replace("comment-btn-hide", "comment-btn-show");
    } else {
      evt.target.parentElement.lastElementChild.classList.replace("comment-btn-show", "comment-btn-hide");
    }
  }

  /**
   * 点击编辑按钮修改评审信息
   * @param evt
   * @returns
   */
  public editComment(evt: any) {
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
      result['isEditOrDelete'] = true;
      this.answer.emit(result);
    });
  }

  /**
   * 点击删除按钮删除评审信息
   * @param evt
   */
  public deleteComment(evt: any) {
    let that = this;
    const target = this._defectService.deleteTarget(that.curInitData?.range, that.curInitData?.fileName);
    if (!target) {
        JigsawNotification.showWarn('未检索到打靶记录！');
        return;
    }
    target['isEditOrDelete'] = true;
    this.answer.emit(target);
  }
}

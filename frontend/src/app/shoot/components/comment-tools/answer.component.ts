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
          <div class="comment"  [jigsawTooltip]="_$tipContent" [jigsawTooltipTheme]="'light'" 
                                (mouseover)="showBgColor()" (mouseout)="clearBgColor()">
                正确答案: {{_targetAnswer?.DefectDescribe}}
                <i class="plx-btn-icon plx-ico-config-script-16 comment-answer"></i>
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
        .comment-answer{
          padding-left: 10px;
          color: #F7DC6F;
          font-size: 16px;
          font-weight: bold;
          vertical-align: text-top;
          height:18px;
          line-height:18px;
        }

    `]
})
export class AnswerComponent extends BaseOperate implements OnInit {
  public _targetAnswer: TargetInfo;
  public _$tipContent: string;
  public curInitData: any;
  public _$targetScore: number;
  public timer = null;
  constructor(private _defectService: DefectService, private _popupService: PopupService) {
      super();
  }

  ngOnInit(): void {
    this._targetAnswer = this._defectService.getAnswer(this.initData?.range, this.initData?.fileName);
    if (!this._targetAnswer) {
        return;
    }

    this._$tipContent = Utils.stripBlank(`
      起始行: ${this._targetAnswer.StartLineNum}
      结束行: ${this._targetAnswer.EndLineNum}
      缺陷大类：${this._targetAnswer.DefectClass}
      缺陷小类：${this._targetAnswer.DefectSubClass}
      缺陷描述：${this._targetAnswer.DefectDescribe}
      缺陷备注：${this._targetAnswer.Remark || ''}
    `);

    this.curInitData = this.initData;
  }

 

  /**
   * 显示对应代码区的背景色
   */
  public showBgColor() {
    let that = this;
    let target = JSON.parse(JSON.stringify(this._defectService.getAnswer(that.curInitData?.range, that.curInitData?.fileName)))
    target['mouseAction'] = "over";
    target['range'] = that.curInitData?.range;
    this.answer.emit(target);
  }

  /**
   * 清除对应代码区的背景色
   */
  public clearBgColor() {
    let that = this;
    let target = JSON.parse(JSON.stringify(this._defectService.getAnswer(that.curInitData?.range, that.curInitData?.fileName)))
    target['mouseAction'] = "out";
    target['range'] = that.curInitData?.range;
    this.answer.emit(target);
  }
}

import {Component} from '@angular/core';
import {DefectService} from "../../services/defect-service";
import {BaseOperate} from "../../misc/BaseOperate";
import {PopupService} from "@rdkmaster/jigsaw";
import {DefectSelectDialog} from "./defect-select/defect-select.dialog";
import {TargetInfo} from "../../misc/target-types";
import {Utils} from "../../misc/utils";

@Component({
    selector: 'tp-add',
    template: `
        <div class="comment-add-tools" title="点击添加评审意见">
          <button class="plx-icon-btn plx-icon-word-btn plx-ico-add-16 comment-add-btn" (click)="_$addComment()"></button>
        </div>
    `,
    styles: [`
        .comment-add-btn{
          padding: 0px;
          height:16px;
          line-height:16px;
          color:white;

        }
        .comment-add-tools {
          background: #64A3E6;
          border-radius: 3px;
          height: 24px;
          display: flex;
          align-items: center;
          padding: 0 8px;
        }
    `]
})
export class AddComponent extends BaseOperate {
    constructor(private _defectService: DefectService, private _popupService: PopupService) {
        super();
    }

    public _$addComment(): void {
        const selectDialog = this._popupService.popup(DefectSelectDialog, Utils.getModalOptions(), this.initData);
        const selectDialogHandler = selectDialog.answer.subscribe((target: TargetInfo) => {
            selectDialogHandler.unsubscribe();
            this._defectService.addTarget(target, this.initData?.range);
            this.answer.emit(target);
        });
    }
}

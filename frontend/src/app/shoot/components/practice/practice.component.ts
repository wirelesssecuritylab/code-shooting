import { Component, Input, NgZone, ViewChild } from '@angular/core';
import { PopupService } from "@rdkmaster/jigsaw";
import { MonacoCodeEditor } from "../code-editor/code-editor";
import { EditComponent } from "../comment-tools/edit.component";
import { AddComponent } from "../comment-tools/add.component";
import { DOMService } from "../../services/dom-service";
import { TargetInfo } from "../../misc/target-types";
import { DecorateInfo, FileTreeNode, Range } from "../../misc/types";
import { DefectService } from "../../services/defect-service";
import { MarkedInfoService } from "../../services/marked-info-service";
import { Utils } from "../../misc/utils";
import { ShowDefectsDialog } from "../show-defects/show-defects.dialog";
import { CommentComponent } from "../comment-tools/comment.component";
import { ManageService } from '../../../admin/manage.service';
import { PlxMessage } from 'paletx';
import { AnswerComponent } from '../comment-tools/answer.component';

declare const monaco: any;

type TabTitle = {
  fileName: string,
  html?: string,
  htmlContext?: any
}

@Component({
  selector: 'tp-practice',
  templateUrl: './practice.component.html',
  styleUrls: ['./practice.component.scss']
})
export class PracticeComponent {
  // 输入属性：靶子id,使用时用split(',')拆分成列表形式，因为进入靶场选择某个语言后，有可能有多个靶子，打靶时可以点击下一个靶子进行继续打靶；
  @Input() _targetId: string;
  @Input() _rangeId: string;
  @Input() _shootType: string;
  @Input() _targetLang: string;
  @Input() _fromPage: string;
  // 已经打开 的文件树节点列表
  private _openedFiles: FileTreeNode[] = [];
  public _$currentFile: FileTreeNode;
  public _$code: string;
  public _$tabTitles: TabTitle[] = [];
  public _$fileSelectedIndex: number = 0;
  private _currentFileName: string;
  private hasRunRender: boolean = false;
  public targetId: string;
  public rangeId: string;

  @ViewChild('codeEditor')
  private _codeEditor: MonacoCodeEditor;
  private _cursorSelection: Range;
  public curShootResult: any;

  constructor(private _domService: DOMService,
    private _zone: NgZone,
    private _markedInfoService: MarkedInfoService,
    private _defectService: DefectService,
    private _popupService: PopupService,
    private manageService: ManageService,
    private plxMessageService: PlxMessage) {
  }

  /**
   * 编辑器准备就绪事件函数
   */
  public _$coderReady(): void {
    // 光标在编辑区点击或者选中
    this.curShootResult = this._defectService.shootingResult;
    if (!this.hasRunRender && this._$currentFile?.fileName) {
      this._renderDefects();
    }

    this._codeEditor.editor.onDidChangeCursorSelection((e) => {
      if (this._shootType === 'view') {
        return;
      }
      this._cursorSelection = e.selection;
    });

    // 鼠标按键弹起
    this._codeEditor.editor.onMouseUp(e => {
      // 查看答卷场景下不允许添加评审意见
      if (this._shootType === 'view') {
        return;
      }
      const endPosition = e.target.position;
      if (!endPosition || !this._cursorSelection) {
        return;
      }
      // 空行时也能添加评审
      if (this._cursorSelection.startLineNumber == this._cursorSelection.endLineNumber &&
        this._cursorSelection.startColumn == this._cursorSelection.endColumn &&
        this._cursorSelection.startColumn !== 1) {
        // 取消选择
        // this._clearSelectedArea();
        // this._renderDefects();
        return;
      }
      //编辑器从行号选择时光标会自动移至下一行，若选择多行，需要将光标所在行减去
      let endLineNumber = this._cursorSelection.endLineNumber;
      if (this._cursorSelection.startLineNumber != this._cursorSelection.endLineNumber
        && this._cursorSelection.endColumn === 1) {
        endLineNumber--;
      }
      const range: Range = new monaco.Range(
        this._cursorSelection.startLineNumber,
        this._cursorSelection.startColumn,
        endLineNumber,
        this._cursorSelection.endColumn
      );
      const target = this._defectService.getTarget(range, this._$currentFile.fileName);
      if (target) {
        // 如果存在交叉的target
        this._codeEditor.editor.setSelection(new monaco.Selection(0, 0, 0, 0));
        return;
      }
      this._removeAddCommentButton();
      this._decorateSelectedArea(range, true);
    });

    // 编辑器内任何地方鼠标按下后把已经显示的编辑和删除按钮隐藏掉，把点击选中的评审信息颜色重置掉；
    // 如果鼠标点击区正好是当前点击选中的评审信息所在行（点击的评审信息或者编辑和删除按钮）则不进行处理。
    this._codeEditor.editor.onMouseDown(e => {
      let curTargetElement = e.target.element;
      this._removeAddCommentButton(curTargetElement);
      // 把已经显示出编辑和删除按钮的地方 隐藏掉
      let btnElements = document.getElementsByClassName("comment-btn-show");
      for (let i = 0; i < btnElements.length; i++) {
        if (curTargetElement.parentNode.getAttribute("class") == "comment-btn-show" || e.target.element.parentElement.lastElementChild.getAttribute("class") == "comment-btn-show") {
          continue;
        }
        btnElements[i].classList.replace("comment-btn-show", "comment-btn-hide");
      }
      let commentElements = document.getElementsByClassName("focusComment");
      for (let j = 0; j < commentElements.length; j++) {
        if (curTargetElement.getAttribute("class") == "focusComment" || curTargetElement.parentElement.parentElement.firstElementChild.getAttribute("class") == "focusComment") {
          continue;
        }
        // 把有点击选中的评审意见区重置为浅底色
        commentElements[j].classList.replace("focusComment", "comment");
      }
    });
  }

  /**
   * 渲染缺陷信息即添加过的评审意见
   * @param fileName
   */
  public _renderDefects(fileName = this._$currentFile.fileName): void {
    // 加这个判断主要是因为这个函数有时会执行两次，导致一个评论上方会出来一个空的地区
    if (this.hasRunRender && this._currentFileName == fileName) {
      return;
    }

    this._defectService.targets
      .filter(target => target.FileName == fileName)
      .forEach(target => {
        const range: Range = new monaco.Range(target.StartLineNum, target.StartColNum, target.EndLineNum, target.EndColNum);
        this._decorateSelectedArea(range, false);
        this.hasRunRender = true;
      })

    if (this._shootType == 'view') {
      this._defectService.targetAnswers.filter(answer => answer.FileName == fileName).forEach(answer => {
        const range: Range = new monaco.Range(answer.StartLineNum, answer.StartColNum, answer.EndLineNum, answer.EndColNum);
        this._decorateSelectedArea2(range, false);
      });
      this.hasRunRender = true;
    }

  }

  /**
   * 提交答卷
   */
  public _$submit(): void {
    if (this._defectService.shootingResult.targets.length <= 0) {
      this.plxMessageService.info('你还未进行打靶，请先打靶再提交答卷', '');
      return;
    }
    this._popupService.popup(ShowDefectsDialog, Utils.getModalOptions());
  }

  /**
   * 装饰选中的区域
   * @param range
   * @param add 取值：true:添加操作 false:非添加操作
   */
  private _decorateSelectedArea(range: Range, add: boolean): void {
    // 根据操作类型设置选中的代码的底色
    if (add) {
      const decorates = this._codeEditor.editor.deltaDecorations([], [{
        range: range,
        options: {
          isWholeLine: true, // 是否整行选中
          className: 'selected-code'  // 选中行的类名，以关联选中底色
        }
      }]);
      this._markedInfoService.addDecoration(range, decorates[0], this._$currentFile.fileName);
    } else {
      const decorates = this._codeEditor.editor.deltaDecorations([], [{
        range: range,
        options: {
          isWholeLine: true,
          className: 'non-selected-code'
        }
      }]);
      this._markedInfoService.addDecoration(range, decorates[0], this._$currentFile.fileName);
    }
    this._zone.run(() => {
      if (add) {
        this._showAddCommentButton(range);
      } else {
        // this._showEditCommentTools(range); // 不再显示原来的编辑评审工具条组件
        this._showCommentTip(range);
        this._showCommentZone(range);
      }
    });
  }

  private _decorateSelectedArea2(range: Range, add: boolean): void {
    const decorates = this._codeEditor.editor.deltaDecorations([], [{
      range: range,
      options: {
        isWholeLine: true,
        className: 'non-selected-code'
      }
    }]);
    this._markedInfoService.addDecoration(range, decorates[0], this._$currentFile.fileName);
    this._zone.run(() => {
      this._showAnswerCommentTip(range);
      this._showCommentZone(range);
    });
  }

  /**
   * 设置选中的代码段的底色
   * @param range
   * @param fileName
   */
  private setSelectedAreaColor(range: Range, fileName = this._$currentFile.fileName): void {
    const marked = this._markedInfoService.getDecoration(range, this._$currentFile.fileName);
    this._codeEditor.editor.deltaDecorations([marked.decoration], [{
      range: range,
      options: {
        isWholeLine: true, // 是否整行选中
        className: 'selected-code'  // 选中行的类名，以关联选中底色
      }
    }]);
  }

  /**
   *
   * @param range 清除选中的代码段的底色
   * @param fileName
   */
  private clearSelectedAreaColor(range: Range, fileName = this._$currentFile.fileName): void {
    const marked = this._markedInfoService.getDecoration(range, this._$currentFile.fileName);
    this._codeEditor.editor.deltaDecorations([marked.decoration], [{
      range: range,
      options: {
        isWholeLine: true,
        className: 'non-selected-code'
      }
    }]);
  }

  /***
   * 显示添加评审意见的按钮
   */
  private _showAddCommentButton(range: Range): void {
    const lineNr = range.endLineNumber;
    const colNr = range.endColumn;
    const contentWidget = {
      getId: () => {
        return `target-practice.add.widget-${lineNr}-${colNr}`;
      },
      getDomNode: () => {
        const componentRef = this._domService.getComponentRef(AddComponent, {
          range: range, fileName: this._$currentFile?.fileName
        }, this._commentedCallback.bind(this));
        return this._domService.getDomElement(componentRef);
      },
      getPosition: () => {
        return {
          position: {
            lineNumber: lineNr,
            column: colNr
          },
          preference: [
            monaco.editor.ContentWidgetPositionPreference.ABOVE,
            monaco.editor.ContentWidgetPositionPreference.BELOW
          ]
        };
      }
    };
    this._codeEditor.editor.addContentWidget(contentWidget);
    this._markedInfoService.addWidget(range, contentWidget, this._$currentFile.fileName, false);
  }

  /**
   * 这个工具条已经废弃不再使用
   * 显示编辑评审意见工具条(添加完评审后去掉添加按钮，显示出编辑和删除按钮，并把评审意见显示出来)
   * @param range
   */
  private _showEditCommentTools(range: Range): void {
    const lineNr = range.startLineNumber;
    const colNr = range.startColumn;
    const contentWidget = {
      getId: () => {
        return `target-practice.edit.widget-${lineNr}-${colNr}`;
      },
      getDomNode: () => {
        const componentRef = this._domService.getComponentRef(EditComponent, {
          range: range, fileName: this._$currentFile?.fileName, shootType: this._shootType
        }, this._commentedCallback.bind(this));
        return this._domService.getDomElement(componentRef);
      },
      getPosition: () => {
        return {
          position: {
            lineNumber: lineNr,
            // column: colNr
            column: 120
          },
          preference: [
            monaco.editor.ContentWidgetPositionPreference.ABOVE,
            monaco.editor.ContentWidgetPositionPreference.BELOW
          ]
        };
      }
    };
    this._codeEditor.editor.addContentWidget(contentWidget);
    this._markedInfoService.addWidget(range, contentWidget, this._$currentFile.fileName, true);
    this._showCommentTip(range);
  }

  /**
   * 在所选择的代码段的下方增加一个视图区（不影响原来代码行），参考gerrit，同时把评审信息窗体也显示到该区域
   * @param range
   */
  private _showCommentZone(range: Range): void {
    var viewZoneId = null;
    let that = this;
    this._codeEditor.editor.changeViewZones(function (changeAccessor) {
      var domNode = document.createElement('div' + range.startLineNumber + range.startColumn + range.endLineNumber + range.endColumn);
      domNode.style.background = '#E6F5FF';
      viewZoneId = changeAccessor.addZone({
        afterLineNumber: range.endLineNumber,
        heightInLines: 1,
        domNode: domNode
      });
      // 把viewZoneId保存起来，以备修改或删除时使用
      that._markedInfoService.addViewZone(range, viewZoneId, that._$currentFile.fileName);
    });
  }

  /**
   * 展示评审信息提示
   * @param range
   */
  private _showCommentTip(range: Range): void {
    const lineNr = this.getRowByMaxColumn(range);
    const colNr = this._codeEditor.editor.getModel().getLineMaxColumn(lineNr);
    const endLineNumber = range.endLineNumber;
    const contentWidget = {
      getId: () => {
        return `target-practice.tip.widget-${range.startLineNumber}-${range.startColumn}-${range.endLineNumber}-${range.endColumn}`;
      },
      getDomNode: () => {
        const componentRef = this._domService.getComponentRef(CommentComponent, {
          range: range, fileName: this._$currentFile?.fileName, shootType: this._shootType
        }, this._commentedCallback.bind(this));

        let domeNode = this._domService.getDomElement(componentRef);
        domeNode.setAttribute('class', 'comment-msg-' + endLineNumber.toString());
        let length = document.getElementsByClassName('comment-msg-' + endLineNumber.toString()).length;
        if (length > 0) {
          domeNode.style['margin-top'] = (length * 18) + 'px';
        }
        return domeNode;
      },
      getPosition: () => {
        return {
          position: {
            lineNumber: lineNr,
            column: 1
          },
          preference: [
            monaco.editor.ContentWidgetPositionPreference.BELOW
          ]
        };
      }
    };
    this._codeEditor.editor.addContentWidget(contentWidget);
    this._markedInfoService.addCommentInfo(range, contentWidget, this._$currentFile.fileName);
  }

  /**
   * 展示答案信息提示
   * @param range
   */
  private _showAnswerCommentTip(range: Range): void {
    const lineNr = this.getRowByMaxColumn(range);
    //const colNr = this._codeEditor.editor.getModel().getLineMaxColumn(lineNr);
    const endLineNumber = range.endLineNumber;
    const contentWidget = {
      getId: () => {
        return `target-practice.tip.widget-${range.startLineNumber}-${range.startColumn}-${range.endLineNumber}-${range.endColumn}`;
      },
      getDomNode: () => {
        const componentRef = this._domService.getComponentRef(AnswerComponent, {
          range: range, fileName: this._$currentFile?.fileName, shootType: this._shootType
        }, this._commentedCallback.bind(this));

        let domeNode = this._domService.getDomElement(componentRef);
        domeNode.setAttribute('class', 'comment-msg-' + endLineNumber.toString());
        let length = document.getElementsByClassName('comment-msg-' + endLineNumber.toString()).length;
        if (length > 0) {
          domeNode.style['margin-top'] = (length * 18) + 'px';
        }
        return domeNode;
      },
      getPosition: () => {
        return {
          position: {
            lineNumber: lineNr,
            column: 1
          },
          preference: [
            monaco.editor.ContentWidgetPositionPreference.BELOW
          ]
        };
      }
    };
    this._codeEditor.editor.addContentWidget(contentWidget);
    this._markedInfoService.addCommentInfo(range, contentWidget, this._$currentFile.fileName);
  }

  public getRowByMaxColumn(range: Range): number {
    if (range.startLineNumber == range.endLineNumber) {
      return range.startLineNumber;
    }
    let currentColNr: number = 0, currentRowNr: number = 0;
    Array(range.endLineNumber - range.startLineNumber).fill(0)
      .forEach((item, index) => {
        const colNr = this._codeEditor.editor.getModel().getLineMaxColumn(range.startLineNumber + index);
        if (colNr > currentColNr) {
          currentColNr = colNr;
          // currentRowNr = range.startLineNumber + index;
          currentRowNr = range.endLineNumber;
        }
      });
    return currentRowNr;
  }

  /**
   * 确认添加或删除评审意见后的回调
   * @param target
   */
  private _commentedCallback(target: TargetInfo): void {
    // 处理鼠标移入移出评审意见区域时的代码背景色问题
    if (target['mouseAction']) {
      if (target['mouseAction'] == 'over') {
        this.setSelectedAreaColor(target['range']);
      } else {
        this.clearSelectedAreaColor(target['range']);
      }
      return;
    }
    const range = this._defectService.toRange(target);
    const marked = this._markedInfoService.getDecoration(range, this._$currentFile.fileName);
    this._clearDecoration(marked);
    const existTarget = this._defectService.getTarget(this._defectService.toRange(target), target.FileName);
    if (existTarget) {
      this._decorateSelectedArea(range, false);
    }
    // 自动调用 一下保存草稿接口
    this.curShootResult = this._defectService.shootingResult;
    this._$save();
    // 如果对已有的评审意见进行了编辑或删除，则移除当前窗体后重新渲染呈现评审信息
    if (target['isEditOrDelete'] == true) {
      let sameLineTarget = this._defectService.getSameEndLineTarget(this._defectService.toRange(target), target.FileName);
      if (sameLineTarget.length > 0) {
        sameLineTarget.forEach(targetItem => {
          let range = this._defectService.toRange(targetItem);
          let markItem = this._markedInfoService.getDecoration(range, this._$currentFile.fileName);
          if (markItem.commentInfo) {
            this._codeEditor.editor.removeContentWidget(markItem.commentInfo);
          }
        });
        sameLineTarget.forEach(targetItem => {
          let range = this._defectService.toRange(targetItem);
          this._showCommentTip(range);
        });
      }
    }
  }

  /**
   * 清除选中区域的装饰
   * @param fileName
   */
  private _clearSelectedArea(fileName = this._$currentFile.fileName): void {
    this._markedInfoService.getDecorations(fileName).forEach(item => this._clearDecoration(item));
    this._markedInfoService.clearDecorations(fileName);
  }

  /**
   * 移除添加评审的按钮（如果点击处是添加按钮则不移除）
   * @param clickedElement
   */
  private _removeAddCommentButton(clickedElement?: any): void {
    const info = this._markedInfoService.getDecorations(this._$currentFile.fileName).filter(item => item.done == false);
    info.forEach(item => {
      if (!clickedElement || clickedElement.parentNode.getAttribute("class") !== "comment-add-tools") {
        this._clearDecoration(item);
      }
    });
  }

  private _clearDecoration(item: DecorateInfo): void {
    let that = this;
    if (!item) {
      return;
    }
    this._codeEditor.editor.deltaDecorations([item.decoration], []);
    if (item.widget) {
      this._codeEditor.editor.removeContentWidget(item.widget);
    }
    if (item.commentInfo) {
      this._codeEditor.editor.removeContentWidget(item.commentInfo);
    }
    if (item.viewzone) {
      this._codeEditor.editor.changeViewZones(function (changeAccessor) {
        changeAccessor.removeZone(item.viewzone);
      });
    }
    this._markedInfoService.delDecoration(item.range, this._$currentFile.fileName);
  }

  /**
   * 左侧文件树节点变动
   * @param fileTreeNode
   * @returns
   */
  public _$fileNodeChange(fileTreeNode: FileTreeNode): void {
    if (!fileTreeNode) {
      return;
    }
    let index = this._$tabTitles.findIndex(item => item.fileName == fileTreeNode.fileName);
    if (index == -1) {
      index = this._$tabTitles.push(this._getTabTitle(fileTreeNode.fileName)) - 1;
      this._openedFiles.push(fileTreeNode);
    }
    this._$fileSelectedIndex = index;
    this._setSelectFile(fileTreeNode);
  }

  /**
   * 编辑器顶部tab切换
   * @param index
   * @returns
   */
  public _$selectedIndexChange(index: number): void {
    const title = this._$tabTitles[index];
    const file = this._openedFiles.find(item => item.fileName == title?.fileName);
    if (!file) {
      return;
    }
    this._setSelectFile(file);
  }

  private _setSelectFile(fileTreeNode: FileTreeNode): void {
    this.hasRunRender = false;
    this._clearSelectedArea(this._currentFileName);
    this._$currentFile = fileTreeNode;
    this._currentFileName = this._$currentFile?.fileName || '';

    if (!this._$currentFile) {
      this._$code = '';
      return;
    }
    if (fileTreeNode.code) {
      this._$code = fileTreeNode.code;
      setTimeout(() => this._renderDefects(this._$currentFile.fileName));
    }
  }

  /**
   * 获取右侧编辑器顶部Tab页的标题名称
   * @param fileName
   * @returns
   */
  private _getTabTitle(fileName: string): TabTitle {
    return {
      fileName: fileName,
      html: `
                <span>${fileName}</span>
                <span (click)="this._$closeFile('${fileName}')" class="iconfont iconfont-e923" style="width: 20px; position: relative; right: -5px; top: 2px;"></span>
            `,
      htmlContext: this
    }
  }

  public _$closeFile(fileName: string): void {
    const index = this._openedFiles.findIndex(item => item.fileName == fileName);
    if (index == -1) {
      return;
    }
    this._openedFiles.splice(index, 1);
    this._$tabTitles.splice(index, 1);
    if (index > this._$fileSelectedIndex) {
      return;
    }
    if (index == this._$fileSelectedIndex) {
      this._$fileSelectedIndex = 0;
      this._setSelectFile(this._openedFiles[0]);
      return;
    }
    this._$fileSelectedIndex -= 1;
  }

  /**
   * 保存当前打靶草稿，以防止网页突然被关闭
   */
  public _$save() {
    let userId: string = localStorage.getItem('user');
    let userName: string = localStorage.getItem('name');
    let reqBody: any = {
      'userid': userId,
      'username': userName,
      'targetid': this._targetId,
      'rangeid': this._rangeId,
      'targets': []
    };
    let targetComment = this._defectService.shootingResult.targets;
    reqBody.targets = targetComment.map(item => {
      return {
        filename: item.fileName,
        startlinenum: item.startLineNum,
        endlinenum: item.endLineNum,
        startcolnum: item.startColNum,
        endcolnum: item.endColNum,
        defectclass: item.defectClass,
        defectsubclass: item.defectSubClass,
        defectdescribe: item.defectDescribe,
        remark: item.remark
      }
    });
    this.manageService.saveTargetDraft(reqBody).subscribe(res => {
      // console.log("草稿保存成功");
    }, error => {
      console.log("草稿保存失败");
    });
  }

  /**
   * 进入下一个靶子进行打靶
   */
  public _$gotoNextTarget() {

  }
}

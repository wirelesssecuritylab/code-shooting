import {
    AfterViewInit,
    ChangeDetectorRef,
    Component,
    EventEmitter,
    Input,
    NgModule,
    NgZone,
    OnDestroy,
    Output,
    Renderer2
} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FormsModule} from '@angular/forms';
import {PopupEffect, PopupInfo, PopupService} from '@rdkmaster/jigsaw';
import {MonacoCodeBase} from "./code-base";
import {Utils} from "../../misc/utils";

export type CodeLineInfo = {
    lineNumber: number,
    message: string
}

declare const monaco;

@Component({
    selector: 'tp-code',
    template: `
        <div #editorContainer style="position:relative;height:{{autoHeight}}">
            <a class="fullscreen-icon iconfont iconfont-e20e" *ngIf="fullscreenSwitchable" title="进入全屏编辑模式"
               (click)="popupFullscreenEditor(popupTarget)" style="display: none"></a>
        </div>

        <ng-template #popupTarget>
            <div style="position:relative; border: 4px solid #777; background-color: #1e1e1e;">
                <tp-code [(code)]="codeBackup" [language]="language" [fullscreenSwitchable]="false"
                          [readonly]="readonly" [options]="options" [environment]="environment"
                          width="95vw" height="95vh"></tp-code>
                <a class="fullscreen-icon iconfont iconfont-e20f" style="color:#eee" title="返回正常编辑模式"
                   (click)="closeFullscreenEditor()"></a>
            </div>
        </ng-template>
    `,
    styles: [`
        .fullscreen-icon {
            position: absolute;
            color: #aaa;
            right: 16px;
            top: 0;
            font-size: 28px;
            z-index: 1;
        }
    `]
})
export class MonacoCodeEditor extends MonacoCodeBase implements AfterViewInit, OnDestroy {
    public codeBackup: string = '';
    private _code: string = '';

    @Input()
    public get code(): string {
        return this._code;
    }

    public set code(value: string) {
        value = Utils.isDefined(value) ? value : "";
        if (typeof value !== "string") {
            value = JSON.stringify(value, null, '  ');
        }
        if (value === this.code) {
            return;
        }
        this._code = value;
        if (this.editor) {
            this.editor.setValue(this._code);
        }
    }

    @Input()
    public fullscreenSwitchable: boolean = true;

    @Output()
    public codeChange = new EventEmitter<string>();

    @Output()
    public blur = new EventEmitter<string>();

    @Output()
    public ready = new EventEmitter<void>();

    constructor(protected _renderer: Renderer2, private _zone: NgZone, private _popupService: PopupService, private _cdr: ChangeDetectorRef) {
        super(_renderer);
    }

    public get autoHeight(): string {
        return this.height ? this.height : (this._code.split('\n').length + 2) * 19 + 'px';
    }

    protected initMonaco() {
        this.editorOptions.value = this.code;
        super.initMonaco();
        this._zone.runOutsideAngular(() => {
            this.editor = monaco.editor.create(this._editorContainer.nativeElement, this.editorOptions);
        });
        this.editor.onDidChangeModelContent((e: any) => {
            this._code = this.editor.getValue();
            this.codeChange.emit(this._code);
            this._cdr.detectChanges();
        });
        this.editor.onDidBlurEditorWidget(() => {
            this.blur.emit();
        });
        this.ready.emit();
        if (this._errorLineInfo) {
            this.jumpToErrorLine(this._errorLineInfo);
        }
    }

    private _popupInfo: PopupInfo;

    public popupFullscreenEditor(target: any) {
        if (this._popupInfo) {
            this._popupInfo.dispose();
        }
        this.codeBackup = this.code;
        this._popupInfo = this._popupService.popup(target, {
            modal: true, useCustomizedBackground: true,
            showEffect: PopupEffect.fadeIn,
            hideEffect: PopupEffect.fadeOut
        });
    }

    public closeFullscreenEditor() {
        if (!this._popupInfo) {
            return;
        }
        this._popupInfo.dispose();
        this._popupInfo = null;
        this.code = this.codeBackup;
        this.codeBackup = '';
    }

    private _decorationIds: string[];

    private _errorLineInfo: CodeLineInfo;

    public jumpToErrorLine(errorInfo: CodeLineInfo) {
        if (!errorInfo) {
            return;
        }
        let {lineNumber, message} = errorInfo;
        lineNumber = Number(lineNumber);
        if (isNaN(lineNumber) || lineNumber <= 0) {
            return;
        }
        setTimeout(() => {
            // 等待编辑器绘制完成
            if (!this.editor) {
                // 保存错误信息，等到editor被创造再跳转
                this._errorLineInfo = errorInfo;
                return;
            }
            const model = this.editor.getModel();
            const lineContent = model.getLineContent(lineNumber);
            this.editor.revealLinesInCenter(lineNumber, lineNumber);
            this._decorationIds = this.editor.deltaDecorations([], [
                {
                    range: {
                        startLineNumber: lineNumber,
                        startColumn: 0,
                        endLineNumber: lineNumber,
                        endColumn: lineContent ? lineContent.length : 1000,
                    },
                    options: {
                        inlineClassName: 'code-editor-error-emphasis',
                        isWholeLine: true,
                        hoverMessage: {
                            value: message,
                            isTrusted: true
                        }
                    }
                },
            ]);
            this._errorLineInfo = null;
        })
    }

    ngOnDestroy(): void {
        super.ngOnDestroy();
        if (this._popupInfo) {
            this._popupInfo.dispose();
        }
    }
}

@NgModule({
    imports: [
        CommonModule, FormsModule
    ],
    declarations: [MonacoCodeEditor],
    exports: [MonacoCodeEditor]
})
export class MonacoCodeEditorModule {
}

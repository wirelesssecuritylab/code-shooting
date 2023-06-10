import {AfterViewInit, ElementRef, Input, OnDestroy, Renderer2, ViewChild, Component} from '@angular/core';
import {Utils} from "../../misc/utils";

declare const monaco;
let loadedMonaco = false;
let loadPromise: Promise<void>;

@Component({
  template: ''
})

export class MonacoCodeBase implements AfterViewInit, OnDestroy {
    public editor: any;
    private _width: string;
    protected _language: string;

    @ViewChild('editorContainer')
    public _editorContainer: ElementRef;

    /**
     * monaco-editor的options 请参考https://github.com/Microsoft/monaco-editor使用
     */
    @Input()
    public options: any;

    @Input()
    public get language(): string {
        return this._language;
    }

    public set language(value: string) {
        this._language = this._getLanguageValue(value);
        if (!loadedMonaco || !this.editor) {
            return;
        }
        monaco.editor.setModelLanguage(this.editor.getModel(), this._language);
    }

    private _readonly: boolean | string;
    @Input()
    public get readonly(): boolean | string {
        return this._readonly
    }

    public set readonly(value: boolean | string) {
        if (this._readonly == value) {
            return;
        }
        this._readonly = value;
        if (this.editor) {
            this.editor.updateOptions({readOnly: value});
        }
    }

    @Input()
    public environment: 'frontend' | 'backend';

    @Input()
    public get width(): string {
        return this._width;
    }

    public set width(value: string) {
        this._width = Utils.getCssValue(value);
    }

    @Input()
    height: string | undefined;

    constructor(protected _renderer: Renderer2) {
    }

    ngAfterViewInit(): void {
        this.editorOptions = {
            // 语言环境
            language: this.language,
            // 只读模式
            readOnly: this.readonly,
            // 代码长度标尺
            rulers: [120],
            // 是否换行
            wordWrap: 'off',
            // 自动Layout
            automaticLayout: true,
            // 迷你地图
            minimap: {enabled: false},
            scrollbar: {vertical: this.height ? 'visible' : 'hidden', handleMouseWheel: Boolean(this.height)},
            scrollBeyondLastLine: false
        };
        if (loadedMonaco) {
            loadPromise.then(() => {
                this.initMonaco();
            });
        } else {
            loadedMonaco = true;
            loadPromise = new Promise<void>((resolve: any) => {
                if (typeof ((<any>window).monaco) === 'object') {
                    resolve();
                    return;
                }
                const baseUrl = './assets/monaco-editor/min/vs';
                const onGotAmdLoader: any = () => {
                    (<any>window).require.config({paths: {'vs': `${baseUrl}`}});
                    (<any>window).require(['vs/editor/editor.main'], () => {
                        this.initMonaco();
                        resolve();
                    });
                };
                if (!(<any>window).require) {
                    const loaderScript: HTMLScriptElement = document.createElement('script');
                    loaderScript.type = 'text/javascript';
                    loaderScript.src = `${baseUrl}/loader.js`;
                    loaderScript.addEventListener('load', onGotAmdLoader);
                    document.body.appendChild(loaderScript);
                } else {
                    onGotAmdLoader();
                }
            });
        }
    }

    public editorOptions: any;

    /**
     * 解析父组件传值的语言类型，当前支持语言
     * 'apex', 'azcli', 'bat', 'clojure', 'coffee', 'cpp', 'csharp', 'csp', 'css', 'dockerfile', 'fsharp', 'go',
     *  'handlebars', 'html', 'ini', 'java', 'javascript', 'json', 'less', 'lua', 'markdown', 'msdax',
     *  'mysql', 'objective', 'perl', 'pgsql', 'php', 'postiats', 'powerquery', 'powershell', 'pug', 'python',
     *  'r', 'razor', 'redis', 'redshift', 'ruby', 'rust', 'sb', 'scheme', 'scss', 'shell', 'solidity', 'sql', 'st',
     *  'swift', 'typescript', 'vb', 'xml', 'yaml'
     */
    protected _getLanguageValue(key: string): string {
        let res: string;
        switch (key) {
            case 'js':
                res = 'javascript';
                break;
            case 'ts':
                res = 'typescript';
                break;
            default:
                res = key;
                break;
        }
        return res;
    }

    ngOnDestroy(): void {
        if (this.editor) {
            this.editor.dispose();
            this.editor = null;
        }
    }

    protected initMonaco() {
        if (this.options) {
            Object.assign(this.editorOptions, this.options);
        }
        // this.editorOptions['theme'] =  'vs-dark';
        this._renderer.setStyle(this._editorContainer.nativeElement, 'width', this.width);
        if (this.height) {
            this._renderer.setStyle(this._editorContainer.nativeElement, 'height', this.height);
        }
    }
}

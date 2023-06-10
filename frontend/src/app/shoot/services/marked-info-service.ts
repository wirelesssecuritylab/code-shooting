import {Injectable} from '@angular/core';
import {DecorateInfo, Range} from "../misc/types";

type MarkedInfo = {
    fileName: string,
    decorations: Map<string, DecorateInfo>
}

@Injectable()
export class MarkedInfoService {
    private _markedInfos: MarkedInfo[] = [];

    public addDecoration(range: Range, decoration: any, fileName: string): void {
        const markedInfo = this._getMarkedInfo(fileName);
        markedInfo.decorations.set(this._getMarkedKey(range), {range: range, decoration: decoration})
    }

    public addWidget(range: Range, widget: any, fileName: string, done?: boolean): void {
        const markedInfo = this._getMarkedInfo(fileName);
        const markedTarget = markedInfo.decorations.get(this._getMarkedKey(range));
        markedTarget.widget = widget;
        markedTarget.done = done;
    }

    public addCommentInfo(range: Range, commentInfo: any, fileName: string): void {
        const markedInfo = this._getMarkedInfo(fileName);
        const key = this._getMarkedKey(range);
        const markedTarget = markedInfo.decorations.get(this._getMarkedKey(range));
        markedTarget.commentInfo = commentInfo;
    }

    /**
     * 把viewzone的ID保存起来
     * @param range
     * @param commentInfo
     * @param fileName
     */
    public addViewZone(range: Range, viewZoneId: any, fileName: string): void {
      const markedInfo = this._getMarkedInfo(fileName);
      const markedTarget = markedInfo.decorations.get(this._getMarkedKey(range));
      markedTarget.viewzone = viewZoneId;
    }

    /**
     *
     * @param range 获取装饰信息：窗体：按钮和评审信息
     * @param fileName
     * @returns
     */
    public getDecoration(range: Range, fileName: string): DecorateInfo {
        const markedInfo = this._getMarkedInfo(fileName);
        if (!markedInfo) {
            return null;
        }
        const key = this._getMarkedKey(range);
        return markedInfo.decorations.get(key);
    }

    /**
     * 删除文件内指定区域内的装饰
     * @param range
     * @param fileName
     * @returns
     */
    public delDecoration(range: Range, fileName: string): void {
        const markedInfo = this._getMarkedInfo(fileName);
        if (!markedInfo) {
            return null;
        }
        const key = this._getMarkedKey(range);
        markedInfo.decorations.delete(key);
    }

    /**
     * 获取某个文件中所有的markinfo
     * @param fileName
     * @returns
     */
    public getDecorations(fileName: string): DecorateInfo[] {
        const markedInfo = this._getMarkedInfo(fileName);
        if (!markedInfo) {
            return [];
        }
        return Array.from(markedInfo.decorations.values());
    }

    public clearDecorations(fileName: string): void {
        const markedInfo = this._getMarkedInfo(fileName);
        if (!markedInfo) {
            return;
        }
        markedInfo.decorations.clear();
    }

    /**
     * 获取某个文件已添加的评审信息
     * @param fileName
     * @returns
     */
    private _getMarkedInfo(fileName: string): MarkedInfo {
        let markedInfo = this._markedInfos.find(item => item.fileName == fileName);
        if (!markedInfo) {
            markedInfo = {
                fileName: fileName,
                decorations: new Map<string, DecorateInfo>()
            }
            this._markedInfos.push(markedInfo);
        }
        return markedInfo;
    }

    private _getMarkedKey(range: Range): string {
        return `${range.startLineNumber}-${range.startColumn}-${range.endLineNumber}-${range.endColumn}`;
    }
}

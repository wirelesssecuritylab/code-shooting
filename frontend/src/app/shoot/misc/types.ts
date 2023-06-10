import {TargetInfo} from "./target-types";

export type Position = {
    column: number,
    lineNumber: number
}

export type FileInfo = {
    fileName: string,
    code: string
}

export type Range = {
    startLineNumber: number;
    startColumn?: number;
    endLineNumber: number;
    endColumn?: number;
}

export type InitData = {
    range: Range,
    fileName?: string,
    defectTarget?: TargetInfo,
    shootType?: string
}

export type DefectClass = {
    DefectClass: string,
    defectSubClasses?: DefectSubClass[]
}

export type DefectSubClass = {
    DefectSubClass: string,
    describes?: any[]
}

export type FileTreeNode = {
    fileName: string,
    open?: boolean;
    type?: 'src' | 'img' | 'md';
    code?: string,
    imageSrc?: string,
    iconUnicode?: string,
    nodes?: FileTreeNode[]
}

/**
 * 装饰的结构体类型
 *
 */
export type DecorateInfo = {
    range?: Range,
    done?: boolean,
    decoration?: any,
    widget?: any,
    commentInfo?: any,
    viewzone?: any
};

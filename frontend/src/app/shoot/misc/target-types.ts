export class TargetResult {
    targetId: string;
    fileName?: string;
    startLineNum?: number;
    endLineNum?: number;
    defectClass: string;
    defectSubClass: string;
    defectDescribe: string;
    remark: string;
    startColNum?: number;
    endColNum?: number;

    constructor(target: TargetInfo) {
        this.targetId = target.TargetId;
        this.fileName = target.FileName;
        this.startLineNum = target.StartLineNum;
        this.endLineNum = target.EndLineNum;
        this.defectClass = target.DefectClass;
        this.defectSubClass = target.DefectSubClass;
        this.defectDescribe = target.DefectDescribe;
        this.remark = target.Remark;
        this.startColNum = target.StartColNum;
        this.endColNum = target.EndColNum;
    }
}

export class ShootingResult {
    userName: string;
    userId: string;
    targets: TargetResult[];

    constructor(userName: string, userId: string) {
        this.userName = userName;
        this.userId = userId;
        this.targets = [];
    }
}

export type TargetInfo = {
    TargetId: string,
    FileName?: string,
    StartLineNum?: number,
    StartColNum?: number,
    EndLineNum?: number,
    EndColNum?: number,
    DefectClass: string,
    DefectSubClass: string,
    DefectDescribe: string,
    Remark: string
}

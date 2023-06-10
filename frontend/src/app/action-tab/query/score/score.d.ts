export interface ScoreResponse {
    hitNum: string;
    hitScore: string;
    targets: Target[]
}

export interface Target {
    filename: string;
    startLineNum: number;
    endLineNum: number;
    defectClass: string;
    defectSubClass: string;
    defectDescribe: string;
    score: number;
}
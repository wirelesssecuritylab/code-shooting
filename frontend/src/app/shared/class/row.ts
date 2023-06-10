export class RowParams {
  public targetId!: string | null;
  public type!: string | null;
  public languages!: string[] | null;
  public start_at!: Date | null;
  public end_at!: Date | null;

  constructor() {

  }

  getParams() {
    return {
      targetId: this.targetId,
      type: this.type,
      languages: this.languages,
      start_at: this.start_at,
      end_at: this.end_at,
    };
  }

  setParams(row: any) {

    this.targetId = row.id;
    this.type = row.type;
    this.languages = row.languages;
    this.start_at = row.start_at;
    console.log(this.getParams())
  }

  reset() {
    this.targetId = null;
    this.type = null;
    this.languages = null;
    this.start_at = null;
    this.end_at = null;
  }
}

export const typeAdapter = [
  { index: 0, text: '练习', value: 'test'},
  { index: 1, text: '比赛', value: 'compete'}
];

export const languageAdapter = {
  Go: 'go',
  'C/C++': 'cpp',
  C: 'c',
  Python: 'python',
  Java:'java',
};

export const rangeType = {
  test: '练习',
  compete: '比赛'
};

export const langReverseAdapter = {
  go : 'Go',
  cpp: 'C/C++',
  c: 'C',
  python: 'Python',
  java: 'Java'
};

export const errCodeAdapter = {
  1001: '参数错误',
  1002: '文件系统访问失败',
  1003: '数据库访问失败',
  1004: '记录不存在',
  1099: '未知错误',
  1101: '用户不存在',
  1102: '认证失败',
  1103: '权限拒绝',
  1201: '答卷文件不合法',
  1202: '靶标不存在',
  1203: '阅卷失败，请确保文件格式正确或提前解密',
};

/**
 * 新建靶子的请求body体
 */
export class AddTargetBody {
  public workspace: string;
  public name: string;
  public owner: string;
  public ownerName: string;
  public language: string;
  public isShared: boolean;
  public targets: string[];
  public template: string;
  // public tagId: string; //
  /**
   * 标签对象,结构如下：里面的key和value都是直接使用中文
   * tagName: {
   *   缺陷大类：性能,
   *   缺陷小类：所有，
   *   缺陷细项：所有
   * }
   */
  public tagName: any = {
    'mainCategory': '所有',
    'subCategory': '所有',
    'defectDetail': '所有'
  }
  public extendedLabel: string[];
  public customLable: string;
  public instituteLabel: string[];
  public answer: string;
}

/**
 * 用于表示添加或编辑靶子表单时的默认值结构体
 */
export class TargetFormField {
  public workspace: string = '';
  public name: string = '';
  public language: string = '';
  public template: string = '';
  public isShared: boolean = false;
  public targets: string = '';
  public bugType: string = '所有';
  public subBugType: string = '所有';
  public bugDetail: string = '所有';
  public targetTag: string[] = [];
  public answer: string = '';
  public extendedLabel: string[] = [];
  public customLable: string = '';
  public instituteLabel: string[] = [];
}

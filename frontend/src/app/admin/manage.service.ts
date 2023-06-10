import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

const DPTURL = '/api/code-shooting/v1/project/';
const SCOREURL = '/api/code-shooting/v1/results/range/';
const TARGETURL = '/api/code-shooting/v1/actions/target';
const DEFECTURL = '/api/code-shooting/v1/actions/target/';
const RANGEURL = '/api/code-shooting/v1/actions/range';
const PROJECTURL = '/api/code-shooting/v1/projects';
const TEMPLATEURL = '/api/code-shooting/v1/actions/template';
@Injectable({
  providedIn: 'root'
})
export class ManageService {

  constructor(private http: HttpClient) { }

  /**
   * 根据项目id获取其下的部门列表
   * @param projectId 项目id
   * @returns
   */
  public getDepartListApi(projectId: string): Observable<any> {
    let url: string = DPTURL + projectId + '/depts';
    return this.http.get(url);
  }


  /**
   * 查询成绩
   * @param rangeId 靶场ID
   * @param langType 打靶语言
   * @param params 其它过滤参数如部门
   * @returns
   */
  public queryScoreApi(rangeId: string, langType: string, params: any): Observable<any> {
    let url: string = SCOREURL + rangeId + '/language/' + langType;
    let paramArray: string[] = [];
    for (const key of Object.keys(params)) {
      paramArray.push(`${key}=${params[key]}`);
    }
    if (paramArray.length > 0) {
      url += '?' + paramArray.join('&')
    }
    return this.http.get(url);
  }

  /**
   * 导出excel成绩
   * @param rangeId 靶场ID
   * @param langType 打靶语言
   * @param params 其它过滤参数如部门
   * @returns
   */
  public exportScoreApi(rangeId: string, langType: string, params: any): Observable<any> {
    let url: string = SCOREURL + rangeId + '/language/' + langType + '/excel';
    let paramArray: string[] = [];
    for (const key of Object.keys(params)) {
      paramArray.push(`${key}=${params[key]}`);
    }
    if (paramArray.length > 0) {
      url += '?' + paramArray.join('&')
    }
    return this.http.get(url, { responseType: 'blob' });
  }

  /**
   * 操作靶子：增删改查
   * @param body
   * @returns
   */
  public manageTargetApi(body: any): Observable<any> {
    return this.http.post(TARGETURL, body);
  }

  /**
   * 批量导出靶子API
   * @param data
   */
  public exportBatchTargetApi(data: any, role: string, userId: string): Observable<any> {
    let url: string = TARGETURL + '/exportbatchtarget';
    var reqBody: any;
    let param: string = '';
    if (role == 'admin') {
      for (var i = 0; i < data.length; i++) {
        param += data[i].id + ","
      }
    } else {
      for (var i = 0; i < data.length; i++) {
        if (data[i].owner == userId) {
          param += data[i].id + ","
        }
      }
    }
    reqBody = {
      targetIds: param,
    };
    return this.http.post(url, reqBody, { responseType: 'blob' });
  }

  /**
   * 获取指定语言下的编码信息(即规范信息)
   * @param body
   * @returns
   */
  public getDefectApi(targetId: string, body: any): Observable<any> {
    return this.http.post(DEFECTURL + targetId + "/defect", body);
  }

  /**
   * 操作靶场api:增删改查
   * @param body 请求body体
   * @returns
   */
  public manageRangeApi(body: any): Observable<any> {
    return this.http.post(RANGEURL, body);
  }

  /**
   * 获取项目列表
   * @param userId
   * @returns
   */
  public getProjectListApi(userId: string): Observable<any> {
    let url: string = userId ? (PROJECTURL + '?userId=' + userId) : PROJECTURL;
    return this.http.get(url);
  }

  /**
   * 规范管理api:启用/停用/删除/查询
   * @param body 请求body体
   * @returns
   */
  public manageTemplateApi(body: any): Observable<any> {
    return this.http.post(TEMPLATEURL, body);
  }

  public downloadTemplateApi(body: any): Observable<any> {
    return this.http.post(TEMPLATEURL, body, { responseType: 'blob' });
  }

  public getWorkspace(): Observable<any> {
    return this.http.get('/api/code-shooting/v1/workspaces');
  }

  /**
   * 获取靶子下某个靶子文件的内容
   * 20230201郭康旭10243452 由于文件存在特殊字符，如"#"，修改使用URL的Get方式为post获取文件内容
   * @param targetId 靶子id
   * @param fileName 靶子文件名
   * @returns
   */
  public getTargetFileApi(targetId: string, fileName: string): Observable<any> {
    let url: string = '/api/code-shooting/v1/targets/' + targetId + '/files';
    // 因为该接口返回的响应头中的Content-Type: text/plain，所以我们这里get请求时要携带{responseType: 'text'}参数
    return this.http.post(url, { "filename": fileName }, { responseType: 'text' });
  }

  /**
   * 提交打靶答案
   * @param rangeId
   * @param language
   * @param reqBody
   * @returns
   */
  public submitTargetAnswer(rangeId: string, language: string, reqBody: any): Observable<any> {
    let url: string = '/api/code-shooting/v1/answers/submit?rangeId=' + rangeId + '&language=' + language;
    return this.http.post(url, reqBody);
  }

  /**
   * 保存打靶草稿， 打靶自动保存时使用
   * @param reqBody
   * @returns
   */
  public saveTargetDraft(reqBody: any): Observable<any> {
    const url: string = '/api/code-shooting/v1/answers/savedraft';
    return this.http.post(url, reqBody);
  }

  /**
   * 获取靶场答案
   * @param targetId
   * @param rangeId
   * @returns
   */
  public getTargetsAnswer(targetId: string, rangeId: string): Observable<any> {
    let url: string = `/api/code-shooting/v1/targets/${targetId}/answers/${rangeId}`;
    return this.http.get(url);
  }

  /**
   * 获取打靶草稿,继续打靶时使用
   * @param userId
   * @param targetId
   * @returns
   */
  public getShootDraftApi(userId: string, targetId: string, rangeId: string): Observable<any> {
    let url: string = `/api/code-shooting/v1/answers/loaddraft?userId=${userId}&targetId=${targetId}&rangeId=${rangeId}`;
    return this.http.get(url);
  }

  /**
   * 获取打靶记录,查看答卷时使用
   * @param userId
   * @param targetId
   * @returns
   */
  public getShootRecordsApi(userId: string, targetId: string, rangeId: string): Observable<any> {
    let url: string = `/api/code-shooting/v1/answers/load?userId=${userId}&targetId=${targetId}&rangeId=${rangeId}`;
    return this.http.get(url);
  }


  /**
   * 获取指定工作空间的规范中的语言
   * @param workspace
   * @returns
   */
  public getWorkSpaceLang(workspace: string): Observable<any> {
    // /actions/target/:workspace/defect/language
    let url: string = `/api/code-shooting/v1/actions/target/${workspace}/defect/language`;
    return this.http.get(url);
  }
  /**
   * 获取指定工作空间的规范中的指定语言与通用语言的缺陷分类信息
   * @param workspace
   * @returns
   */
  public getWorkSpaceLangDefect(workspace: string, language: string, templateVersion: string, needCode: boolean = false): Observable<any> {
    let url: string = `/api/code-shooting/v1/actions/target/${workspace}/defect/language`;
    return this.http.post(url, { language: language, needCode: needCode, templateVersion: templateVersion });
  }

  /**
   * 根据工作空间选择打靶规范
   * @param workspace
   * @returns
   */
  public getTemplateByWorkspace(workspace: string): Observable<any> {
    let url: string = `/api/code-shooting/v1/actions/target/${workspace}`;
    return this.http.get(url);
  }

  public getLangByworkspaceAndTemplate(workspace: string, template: string): Observable<any> {
    let url: string = `/api/code-shooting/v1/actions/target/${workspace}/${template}`;
    return this.http.get(url);
  }
}



import { Injectable } from '@angular/core';
import { Observable, Subject,zip } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { ScoreResponse } from './score';
const SPECIALCHARACTER = [';', ':', ',', '/', '?', '@', '&', '=', '+', '$', '#'];
@Injectable(
  {
    providedIn: 'root'
  }
)
export class QueryService {

    constructor(private http: HttpClient) {

    }

  getScore(rangeId, userId: string, language: string, targetId?: string): Observable<any> {
        let url = `/api/code-shooting/v1/results/range/${rangeId}/language/${language}?userId=${userId}&verbose=true`;
        if (targetId) {
          url += '&targetId=' + targetId;
        }
        return this.http.get<any>(url)
  }

  getTotalScore(rangeId, language: string): Observable<any> {
    let url = `/api/code-shooting/v1/answers/range/${rangeId}/language/${language}`;
    return this.http.get<any>(url)
  }

    /**
     * 对一些特殊字符（主要是用于分隔 URI 组件的标点符号）进行转义处理:位于url路径中的暂时不转义，位于?后边的要转义
     * @param character
     * @returns
     */
    transferCharacter(character: string): string {
      if (typeof(character) !== 'string') {
        return character;
      }
      let specialCharacter = '';
      for(let i = 0; i < character.length; i++) {
        if (SPECIALCHARACTER.indexOf(character[i]) >= 0) {
          specialCharacter += encodeURIComponent(character[i]);
        } else {
          specialCharacter += character[i];
        }
      }
      return specialCharacter;
    }
}

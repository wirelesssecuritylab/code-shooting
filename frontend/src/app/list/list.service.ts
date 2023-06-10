import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

const RANGEURL = '/api/code-shooting/v1/actions/range';

@Injectable({
  providedIn: 'root',
})
export class ListService {
  public type: string;
  public lan: string;
  constructor(private http: HttpClient) {}

  public getRangeList(userId: string, rangeId?: string): Observable<any> {
    let reqBody: any = {
      name: 'query',
      parameters: {
        id: rangeId,
        user: userId,
      },
    };
    console.log(reqBody);
    return this.http.post(RANGEURL, reqBody).pipe(
      map((res: any) => {
        if (
          res &&
          res.result == 'success' &&
          res.detail &&
          Array.isArray(res.detail)
        )
          console.log(res.detail);
        {
          return res.detail.map((range: any) => {
            range.targetId = range.id;
            range.start_at = new Date(range.startTime * 1000);
            range.end_at = new Date(range.endTime * 1000);
            range.languages = range.language;
            const languages = [];
            range.targets.forEach(function (value) {
              languages.push(value.language);
            });
            range.languages = [...new Set(languages)].sort();
            return range;
          });
        }
      })
    );
  }

  public queryLocalData(data: any[], requestBody: any) {
    if (requestBody == null) {
      return data;
    }
    let reqMap = new Map(Object.entries(requestBody));
    for (let [key, value] of reqMap) {
      if (value != null && value != '') {
        data = data.filter((item) => item[key].indexOf(value) >= 0);
      } else {
      }
    }
    return data;
  }

  public getFilterReq(parent: any) {
    let typeName = '';
    let lanName = '';
    if (parent?.key == 'type') {
      parent?.values.forEach(function (value) {
        if (value.isSelected == true) {
          typeName = value.key;
        }
      });
      this.type = typeName;
    } else if (parent?.key == 'languages') {
      parent?.values.forEach(function (value) {
        if (value.isSelected == true) {
          lanName = value.key;
        }
      });
      this.lan = lanName;
    }
    return {
      type: this.type,
      languages: this.lan,
    };
  }

  typeItem: any = {};
  public typeMap(data: any[], status: any[]) {
    return data.map((i) => {
      if (i.type === 'test') {
        //i.status = 0;
        this.typeItem = status[0];
      } else {
        //i.status = 1;
        i.disabled = false;
        this.typeItem = status[1];
      }
      //const statusItem = status[i.status];
      i.type = this.typeItem.text;
      return i;
    });
  }
}

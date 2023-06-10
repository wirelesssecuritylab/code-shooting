import {Injectable} from "@angular/core";

@Injectable()
export class CommonService {

  constructor() {
  }

  /**
   * 格式化时间戳为YYYY-MM-DD HH:MM:SS格式
   * @param stamp
   * @returns
   */
  public formatDateTime(stamp: number) {
    //stamp:是整数，否则要parseInt转换
    var time = new Date(stamp * 1000);
    var y = time.getFullYear();
    var m = time.getMonth() + 1;
    var d = time.getDate();
    var h = time.getHours();
    var mm = time.getMinutes();
    var s = time.getSeconds();
    return y + '-' + this.add0(m) + '-' + this.add0(d) + ' ' + this.add0(h) + ':' + this.add0(mm) + ':' + this.add0(s);
  }

  add0(number) {
    return number < 10 ? '0' + number : number;
  }
}

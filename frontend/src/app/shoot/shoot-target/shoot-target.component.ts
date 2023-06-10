import { Component, OnInit, OnDestroy, AfterViewInit, ViewChild } from '@angular/core';
import { PlxModal, FieldType, PlxBreadcrumbItem, IOption } from 'paletx';
import { ActivatedRoute } from '@angular/router';
import { languageAdapter } from "../../shared/class/row";
import { ManageService } from '../../admin/manage.service';
import { StoreService } from "../../shared/service/store.service";

@Component({
  selector: 'app-shoot-target',
  templateUrl: './shoot-target.component.html',
  styleUrls: ['./shoot-target.component.css']
})
export class ShootTargetComponent implements OnInit, OnDestroy {

  // 组件变量
  public breadModel: PlxBreadcrumbItem[] = [];
  public formSetting: any;
  public rangeId: string;
  public rangeName: string;
  public language: string;
  public shootType: string; // 打靶还是查看答卷
  public targetIds: string;  // 靶场下靶子的id（由于一个靶场可关联一个语言下的多个靶子，所以这里会有多个靶子id，需要传入时进行join）
  public fromPage: string;  // 打靶页面的跳转源页面
  public targetName: string;
  public rangeType: string;
  public isMyRange: boolean;
  constructor(private activatedRoute: ActivatedRoute,
    private manageService: ManageService,
    private storeService: StoreService) { }

  ngOnInit(): void {
    this.activatedRoute.queryParamMap.subscribe((paramMap) => {
      this.rangeId = paramMap.get('rangeId');
      this.rangeName = paramMap.get('rangeName');
      this.language = paramMap.get('language');
      this.shootType = paramMap.get('shootType');
      this.fromPage = paramMap.get('from');
      this.targetName = paramMap.get('targetName');
      this.rangeType = paramMap.get('rangeType');
      this.isMyRange = paramMap.get('isMyRange') == 'true';
      // 这里先实现单个靶子的打靶场景，默认取第一个靶子的
      // this.targetIds = paramMap.get('targetIds').split(',')[0] || '';
      this.targetIds = paramMap.get('targetId');
      this.storeService.setItem("currentTargetId", this.targetIds);
      this.storeService.setItem("currentDefectInfo", "");
      this.storeService.setItem("currentTargetLang", this.language);
      this.storeService.setItem("currentRangeId", this.rangeId);
      this.storeService.setItem("currentRangeName", this.rangeName);
      if (this.fromPage === 'test') {
        this.breadModel = [
          { label: '自由练习', routerLink: '../test', name: 'test' },
          { label: this.targetName, name: 'curTarget' }
        ];
      } else {
        this.breadModel = [
          { label: '靶场', routerLink: '/main/user/list', name: 'rangeList' },
          {
            label: this.rangeName + '-' + this.language, routerLink: '/main/user/shoot-range', name: 'rangeList',
            queryParams: { rangeId: this.rangeId, language: this.language, rangeName: this.rangeName, rangeType: this.rangeType }
          },
        ];
        if (this.isMyRange) {
          this.breadModel = [
            { label: '我的靶场', routerLink: '/main/user/personalcenter/myrange', name: 'myRangeList' },
            {
              label: this.rangeName + '-' + this.language, routerLink: '/main/user/shoot-range', name: 'myRangeList',
              queryParams: { rangeId: this.rangeId, language: this.language, rangeName: this.rangeName, rangeType: this.rangeType, isMyRange: true }
            },
          ];
        }
        if (this.shootType === 'view') {
          this.breadModel.push({ label: '查看答卷', name: 'link2' });
        } else {
          this.breadModel.push({ label: '打靶', name: 'link2' });
        }
      }
    });
    this.getDefectClass();
  }

  /**
   * 获取某个语言下的缺陷类型并保存到store
   */
  private getDefectClass() {
    let reqBody: any = {
      language: this.language,
      needCode: true
    };
    this.manageService.getDefectApi(this.targetIds, reqBody).subscribe(res => {
      if (res && res.result == 'success' && res.detail) {
        let defectDetail = res.detail;
        this.storeService.setItem('defectDetail', JSON.stringify(defectDetail));
      }
    }, error => { });
  }

  public ngOnDestroy() {
    this.storeService.removeItem('currentTargetId');
    this.storeService.removeItem('currentTargetLang');
    this.storeService.removeItem('currentRangeId');
    this.storeService.removeItem('defectDetail');
    this.storeService.removeItem('curRangeName');
  }

}

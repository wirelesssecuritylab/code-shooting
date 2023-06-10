import { ChangeDetectorRef, Component, OnInit, Input } from '@angular/core';
import { PlxBreadcrumbItem } from 'paletx';
import { ActivatedRoute } from '@angular/router';
import { ListService } from '../../list/list.service';
import { QueryService } from '../query/score/score.service';
import { Observable, zip } from 'rxjs';

const USER = 'user';
const TEST = 'test';

@Component({
  selector: 'app-shoot-range',
  templateUrl: './shoot-range.component.html',
  styleUrls: ['./shoot-range.component.css']
})
export class ShootRangeComponent implements OnInit {
  public userId: string;
  public breadModel: PlxBreadcrumbItem[] = [];
  public rangeName: string;
  public language: string;
  public rangeId: string;
  public rangeType: string;
  public isEnd: boolean = true;
  public targetId: string;
  public isMyRange: boolean = false;

  //四个用作靶场数据的参数
  hitNum: number;
  hitScore: number;
  totalNum: number;
  totalScore: number;
  hundredScore: number;

  constructor(private route: ActivatedRoute,
    private listService: ListService,
    private rangescoreService: QueryService,
    private cdr: ChangeDetectorRef, ) { }

  ngOnInit(): void {
    this.userId = localStorage.getItem(USER);
    //从本地存储取到userID
    this.route.queryParams.subscribe((params) => {
      //从路由取到以下字段
      this.rangeId = params.rangeId;
      this.language = params.language;
      this.rangeName = params.rangeName;
      this.rangeType = params.rangeType;
      if (params.isMyRange == 'true') {
        this.isMyRange = true;
      }
    });
    this.breadModel = [
      { label: '靶场', routerLink: '/main/user/list', name: 'rangeList' },
      { label: this.rangeName + '-' + this.language, name: 'link2' },
    ];
    if (this.isMyRange) {
      this.breadModel[0] = { label: '我的靶场', routerLink: '/main/user/personalcenter/myrange', name: 'myRangeList' };
    }

    this.getRangeInfo();
  }

  /**
   * 获取靶场详情，以用来判断打靶是否结束
   */
  public getRangeInfo() {
    let slbInfo$: Observable<any>;
    let pvcInfo$: Observable<any>;
    slbInfo$ = this.rangescoreService.getScore(this.rangeId, this.userId, this.language);
    pvcInfo$ = this.rangescoreService.getTotalScore(this.rangeId, this.language);
    //使用接口取到靶场数据
    // this.rangescoreService.getScore(this.rangeId, this.userId, this.language)
    //   .subscribe(res => {
    //     if (res && Array.isArray(res) && res.length == 1) {
    //       let result = res[0];
    //       if (result.rangeScore) {
    //         this.hitNum = result.rangeScore.hitNum ? result.rangeScore.hitNum : 0;
    //         this.hitScore = result.rangeScore.hitScore ? result.rangeScore.hitScore : 0;
    //         console.log(this.hitScore)
    //       }
    //     }
    //     this.cdr.detectChanges();
    //   },
    //   )

    //   this.rangescoreService.getTotalScore(this.rangeId, this.language)
    //   .subscribe(res => {
    //     if (res) {
    //         this.totalNum = res.detail.totalNum ? res.detail.totalNum : 0;
    //       this.totalScore = res.detail.totalScore ? res.detail.totalScore : 0;
    //       console.log(this.totalScore)
    //     }

    //   })
    zip(slbInfo$, pvcInfo$).subscribe((result: any) => {
      if (result[0][0].rangeScore == undefined) {
        this.hitNum = 0;
        this.hitScore = 0;
        this.totalNum = 0;
        this.totalScore = 0;
      } else {
        this.hitNum = result[0][0].rangeScore.hitNum ? result[0][0].rangeScore.hitNum : 0;
        this.hitScore = result[0][0].rangeScore.hitScore ? result[0][0].rangeScore.hitScore : 0;
        this.totalNum = result[1].detail.totalNum ? result[1].detail.totalNum : 0;
        this.totalScore = result[1].detail.totalScore ? result[1].detail.totalScore : 0;
      }

      this.hundredScore = Math.trunc(this.hitScore / this.totalScore * 100);
    }
    )



    // );
    //判断打靶比赛是否结束，传出参数isEnd
    this.listService.getRangeList(this.userId, this.rangeId).subscribe((res) => {
      if (res && Array.isArray(res) && res.length == 1) {
        let endTime = res[0].endTime * 1000;
        let nowTime = new Date(Date.parse(new Date().toString()));
        if (nowTime.getTime() <= endTime) {
          this.isEnd = false;
        } else {
          this.isEnd = true;
        }
      }
    });
  }

  changeTargetId(id: any) {
    this.targetId = id;
  }
  changeTargetId1(id: any) {
    this.targetId = id;
  }


}

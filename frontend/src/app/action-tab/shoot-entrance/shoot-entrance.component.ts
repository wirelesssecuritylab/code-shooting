import {
  Component,
  Input,
  OnInit,
  TemplateRef,
  ViewChild,
} from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { DefectService } from '../../shoot/services/defect-service';
import { PlxModal } from 'paletx';
import { ManageService } from '../../admin/manage.service';

@Component({
  selector: 'app-shoot-entrance',
  templateUrl: './shoot-entrance.component.html',
  styleUrls: ['./shoot-entrance.component.css'],
})
export class ShootEntranceComponent implements OnInit {
  @Input() public targets: any[];
  @Input() public userId: string;
  @Input() public disabled: boolean;
  @Input() public type: string;
  @Input() public language: string;
  @Input() public rangeName: string;
  @Input() public rangeId: string;
  @Input() public records: any[]; // 打靶记录（草稿）
  @Input() public targetAnswers: any[]; // 靶标答案
  @Input() public rangeType: string;
  @Input() public isSubmitAnswer: boolean;
  @Input() public targetId: string;
  @Input() public isMyRange: boolean;
  btnLabel: string;
  btnWarnTip: string;
  targetIds: string;
  modal: any;

  @ViewChild('confirmInfo', { static: true })
  public confirmInfo: TemplateRef<any>;

  constructor(
    private router: Router,
    private activatedRoute: ActivatedRoute,
    private defectService: DefectService,
    private manageService: ManageService,
    private modalService: PlxModal
  ) { }

  ngOnInit(): void {
    this.initBtn(this.type);
  }

  public initBtn(type: string) {
    if (type == 'startShoot') {
      this.btnLabel = '开始打靶';
      this.btnWarnTip = '打靶比赛已结束，不允许开始打靶';
    } else if (type == 'continueShoot') {
      this.btnLabel = '继续打靶';
      this.btnWarnTip = '打靶比赛已结束，不允许继续打靶';
    } else if (type == 'restartShoot') {
      this.btnLabel = '重新打靶';
      this.btnWarnTip = '打靶比赛已结束，不允许重新打靶';
    } else if (type == 'view') {
      this.btnLabel = '查看答卷';
    }
  }

  navigator() {
    if (this.type === 'restartShoot' || this.type === 'continueShoot') {
      this.modal = this.modalService.open(this.confirmInfo, { size: 'xs' });
    } else {
      // 如果是查询答卷，则把defect服务中的targets初始化为从接口读取的打靶记录
      if (this.type === 'view') {
        let reqBody: any = {
          name: 'get',
          parameters: {
            Id: this.rangeId
          }
        };
        this.manageService.manageRangeApi(reqBody).subscribe(res => {
          if ((new Date().getTime() / 1000) >= (new Date(res.detail.EndTime).getTime() / 1000)) {
            this.defectService.targets = this.records;

          } else {
            this.defectService.targets = [];
          }
          this.defectService.targetAnswers = this.targetAnswers;
        });
      } else {
        this.defectService.targets = [];
        this.defectService.targetAnswers = [];
      }
      this.gotoShoot();
    }
  }

  confirmNavigator(isLoadRecord: boolean) {
    this.modal.close();
    if (isLoadRecord) {
      this.defectService.targets = this.records;
    } else {
      this.defectService.targets = [];
    }
    this.gotoShoot();
  }

  /**
     * 点击打靶进入打靶页面
     */
  gotoShoot() {
    // const targetIds = this.targets.map((target) => target.targetId);
    // this.targetIds = targetIds.join(',');
    let routeStr: string = '../shoot';
    this.router.navigate([routeStr], {
      queryParams: {
        rangeName: this.rangeName,
        rangeId: this.rangeId,
        language: this.language,
        rangeType: this.rangeType,
        targetId: this.targetId,
        shootType: this.type,
        isMyRange: this.isMyRange,
      },
      relativeTo: this.activatedRoute,
    });
  }
}


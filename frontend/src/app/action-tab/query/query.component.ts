import { Target } from './score/score.d';
import { ChangeDetectorRef, Component, Input, OnInit } from '@angular/core';
import { FieldType, PlxMessage } from "paletx";
import { ScoreResponse } from './score/score'
import { QueryService } from "./score/score.service";
import { ManageService } from '../../admin/manage.service';
@Component({
  selector: 'app-query',
  templateUrl: './query.component.html',
  styleUrls: ['./query.component.css']
})
export class QueryComponent implements OnInit {
  formSetting: any;
  lansSelector: any[];
  shotTargets: Target[];
  result: ScoreResponse;
  hitNum: number;
  hitScore: number;
  totalNum: number;
  totalScore: number;
  showResult: boolean;
  hoverStr: string;
  constructor(
    private queryService: QueryService,
    private manageService: ManageService,
    private cdr: ChangeDetectorRef,
    private plxMessageService: PlxMessage
  ) { }
  @Input() public disabled: boolean;
  @Input() public rangeId: string;
  @Input() public language: string;
  @Input() public userId: string;
  @Input() public targetURL: string;
  @Input() public submit: string;
  @Input() public rangeType: string;
  @Input() public targetId: string;
  ngOnInit(): void {
    this.setForm()
    this.hoverStr = '查询答卷的得分及详情，以便持续改进';
    // 自由练习提交答卷后自动显示成绩
    if (this.submit == 'true' && this.rangeType == 'test') {
      this.getResult();
    }
  }

  /**
     * 查询个人打靶成绩（目前先按支持一个靶子的解析）
     */
  getResult() {
    let reqBody: any = {
      name: 'get',
      parameters: {
        Id: this.rangeId
      }
    };
    this.manageService.manageRangeApi(reqBody).subscribe(res => {
      if ((new Date().getTime() / 1000) >= (new Date(res.detail.EndTime).getTime() / 1000)) {
        this.queryService.getScore(this.rangeId, this.userId, this.language, this.targetId).subscribe(
          res => {
            if (res && Array.isArray(res) && res.length == 1) {
              this.result = res[0];
              if (this.result.targets && Array.isArray(this.result.targets) && this.result.targets.length >= 0) {
                let target: any = this.result.targets[0];
                this.shotTargets = target.detail;
                this.hitNum = target.hitNum;
                this.hitScore = target.hitScore;
                this.totalNum = target.totalNum;
                this.totalScore = target.totalScore;
                this.showResult = true;
              }
            }
            this.cdr.detectChanges();
          },
          error => {
            console.log(error)
            if (error) {
              // this.plxMessageService.show('error',
              //   {
              //     title: '查看成绩失败，请先提交答卷。',
              //     isLightweight: true
              //   });
              this.plxMessageService.error('查看成绩失败，请先完成打靶。', '');
            }
          }
        );
      }
    });
  }

  setForm() {
    this.formSetting = {
      isShowHeader: false,
      header: '查看成绩',
      isGroup: true,
      hideGroup: true,
      advandedFlag: true,
      boldOnValueChange: false,
      buttons: [
        {
          type: 'submit',
          label: '确定',
          class: 'plx-btn plx-btn-primary',
          hidden: false,
          disabled: false,
          callback: (values, $event, controls) => {

          }
        },
        {
          type: 'cancel',
          label: '取消',
          class: 'plx-btn',
          hidden: false,
          disabled: false,
          callback: (values, $event, controls) => {

          }
        }
      ],
      fieldSet: [
        {
          group: '基本信息',
          fields: [
            {
              name: 'language',
              label: '语言',
              type: FieldType.SELECTOR,
              required: true,
              multiple: false,
              valueSet: this.lansSelector,
            },]
        }
      ],
    }
  }

}


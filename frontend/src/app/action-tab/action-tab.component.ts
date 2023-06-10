import { Component, Input, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { AnchorItem, PlxBreadcrumbItem, PlxMessage } from 'paletx';
import { ListService } from '../list/list.service';
import { ActivatedRoute, Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { ManageService } from '../admin/manage.service';
import { QueryService } from './query/score/score.service';
import { DefectService } from '../shoot/services/defect-service';

const TEST = 'test';
const USER = 'user';

@Component({
  selector: 'app-action-tab',
  templateUrl: './action-tab.component.html',
  styleUrls: ['./action-tab.component.css'],
})
export class ActionTabComponent implements OnInit, OnDestroy {
  rangeId: string;
  userId: string;
  language: string;
  rangeInfo: any = {};
  valid: boolean = false;
  errorTip: string;
  errTitle: string;
  errImg: string;
  targets: any[];
  shootType = 'startShoot';
  public nowTime: Date;
  public rangeName: string;

  public targetRecords: any[] = []; // 打靶记录
  public targetAnswers: any[] = []; // 靶标
  public targetDrafts: any[] = []; // 打靶草稿
  public haveSubmitAnswer: boolean = false;
  public submit: string;
  public rangeType: string;
  shootTypeMap = {
    startShoot: '开始打靶',
    continueShoot: '继续打靶',
    restartShoot: '重新打靶',
  };
  @Input() targetId: string;
  @Input() isMyRange: boolean;

  @ViewChild('shootEntrance') shootEntrance;

  constructor(
    private listService: ListService,
    private http: HttpClient,
    public route: ActivatedRoute,
    private manageService: ManageService,
    private plxMessageService: PlxMessage,
    private queryService: QueryService,
    private defectService: DefectService,
    private router: Router
  ) { }

  ngOnInit(): void {
    this.userId = localStorage.getItem(USER);
    this.route.queryParams.subscribe((params) => {
      this.rangeId = params.rangeId;
      this.submit = params.submit;
      this.rangeType = params.rangeType;
      this.language = decodeURIComponent(params.language);
      this.getTargetsInfo();
    });
    console.log(this.rangeType)
  }

  initTargetInfo() {
    if (!this.targetId) {
      return;
    }
    this.getShootDraft(this.targetId);
    this.getShootRecords(this.targetId);
    this.getTargetsAnswers(this.targetId, this.rangeId);
    this.judgeIsSubmitAnswer(this.targetId);
  }

  getTargetsInfo() {
    this.nowTime = new Date(Date.parse(new Date().toString()));
    this.listService
      .getRangeList(this.userId, this.rangeId)
      .subscribe((res) => {
        //取到这个用户ID能看到的靶场ID（相当于过滤一遍数据），再
        if (res && Array.isArray(res) && res.length > 0) {
          this.rangeInfo = {
            start_at: res[0].startTime * 1000,
            end_at: res[0].endTime * 1000,
            type: res[0].type,
            name: res[0].name,
            targets: res[0].targets,
            //从靶场ID取到对应的靶子ID，找到里面要用的数据
          };
          this.targets = this.rangeInfo.targets.filter(
            (target) =>
              target.language.toUpperCase() == this.language.toUpperCase()
          );
          this.valid = true;
          this.initTargetInfo();
          // if (this.targets.length > 0) {
          //   this.language = this.targets[0].language;
          //   this.valid = true;
          //   let targetId = this.targets[0].targetId;
          //   // let targetId = this.
          //   this.getShootDraft(targetId);
          //   this.getShootRecords(targetId);
          //   this.judgeIsSubmitAnswer(targetId);
          // } else {
          //   this.errorTip = `您访问的靶场不支持${this.language}语言`;
          // }
        } else {
          this.errorTip = '您访问的靶场已删除或不存在';
        }
        this.errTitle = '访问失败';
        this.errImg =
          'data:image/png;base64,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';
      });
  }

  showShootTab(): boolean {
    return !(
      this.rangeInfo.type == TEST ||
      (this.nowTime.getTime() >= this.rangeInfo.start_at &&
        this.nowTime.getTime() <= this.rangeInfo.end_at)
    );
  }

  showViewTab(): boolean {
    return (
      !this.haveSubmitAnswer ||
      (this.haveSubmitAnswer &&
        this.rangeInfo.type !== TEST &&
        this.nowTime.getTime() <= this.rangeInfo.end_at)
    );
  }

  showQueryTab(): boolean {
    return (
      !(
        this.rangeInfo.type == TEST ||
        this.nowTime.getTime() >= this.rangeInfo.end_at
      ) ||
      (this.rangeInfo.type == TEST && !this.haveSubmitAnswer)
    );
  }

  public items: Array<AnchorItem> = [
    {
      anchorId: 'anchor1',
      containerId: 'containerId1',
      anchorName: '开始打靶',
    },
    {
      anchorId: 'anchor2',
      containerId: 'containerId2',
      anchorName: '查看成绩',
    },
    {
      anchorId: 'anchor3',
      containerId: 'containerId3',
      anchorName: '查看答卷',
    },
  ];

  ngOnDestroy(): void { }

  /**
   * 获取某个靶子的打靶草稿
   * 继续打靶时所用
   * @param targetId 靶子ID
   * @returns
   */
  public getShootDraft(targetId: string) {
    if (!this.userId || !targetId) {
      return;
    }
    this.manageService
      .getShootDraftApi(this.userId, targetId, this.rangeId)
      .subscribe(
        (res) => {
          if (
            res &&
            res.targets &&
            Array.isArray(res.targets) &&
            res.targets.length > 0
          ) {
            this.targetDrafts = res.targets.map((item) => {
              let record = {
                DefectClass: item.defectClass,
                DefectSubClass: item.defectSubClass,
                DefectDescribe: item.defectDescribe,
                FileName: item.fileName,
                Remark: item.remark,
                StartLineNum: item.startLineNum,
                EndLineNum: item.endLineNum,
                TargetId: res.targetid,
                StartColNum: item.startColNum,
                EndColNum: item.endColNum,
              };
              return record;
            });
            //继续打靶
            this.setShootEntranceLabel('continueShoot');
          } else {
            this.targetDrafts = [];
            this.getShootScore(targetId);
          }
        },
        (error) => {
          this.targetDrafts = [];
          this.plxMessageService.error(
            '获取打靶草稿失败',
            error.message || error.error.message
          );
        }
      );
  }

  getShootScore(targetId: string) {
    this.queryService
      .getScore(this.rangeId, this.userId, this.language, targetId)
      .subscribe(
        (res) => {
          if (res && Array.isArray(res) && res.length > 0) {
            // 重新打靶
            this.setShootEntranceLabel('restartShoot');
          } else {
            // 开始打靶
            this.setShootEntranceLabel('startShoot');
          }
        },
        (error) => {
          // 开始打靶
          this.setShootEntranceLabel('startShoot');
        }
      );
  }

  setShootEntranceLabel(type: string) {
    this.shootType = type;
    this.items[0].anchorName = this.shootTypeMap[this.shootType];
    this.shootEntrance.initBtn(this.shootType);
  }

  /**
   * 通过查询打靶成绩判断是否提交过答卷
   * @param targetId
   */
  public judgeIsSubmitAnswer(targetId: string) {
    this.queryService
      .getScore(this.rangeId, this.userId, this.language, targetId)
      .subscribe(
        (res) => {
          if (res && Array.isArray(res) && res.length > 0) {
            // 有成绩代表已经提交过答卷
            this.haveSubmitAnswer = true;
            this.defectService.targetScore = res;
          } else {
            // 未提交
            this.haveSubmitAnswer = false;
          }
        },
        (error) => {
          // 未提交
          this.haveSubmitAnswer = false;
        }
      );
  }

  /**
   * 获取靶标, 当查看答卷时所用
   * @param targetId 
   * @param rangeId
   * @returns 
   */
  public getTargetsAnswers(targetId: string, rangeId: string) {
    this.manageService.getTargetsAnswer(targetId, rangeId).subscribe(res => {
      if (res && res.targetid && Array.isArray(res.answers) && res.answers.length > 0) {
        this.targetAnswers = res.answers.map((item) => {
          let record = {
            DefectClass: item.defectClass,
            DefectSubClass: item.defectSubClass,
            DefectDescribe: item.defectDescribe,
            FileName: item.fileName,
            Remark: "",
            StartLineNum: item.startLineNum,
            EndLineNum: item.endLineNum,
            TargetId: res.targetid,
            StartColNum: 0,
            EndColNum: item.defectDescribe.length,
          };
          return record;
        });
      } else {
        this.targetAnswers = [];
      }
    });
  }

  /**
   * 获取某个靶子的打靶记录
   * 当查看答卷时所用
   * @param targetId 靶子ID
   * @returns
   */
  public getShootRecords(targetId: string) {
    if (!this.userId || !targetId) {
      return;
    }
    this.manageService
      .getShootRecordsApi(this.userId, targetId, this.rangeId)
      .subscribe(
        (res) => {
          if (
            res &&
            res.targets &&
            Array.isArray(res.targets) &&
            res.targets.length > 0
          ) {
            this.targetRecords = res.targets.map((item) => {
              let record = {
                DefectClass: item.defectClass,
                DefectSubClass: item.defectSubClass,
                DefectDescribe: item.defectDescribe,
                FileName: item.fileName,
                Remark: item.remark,
                StartLineNum: item.startLineNum,
                EndLineNum: item.endLineNum,
                TargetId: res.targetid,
                StartColNum: item.startColNum,
                EndColNum: item.endColNum,
              };
              return record;
            });
          } else {
            this.targetRecords = [];
          }
        },
        (error) => {
          this.targetRecords = [];
          this.plxMessageService.error(
            '获取打靶记录失败',
            error.message || error.error.message
          );
        }
      );
  }
}

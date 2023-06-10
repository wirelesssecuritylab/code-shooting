import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { PlxBreadcrumbItem, PlxMessage} from 'paletx';
import { ScoreResponse, Target } from '../../action-tab/query/score/score';
import {QueryService} from "../../action-tab/query/score/score.service";

const TESTRANGEID = 0;
@Component({
  selector: 'app-test-score',
  templateUrl: './test-score.component.html',
  styleUrls: ['./test-score.component.css']
})
export class TestScoreComponent implements OnInit {
  public targetId: string;
  public language: string;
  public breadModel: PlxBreadcrumbItem[] = [];
  shotTargets: Target[];
  result: ScoreResponse;
  hitNum: number;
  hitScore: number;
  totalNum: number;
  totalScore: number;
  showResult: boolean = false;
  public userId: string;
  constructor(private activatedRoute: ActivatedRoute, private queryService: QueryService,
              private plxMessageService: PlxMessage) { }

  ngOnInit(): void {
    this.userId = localStorage.getItem('user');
    this.breadModel = [
      {label: '自由练习', routerLink: '../test', name: 'test'},
      {label: '查看成绩', name:'curTarget'}
    ];
    this.activatedRoute.queryParamMap.subscribe((paramMap) => {
      this.targetId = paramMap.get('targetId');
      this.language = paramMap.get('language');
      this.getTestScore();
    });
  }

  /**
   * 查看个人打靶成绩
   */
  public getTestScore() {
    this.queryService.getScore(TESTRANGEID, this.userId, this.language, this.targetId).subscribe(
      res => {
        console.log(res)
        if(res && Array.isArray(res) && res.length == 1) {
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
      },
      error => {
        this.showResult = false;
        console.log(error)
        // if (error) {
        //   this.plxMessageService.error('成绩查看失败', error.message || error.error.message);
        // }
      }
    );
  }

}

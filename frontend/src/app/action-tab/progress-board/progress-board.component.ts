import { Target } from './../query/score/score.d';
import { Component, Input, OnInit, Output, EventEmitter, SimpleChanges } from '@angular/core';
import { ListService } from 'src/app/list/list.service';
import { QueryService } from ".././query/score/score.service";

const USER = 'user';
@Component({
  selector: 'app-progress-board',
  templateUrl: './progress-board.component.html',
  styleUrls: ['./progress-board.component.css'],
})
export class ProgressBoardComponent implements OnInit {
  constructor(
    private listService: ListService,
    private queryService: QueryService,
  ) {}
  @Input() public language: string;
  @Input() public rangeId: string;
  @Input() public targetId: string;
  @Output() public targetChanged = new EventEmitter<any>();
  public userId: string;
  public activeId;
  public progressBoardButtons = [];
  public targets: any[];
  public current: number = 1;
  public answers: any[] = [];
  public currentTargetId: string;

  ngOnInit(): void {
    this.userId = localStorage.getItem(USER);
    this.getTargetsInfo();
    this.getResult();
  }

  ngOnChanges(change: SimpleChanges) {
    if (this.targets) {
      this.current = this.targets.findIndex(item => {
        return item.targetId == this.targetId
    }) + 1
    }
  }

  getTargetsInfo() {
    this.listService
      .getRangeList(this.userId, this.rangeId)
      .subscribe((res) => {
        if (res && Array.isArray(res)) {
          this.targets = res[0].targets.filter(
            (language) => {
              return language.language.toUpperCase() === this.language.toUpperCase()
            });
          for (let i = 0; i < this.targets.length; i++) {
            const progressBoardButton = {};
            progressBoardButton['id'] = this.targets[i].targetId; //
            this.progressBoardButtons.push(progressBoardButton);

          }
        }
      });
  }


  gotoTarget(targetId: string, index: number) {
    this.currentTargetId = targetId;
    this.targetChanged.emit(targetId);
    this.current = index + 1;
  }

  getResult() {
    this.queryService
      .getScore(this.rangeId, this.userId, this.language)
      .subscribe(
        (res) => {
          this.answers = res[0].targets.map((item) => {
            return item.targetId;
          });
        },
      );
  }
}



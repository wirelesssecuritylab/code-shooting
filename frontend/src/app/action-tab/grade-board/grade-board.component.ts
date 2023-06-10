import { Component, OnInit, Input, ChangeDetectorRef } from '@angular/core';

const USER = 'user';
const TEST = '练习';

@Component({
  selector: 'app-grade-board',
  templateUrl: './grade-board.component.html',
  styleUrls: ['./grade-board.component.css'],
})
export class GradeBoardComponent implements OnInit {
  @Input() public language: string;
  @Input() public rangeId: string;
  @Input() public rangeType: string;
  @Input() public isEnd: boolean;
  public userId: string;
  @Input() public hitNum: number = 0;
  @Input() public hitScore: number = 0;
  @Input() public totalNum: number = 0;
  @Input() public totalScore: number = 0;
  @Input() public hundredScore: number = 0;

  constructor(
    private cdr: ChangeDetectorRef
  ) {}

  ngOnInit(): void {
    this.userId = localStorage.getItem(USER);
  }
}

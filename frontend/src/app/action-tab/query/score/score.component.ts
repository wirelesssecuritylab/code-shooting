import {Component, OnInit, Input, ChangeDetectorRef} from '@angular/core';
import { Target } from './score'
import {PlxTableService} from "paletx";

@Component({
  selector: 'result-render',
  template: `
    <div>
        <span  [hidden]= "score <= 0">{{value}}</span>
        <span [hidden]= "score > 0" style="color:#FF0000">{{value}}</span>
    </div>
`
})
export class ResultRenderComponent implements OnInit{
  rowData: any;
  score: number;
  value: string;
  constructor(private plxTableService: PlxTableService) {
    this.rowData = this.plxTableService.getPositionInfo().currentRowData;
    this.value  = this.rowData[this.plxTableService.getPositionInfo().currentColumn.key];
  }
  ngOnInit() {
    this.score = this.rowData.score;
  }
}

@Component({
  selector: 'app-score',
  templateUrl: './score.component.html',
  styleUrls: ['./score.component.css']
})
export class ScoreComponent implements OnInit {

  pageSizeSelections = [20, 50, 100, 200];
  scroll={y: '325px'}
  columns = [
    {
      key: 'fileName',
      title: '文件名',
      sort: 'asc',
      filter: true,
      contentType: 'component',
      component: ResultRenderComponent,
    },

    {
      key: 'startLineNum',
      title: '起始行号',
      show: true,
      filter: true,
      sort: 'asc',
      contentType: 'component',
      component: ResultRenderComponent,
    },
    {
      key: 'endLineNum',
      title: '结束行号',
      show: true,
      filter: true,
      sort: 'asc',
      contentType: 'component',
      component: ResultRenderComponent,
    },
    {
      key: 'defectClass',
      title: '缺陷大类',
      show: true,
      filter: true,
      sort: 'asc',
      contentType: 'component',
      component: ResultRenderComponent,
    },
    {
      key: 'defectSubClass',
      title: '缺陷小类',
      show: true,
      filter: true,
      sort: 'asc',
      contentType: 'component',
      component: ResultRenderComponent,
    },
    {
      key: 'defectDescribe',
      title: '缺陷细项',
      show: true,
      filter: true,
      sort: 'asc',
      contentType: 'component',
      component: ResultRenderComponent,
    },
    {
      key: 'remark',
      title: '缺陷备注',
      show: true,
      filter: true,
      sort: 'asc',
      contentType: 'component',
      component: ResultRenderComponent,
    },
    {
      key: 'score',
      title: '得分',
      filter: true,
      show: true,
      sort: 'asc',
      contentType: 'component',
      component: ResultRenderComponent,
    },
  ];

  @Input() data: Target[] = [];
  @Input() showResult: boolean;
  @Input() hitNum: number;
  @Input() hitScore: number;
  @Input() totalNum: number;
  @Input() totalScore: number;

  constructor(private cdr: ChangeDetectorRef) { }

  ngOnInit(): void {
    this.cdr.detectChanges();
  }

}

import {Component, Input, OnInit, ViewChild} from "@angular/core";
import {ColumnDefine, DialogBase, JigsawDialog, TableData,} from "@rdkmaster/jigsaw";
import {DefectService} from "../../services/defect-service";
import {Router, ActivatedRoute} from "@angular/router";
import { ManageService } from '../../../admin/manage.service';
import {QueryService} from '../../../action-tab/query/score/score.service';
import { PlxMessage } from 'paletx';
import {StoreService} from "../../../shared/service/store.service";
import {ShootingResult, TargetInfo, TargetResult} from "../../misc/target-types";
@Component({
    templateUrl: "./show-defects.dialog.html",
    styleUrls: ['./show-defects.dialog.scss']
})
export class ShowDefectsDialog extends DialogBase implements OnInit {
    @Input()
    public initData: any;

    @ViewChild(JigsawDialog)
    public dialog: JigsawDialog;

    public _$defectTableData: TableData;
    public curRangeId: string;
    public curTargetLang: string;
    public curRangeName: string;
    public fromPage: string;
    public rangeType: string;
    public isMyRange: boolean;
    public targetId:string;
    public _$columnDefines: ColumnDefine[] = [
        {
            target: ['StartLineNum', 'EndLineNum'],
            width: 70
        }
    ];

    constructor(private _defectService: DefectService,
                private router: Router, private activatedRoute: ActivatedRoute,
                private manageService: ManageService,
                private plxMessageService: PlxMessage,
                private queryService: QueryService,
                private storeService: StoreService) {
        super();
    }

    ngOnInit(): void {
        this.curRangeId = this.storeService.getItem('currentRangeId');
        this.curTargetLang = this.storeService.getItem('currentTargetLang');
        this.curRangeName = this.storeService.getItem('currentRangeName');
        const data = this._defectService.shootingResult.targets.map(target =>
            [target.fileName, target.startLineNum, target.endLineNum, target.defectClass, target.defectSubClass, target.defectDescribe, target.remark]);
        this._$defectTableData = new TableData(
            data,
            ['FileName', 'StartLineNum', 'EndLineNum', 'DefectClass', 'DefectSubClass', 'DefectDescribe','Remark'],
            ['文件名', '开始行号', '结束行数', '缺陷大类', '缺陷小类', '缺陷描述', '缺陷备注']
        );
        this.activatedRoute.queryParamMap.subscribe((paramMap) => {
          this.fromPage = paramMap.get('from');
          this.rangeType = paramMap.get('rangeType');
          this.isMyRange = paramMap.get('isMyRange') == 'true';
          this.targetId = paramMap.get('targetId');
        });
    }

    /**
     * 确认提交代码答卷
     * @returns
     */
    confirmSubmit() {
      // console.log("确认要提交的评审信息：", this._defectService.shootingResult);
      if (this.fromPage === 'test') {
        this.curRangeId = '0';
      } else {
        if (!this.curRangeId || !this.curTargetLang) {
          return;
        }
      }
      let lang = this.queryService.transferCharacter(this.curTargetLang);

      this.manageService.submitTargetAnswer(this.curRangeId, lang, this._defectService.shootingResult).subscribe(res => {
        this.plxMessageService.success('提交成功', '');
        this.dispose();
       
        if (this.fromPage === 'test') {
          this.router.navigate(['../test'], {
            queryParams:{'targetId':this.targetId},  
            relativeTo: this.activatedRoute
          
          });
        } else {
          this.backtoRangeDetail();
        }
      }, error => {
        this.plxMessageService.error('提交失败', error.error.message);
      });
    }

    /**
     * 提交成功返回靶场详情页面
     */
    backtoRangeDetail() {
      const queryParams = {
        'rangeId': this.curRangeId,
        'language': this.curTargetLang,
        'rangeName': this.curRangeName,
        'submit': 'true',
        'rangeType': this.rangeType,
      };
      if(this.isMyRange) {
        queryParams['isMyRange'] = this.isMyRange;
      }
      let routeStr: string = '../shoot-range';
      this.router.navigate([routeStr], {
        queryParams: queryParams,
        relativeTo: this.activatedRoute
      });
    }
}

import {Component, Input, OnInit, ViewChild} from "@angular/core";
import {DialogBase, JigsawDialog,} from "@rdkmaster/jigsaw";
import * as defect from './defect.json';
import {DefectClass, DefectSubClass, InitData} from "../../../misc/types";
import {StoreService} from "../../../../shared/service/store.service";

@Component({
    templateUrl: "./defect-select.dialog.html",
    styleUrls: ['./defect-select.dialog.scss']
})
export class DefectSelectDialog extends DialogBase implements OnInit {
    @Input()
    public initData: InitData;

    @ViewChild(JigsawDialog)
    public dialog: JigsawDialog;

    public _$currentDefectClass: any;
    public _$defectClasses: any[];
    public _$currentDefectSubClass: any;
    public _$defectSubClasses: any[];
    public _$currentDefectDescribe: string;
    public _$defectDescribes: any[];
    public _$title = '选择缺陷类型';
    public _$optionCount: number  = 12;
    public _$startLineNumber: string;
    public _$endLineNumber: string;
    public _$remark: string;
    public _$defectClassObject: any;
    public curTargetId: string; // 当前正在打靶的靶子id
    public _$quickSearch: string;
    public _$defectAllDescribe: any[];
    public _$defectAllDescribeData: any[];

    constructor(private storeService: StoreService) {
      super();
    }

    ngOnInit(): void {
      let defectDetail = this.storeService.getItem("defectDetail");
      this.curTargetId = this.storeService.getItem("currentTargetId");
      if (defectDetail) {
        this._$defectClassObject = JSON.parse(defectDetail);
      } else {
        this._$defectClassObject = {};
      }
      this._$startLineNumber = this.initData.range.startLineNumber.toString();
      this._$endLineNumber = this.initData.range.endLineNumber.toString();
      this._$defectClasses = Object.keys(this._$defectClassObject);
      this._$defectAllDescribeData = this._$combineAllDescribe();
      this._$defectAllDescribe = this._$defectAllDescribeData.map(describe => describe.description);
      // 编辑态
      if (this.initData?.defectTarget) {
        this._$title = '修改缺陷类型';
        this._$currentDefectClass = this.initData.defectTarget.DefectClass;
        this._$defectSubClasses = Object.keys(this._$defectClassObject[this._$currentDefectClass])
        this._$currentDefectSubClass = this.initData.defectTarget.DefectSubClass;
        let defectDescribes = this._$defectClassObject[this._$currentDefectClass][this._$currentDefectSubClass];
        this._$defectDescribes = defectDescribes.map(item => {
          return item.description;
        });
        this._$currentDefectDescribe = this.initData.defectTarget.DefectDescribe;
        this._$remark = this.initData.defectTarget.Remark;
      } else {
        const currentDefectInfoStr = this.storeService.getItem("currentDefectInfo");
        if(currentDefectInfoStr && currentDefectInfoStr != '') {
          const curDefectInfo = JSON.parse(currentDefectInfoStr);
          this._$currentDefectClass = curDefectInfo.defectClass;
          this._$defectSubClasses = Object.keys(this._$defectClassObject[this._$currentDefectClass])
          this._$currentDefectSubClass = curDefectInfo.defectSubClass;
          let defectDescribes = this._$defectClassObject[this._$currentDefectClass][this._$currentDefectSubClass];
          this._$defectDescribes = defectDescribes.map(item => {
            return item.description;
          });
          this._$currentDefectDescribe = curDefectInfo.defectDescribe;
          this._$remark = curDefectInfo.remark;
        }else{
          this._$currentDefectClass = this._$defectClasses[0];
          this._$defectClassChange();
        }
      }
    }

    private _$combineAllDescribe() {
      const allDescribe = [];
      Object.keys(this._$defectClassObject).forEach(defClass => {
        Object.keys(this._$defectClassObject[defClass]).forEach(defSubClass => {
          this._$defectClassObject[defClass][defSubClass].forEach(describe => {
            allDescribe.push(Object.assign({defClass: defClass, subDefClass: defSubClass},describe));
          });
        });
      });
      return allDescribe;
    }

    /**
     * 提交当前评审意见
     */
    public _$confirm(): void {
      const cacheInfo = {
        defectClass: this._$currentDefectClass,
        defectSubClass: this._$currentDefectSubClass,
        defectDescribe: this._$currentDefectDescribe,
        remark: this._$remark
      };
      this.storeService.setItem('currentDefectInfo', JSON.stringify(cacheInfo));
      this.answer.emit({
        TargetId: this.curTargetId,
        FileName: this.initData.fileName,
        StartLineNum: this.initData.range.startLineNumber,
        StartColNum: this.initData.range.startColumn,
        EndLineNum: this.initData.range.endLineNumber,
        EndColNum: this.initData.range.endColumn,
        DefectClass: this._$currentDefectClass,
        DefectSubClass: this._$currentDefectSubClass,
        DefectDescribe: this._$currentDefectDescribe,
        Remark: this._$remark,
      });
      this.dispose();
    }

    /**
     * 更改缺陷大类
     */
    public _$defectClassChange(): void {
      // console.log("初始化子类选项");
      this._$defectSubClasses = Object.keys(this._$defectClassObject[this._$currentDefectClass]);
      this._$currentDefectSubClass =this._$defectSubClasses[0];
      this._$defectSubClassChange();
    }

    /**
     * 更改缺陷小类
     */
    public _$defectSubClassChange(): void {
      // console.log("初始化缺陷细项选项");
      let defectDescribes = this._$defectClassObject[this._$currentDefectClass][this._$currentDefectSubClass];
      this._$defectDescribes = defectDescribes.map(item => {
        return item.description;
      });
      this._$currentDefectDescribe = this._$defectDescribes[0];
    }

    public _$quickSearchChange(): void {
      const filterData = this._$defectAllDescribeData.filter(describe => describe.description == this._$quickSearch);
      if(filterData && filterData.length > 0) {
        this._$currentDefectClass = filterData[0].defClass;
        this._$defectSubClasses = Object.keys(this._$defectClassObject[this._$currentDefectClass]);
        this._$currentDefectSubClass = filterData[0].subDefClass;
        let defectDescribes = this._$defectClassObject[this._$currentDefectClass][this._$currentDefectSubClass];
        this._$defectDescribes = defectDescribes.map(item => {
          return item.description;
        });
        this._$currentDefectDescribe = this._$quickSearch;
      }
    }
}

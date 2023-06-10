import {EventEmitter, Input, Output, Component} from "@angular/core";
import {TargetInfo} from "./target-types";
import {InitData} from "./types";

@Component({
  template: ''
})

export class BaseOperate {
    @Input()
    public initData: InitData;

    @Output()
    public answer: EventEmitter<TargetInfo> = new EventEmitter<TargetInfo>();
}

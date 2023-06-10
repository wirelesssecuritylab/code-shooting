import { Component, OnInit, ViewChild, Input, Output, EventEmitter, SimpleChanges } from '@angular/core';
import { PlxModal, PlxBreadcrumbItem } from 'paletx';
import { ListService } from '../../list/list.service';

const USER = 'user';
const TEST = 'test';

@Component({
  selector: 'app-tabs',
  templateUrl: './tabs.component.html',
  styleUrls: ['./tabs.component.css'],
})
export class TabsComponent implements OnInit {
  @ViewChild('content') public content;
  @ViewChild('plxTabset') public plxTabset;

  public modal: any;

  public tabs = [];
  public closeIdArr = [];
  public titleDesc = '';
  public title = '确定';
  public userId: string;
  public targets: any[];
  public breadModel: PlxBreadcrumbItem[] = [];
  public rangeName: string ;
  public rangeType: string;
  public nowTime: Date;
  public isEnd: boolean = true;
  public rangeInfo: any = {};
  public valid: boolean = false;
  public activeTabId : string;
  public entranceTabId: string;

  @Input()
  set targetId(value: string) {
    console.log("input targetid changed" + value);
    this.activeTabId = value;
    this.entranceTabId = value;
  }

  @Input() public language: string;
  @Input() public rangeId: string;
  @Input() public isMyRange: boolean;
  @Output() public targetChanged1 = new EventEmitter<any>();
  constructor(
    private modalService: PlxModal,
    private listService: ListService,
  ) { }

  ngOnInit() {
    this.userId = localStorage.getItem(USER);
    this.getTargetsInfo();
  }

  ngOnChanges(changes:SimpleChanges) {
    console.log(changes)
  }
  getTargetsInfo() {
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
          };
          this.targets = this.rangeInfo.targets.filter(
            (target) =>
              target.language.toUpperCase() == this.language.toUpperCase()
          );
          if (this.targets.length > 0) {
            for (let i = 0; i < this.targets.length; i++) {
              const tab = {};
              tab['title'] = this.targets[i].targetName; //targetName
              tab['id'] = this.targets[i].targetId; //
              this.tabs.push(tab);
            }
            this.activeTabId = this.tabs[0].id;
            this.entranceTabId = this.activeTabId;
          }
        }
      });
  }

  tabChange(event: any) {
    this.entranceTabId = event.nextId;
    this.targetChanged1.emit(event.nextId);
    console.log(this.targetChanged1)
  }

  public openModal(content) {
    this.modal = this.modalService.open(content);
  }

  showShootTab(): boolean {
    return !(
      this.rangeInfo.type == TEST ||
      (this.nowTime.getTime() >= this.rangeInfo.start_at &&
        this.nowTime.getTime() <= this.rangeInfo.end_at)
    );
  }

  End(): boolean {
    if (this.nowTime.getTime() >= this.rangeInfo.start_at &&
      this.nowTime.getTime() <= this.rangeInfo.end_at) {
      this.isEnd = true;
    }

    else {
      this.isEnd = false;
    }
    return this.isEnd
  }
}

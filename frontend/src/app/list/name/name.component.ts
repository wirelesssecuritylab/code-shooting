import {ChangeDetectorRef, Component, Injector, OnInit} from '@angular/core';
import {PlxTableService} from "paletx";
const TEST = '练习';
@Component({
  selector: 'app-name',
  templateUrl: './name.component.html',
  styleUrls: ['./name.component.css']
})
export class NameComponent implements OnInit {
  public rowData: any;
  public name: any;
  public start_at: Date;
  public end_at: Date;
  public type: any;

  constructor(private plxTableService: PlxTableService, private inject: Injector, private cdr: ChangeDetectorRef) {
    this.rowData = this.plxTableService.getPositionInfo().currentRowData;
  }

  ngOnInit(): void {
    this.type = this.rowData.type;
    this.name = this.rowData.name;
    this.start_at = this.rowData.start_at;
    this.end_at = this.rowData.end_at;
  }

  showTime(): boolean {
    return !(this.type == TEST)
  }
}

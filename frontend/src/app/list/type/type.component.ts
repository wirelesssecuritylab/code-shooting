import {ChangeDetectorRef, Component, Injector, Input, OnInit} from '@angular/core';
import {PlxTableService} from "paletx";

@Component({
  selector: 'app-type',
  templateUrl: './type.component.html',
  styleUrls: ['./type.component.css']
})
export class TypeComponent implements OnInit {
  public rowData: any;
  public type: any;
  public people: any;
  public languages: string;

  constructor(private plxTableService: PlxTableService, private inject: Injector, private cdr: ChangeDetectorRef) {
    this.rowData = this.plxTableService.getPositionInfo().currentRowData;
  }

  ngOnInit(): void {
    this.type = this.rowData.type;
    this.languages = this.rowData.languages.join();
  }

}

import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { PlxTableService, PlxModal, PlxMessage} from 'paletx';
import { ManageService } from '../../manage.service';

const _DELETE = 'remove';

@Component({
  selector: 'app-range-operation',
  templateUrl: './range-operation.component.html',
  styleUrls: ['./range-operation.component.css']
})
export class RangeOperationComponent implements OnInit {
  public rowData: any;
  public modal:any;
  public isOpen = true;
  constructor(private plxTableService: PlxTableService,
              private router: Router,
              private activatedRoute: ActivatedRoute,
              private modalService: PlxModal,
              private plxMessageService: PlxMessage,
              private manageService: ManageService) {
    this.rowData = this.plxTableService.getPositionInfo().currentRowData;
  }

  ngOnInit(): void {
  }

  /**
   * 管理员进入靶场查询或导出成绩
   */
  goToRange() {
    let routeStr: string = '../query';
    this.router.navigate([routeStr], {
      queryParams: {
        'id': this.rowData.id,
        'name': this.rowData.name,
        'languages': this.rowData.languages,
        'projectId': this.rowData.project
      },
      relativeTo: this.activatedRoute
    });
  }

  /**
   * 编辑靶场
   */
  modifyRange() {
    let routeStr: string = '../edit';
    this.router.navigate([routeStr], {
      queryParams: {
        'id': this.rowData.id
      },
      relativeTo: this.activatedRoute
    });
  }

  /**
   * 删除靶子时打开确认删除模态框
   * @param content
  */
  openConfirmDlg(content) {
    this.modal = this.modalService.open(content, {size: 'xs', enterEventFunc: this.func.bind(this)});
    this.isOpen = true;
  }

  public func(): void {
    if(this.isOpen) {
      this.modal.close();
    }
    this.isOpen = !this.isOpen;
  }

  /**
   * 删除靶场
   * @returns
   */
  deleteRange() {
    if (!this.rowData.id) {
      return;
    }
    let deleteBody: any = {
      name: _DELETE,
      parameters: {
        id: this.rowData.id,
      }
    };
    this.manageService.manageRangeApi(deleteBody).subscribe(res => {
      this.plxMessageService.success("删除成功！", '');
      this.router.navigateByUrl('/main/manage/range/list');
    }, err => {
      this.plxMessageService.error('删除失败！', err.cause);
    });
  }
}

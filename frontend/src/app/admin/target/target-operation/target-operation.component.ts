import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { PlxTableService, PlxMessage, PlxModal} from 'paletx';
import { ManageService } from '../../manage.service';

const _DELETE = 'remove';

@Component({
  selector: 'app-target-operation',
  templateUrl: './target-operation.component.html',
  styleUrls: ['./target-operation.component.css']
})
export class TargetOperationComponent implements OnInit {

  public rowData: any;
  public modal:any;
  public isOpen = true;
  public userId: string;
  constructor(private router: Router, private activatedRoute: ActivatedRoute,
              private plxTableService: PlxTableService,
              private manageService: ManageService,
              private plxMessageService: PlxMessage,
              private modalService: PlxModal) {
    this.userId = localStorage.getItem('user');
    this.rowData = this.plxTableService.getPositionInfo().currentRowData;
   }

  ngOnInit(): void {
  }

  /**
   * 编辑靶子操作
   */
  modifyTarget() {
    let routeStr: string = '../edit';
    this.router.navigate([routeStr], {
      queryParams: {
        'id': this.rowData.id,
        'language': this.rowData.language
      },
      relativeTo: this.activatedRoute
    });
  }

  /**
   * 删除靶子时打开确认删除模态框
   * @param content
   */
  open(content) {
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
   * 确认删除靶子
   * @returns
   */
  deleteTarget() {
    if (!this.rowData.id) {
      return;
    }
    let deleteBody: any = {
      name: _DELETE,
      parameters: {
        id: this.rowData.id,
        owner: this.rowData.owner
      }
    };
    this.manageService.manageTargetApi(deleteBody).subscribe(res => {
      this.plxMessageService.success("删除成功！", '');
      this.router.navigateByUrl('/main/manage/target/list');
    }, err => {
      this.plxMessageService.error('删除失败！', err.cause);
    });
  }
}

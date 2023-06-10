import { Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {JigsawTreeExt, SimpleTreeData, ZTreeSettingSetting} from "@rdkmaster/jigsaw";
import {DefectService} from "../../services/defect-service";
import {FileTreeNode} from "../../misc/types";
import { ManageService } from '../../../admin/manage.service';
import { PlxMessage } from 'paletx';
import { ActivatedRoute } from '@angular/router';

@Component({
    selector: 'tp-files-tree',
    templateUrl: './files-tree.component.html',
    styleUrls: ['./files-tree.component.scss']
})
export class FilesTreeComponent implements OnInit {
    public _$settings: ZTreeSettingSetting;
    public _$treeData: SimpleTreeData = new SimpleTreeData();

    private _selectedFile: FileTreeNode;

    @ViewChild('fileTree')
    private _fileTree: JigsawTreeExt;

    @Input()
    public get selectedFile(): FileTreeNode {
        return this._selectedFile;
    }

    public set selectedFile(value: FileTreeNode) {
        this._selectedFile = value;
        this._selectedNode();
    }

    //  输入属性：靶子id
    @Input() targetId: string;
    public userId: string;
    private fileNameList: string[] = [];
    private fileTree: any[] = [];
    public rangeId: string;

    @Output()
    public selectedFileChange: EventEmitter<FileTreeNode> = new EventEmitter<FileTreeNode>();

    constructor(private _defectService: DefectService, private manageService: ManageService,
                private plxMessageService: PlxMessage, private activatedRoute: ActivatedRoute) {
    }

    ngOnInit(): void {
      this.userId = localStorage.getItem('user');
      this.rangeId = this.activatedRoute.snapshot.queryParams.rangeId;
        // this._$treeData.fromObject(this._defectService.getSources());
      this._$settings = {
          data: {
              key: {
                  children: 'nodes',
                  name: 'fileName'
              }
          },
          callback: {
              onClick: (event, treeId, treeNode) => {
                  this.selectedFileChange.emit(treeNode);
              }
          }
      };
      this.getTargetFileName();
    }

    ngAfterViewInit(): void {
        // 默认选中第一个代码文件
        setTimeout(() => {
            this.selectedFile = this._fileTree.ztree.getNodesByFilter(node => !!node.code, true);
            this.selectedFileChange.emit(this.selectedFile);
        })
    }


    private _selectedNode(): void {
        if (!this._fileTree) {
            return;
        }
        const node = this._fileTree.ztree.getNodeByParam('fileName', this.selectedFile?.fileName);
        if (!node) {
            this._fileTree.ztree.getSelectedNodes().forEach(item => this._fileTree.ztree.cancelSelectedNode(item));
            return;
        }
        this._fileTree.ztree.selectNode(node);
    }

  /**
   * 通过查询靶子详情获取靶子下的靶子文件以及文件内容（如代码或md）
   * 目前先借用靶子管理中的接口，普通用户访问靶子详情应该单独的接口
   */
  private getTargetFileName() {
    if (!this.targetId) {
      return;
    }
    let reqBody: any = {
      name: 'shoot',
      parameters: {
        rangeid: this.rangeId,
        id: this.targetId,
        user: this.userId
      }
    };
    this.manageService.manageTargetApi(reqBody).subscribe(res => {
      if (!res || !res.detail || !Array.isArray(res.detail) || res.detail.length != 1) {
        return;
      }
      let detail: any = res.detail[0];
      this.fileNameList = this.fileNameList.concat(detail.targets);
      this.initFileTree();
    }, error => {
      this.plxMessageService.error('获取靶子详情失败！', JSON.parse(error.error).message);
    });
  }

  /**
   *
   * @returns 初始化左侧文件树
   */
  private initFileTree() {
    if (this.fileNameList.length === 0) {
      return;
    }
    this.fileNameList.forEach(item => {
      this.getTargetFileContent(item)
    });
  }

  /**
   * 获取某个靶子文件的内容（代码）
   * @param fileName
   * @returns
   */
  private getTargetFileContent(fileName: string) {
    if (!this.targetId || !fileName) {
      return;
    }
    this.manageService.getTargetFileApi(this.targetId, fileName).subscribe(res => {
      let node = {
        fileName: fileName,
        code: res
      };
      if (fileName.substr(-3, 3) === '.md') {
        node['type'] = 'md'
      }
      this.fileTree.push(node);
      this._$treeData.fromObject(this.fileTree);
      if (this.fileTree.length == this.fileNameList.length) {
        // 默认选中第一个代码文件
        this.ngAfterViewInit();
      }
    }, error => {

      this.plxMessageService.error('获取靶子文件失败！', JSON.parse(error.error).message);
    });
  }
}

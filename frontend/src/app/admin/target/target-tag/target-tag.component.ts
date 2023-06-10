import { Component, OnInit, ViewChild, Input, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-target-tag',
  templateUrl: './target-tag.component.html',
  styleUrls: ['./target-tag.component.css']
})
export class TargetTagComponent implements OnInit {

  constructor() { }

  @ViewChild('plxSelectTree') plxSelectTree: any;
  @ViewChild('plxSelect') plxSelect: any;
  // 已经选中的节点的key值
  @Input() selectedNodeKey: string;
  // 已经选择的标签选项
  @Input() selectedOptions: Array<string>;
  @Input() templateTreeNodes:any;
  @Output() public tagOptionChanged = new EventEmitter<any>();
  public options = [];

  ngOnInit() {
    const tmpOptions = [];
    this.getOptions(this.treeNodes, tmpOptions);
    this.options = tmpOptions;
    setTimeout(() => {
      this.treeNodes = [].concat(this.treeNodes);
    });
  }

  public filterChanged(filterValue: string): void {
    this.plxSelectTree.search(filterValue);
  }

  public onClose(): void {
    this.plxSelectTree.search('');
  }

  public onOpen(): void {
    // 每次打开时，重新设置树节点引用，才能处理默认选中节点。否则处理不了。
    this.treeNodes = [].concat(this.treeNodes);
  }

  public deselected(value: any): void {
    for (let i = 0; i < this.treeNodes.length; i++) {
      if (this.treeNodes[i].children && this.treeNodes[i].children.length) {
        this.treeNodes[i].children.forEach(node => {
          if (value.value == node.key) {
            node['isSelected'] = false;
            node['parent']['isSelected'] = false;
          }
        });
      } else {
        if (value.value == this.treeNodes[i].key) {
          this.treeNodes[i]['isSelected'] = false;
        }
      }

    }
    this.treeNodes = [].concat(this.treeNodes);
  }

  selectionChange(val) {
    console.log('selectionChange', val);
    console.log(this.options);

    const tmp = [];
    if (val) {
      for (let i = 0; i < val.length; i++) {
        if (val[i]['isLeaf']) {
          tmp.push(val[i].key);
        }
      }
    }
    setTimeout(() => {
      this.selectedOptions = [].concat(tmp);
      this.tagOptionChanged.emit(this.selectedOptions);
    });
  }

  private getOptions(nodes: any, tmpOptions: any): void {
    for (let i = 0; i < nodes.length; i++) {
      tmpOptions.push(this.getSelectItem(nodes[i]));
      if (nodes[i]['children']) {
        this.getOptions(nodes[i]['children'], tmpOptions);
      }
    }
  }

  getSelectItem(node: any): any {
    return {
      value: node['key'],
      label: node['label']
    };
  }

  public treeNodes = [
    {
      label: '功能',
      key: '00',
      isExpanded: true,
      children: [
        {
          label: '测试',
          isLeaf: true,
          key: '01',
          isExpanded: true,
          children: [
            {
              label: '缺少测试用例',
              isLeaf: true,
              key: '0.0.0.0',
            },
            {
              label: '缺少校验点/分支校验',
              isLeaf: true,
              key: '0.0.0.1',
            },
            {
              label: '测试用例不正交',
              isLeaf: true,
              key: '0.0.0.2',
            },
            {
              label: '测试用例冗余',
              isLeaf: true,
              key: '0.0.0.3',
            },
            {
              label: '测试用例职责不单一',
              isLeaf: true,
              key: '0.0.0.4',
            }
          ]
        },
        {
          label: '实现',
          isLeaf: true,
          key: '02',
          isExpanded: true,
          children: [
            {
              label: '正常流程实现不正确',
              isLeaf: true,
              key: '0.0.1.0',
            },
            {
              label: '正常流程实现缺失',
              isLeaf: true,
              key: '0.0.1.1',
            },
            {
              label: '可选流程实现不正确',
              isLeaf: true,
              key: '0.0.1.2',
            },
            {
              label: '可选流程实现缺失',
              isLeaf: true,
              key: '0.0.1.3',
            },
            {
              label: '异常流程实现不正确',
              isLeaf: true,
              key: '0.0.1.4',
            },
            {
              label: '异常流程实现缺失',
              isLeaf: true,
              key: '0.0.1.5',
            },
          ]
        }
      ]
    },
    {
      label: '性能',
      isLeaf: false,
      key: '03',
      isExpanded: true,
      children: [
        {
          label: '计算',
          isLeaf: false,
          key: '031',
          isSelected: false,
          isExpanded: true,
          children: [
            {
              label: '不必要的等待',
              key: '0.4.0.0',
              isLeaf: true,
            },
            {
              label: '重复计算',
              key: '0.4.0.1',
              isLeaf: true,
            },
            {
              label: '没有使用预先计算',
              key: '0.4.0.2',
              isLeaf: true,
            },
            {
              label: '没有使用延迟计算',
              key: '0.4.0.3',
              isLeaf: true,
            },
            {
              label: '没有使用分摊计算',
              key: '0.4.0.4',
              isLeaf: true,
            },
            {
              label: '没有使用批量计算',
              key: '0.4.0.5',
              isLeaf: true,
            },
            {
              label: '不必要的浮点运算',
              key: '0.4.0.6',
              isLeaf: true,
            },
            {
              label: '线程数量与核数配比不合理',
              key: '0.4.0.7',
              isLeaf: true,
            }
          ]
        },
        {
          label: '内存',
          isLeaf: false,
          key: '032',
          isSelected: false,
          isExpanded: true,
          children: [
            {
              label: '不必要的动态内存申请和释放',
              key: '0.4.1.0',
              isLeaf: true
            },
            {
              label: '不合理的传值（pass-by-value）',
              key: '0.4.1.1',
              isLeaf: true
            },
          ]
        },
        {
          label: 'Cache',
          isLeaf: false,
          key: '033',
          isSelected: false,
          isExpanded: true,
          children: [
            {
              label: '不必要的进程/线程上下文切换',
              key: '0.4.2.0',
              isLeaf: true
            },
            {
              label: '不必要的访问数据时缓存切换',
              key: '0.4.2.1',
              isLeaf: true
            },
          ]
        },
        {
          label: '编译',
          isLeaf: false,
          key: '034',
          isSelected: false,
          isExpanded: true,
          children: [
            {
              label: '没有使用编译优化选项',
              key: '0.4.3.0',
              isLeaf: true
            }
          ]
        },
      ]
    },
    {
      label: '可靠性',
      key: '04',
      isExpanded: true,
      isSelected: false,
      children: [
        {
          label: '异常',
          isLeaf: false,
          key: '041',
          isSelected: false,
          isExpanded: true,
          children: [
            {
              label: '死锁',
              key: '0.5.0.0',
              isLeaf: true
            },
            {
              label: '磁盘操作错误',
              key: '0.5.0.1',
              isLeaf: true
            },
            {
              label: '算术操作错误',
              key: '0.5.0.2',
              isLeaf: true
            },
            {
              label: '系统资源耗尽',
              key: '0.5.0.3',
              isLeaf: true
            },
          ]
        },
        {
          label: '错误',
          isLeaf: false,
          key: '042',
          isSelected: false,
          isExpanded: true,
          children: [
            {
              label: '局部变量覆盖同名的全局变量',
              key: '0.5.1.0',
              isLeaf: true
            },
            {
              label: '错误信息丢失',
              key: '0.5.1.1',
              isLeaf: true
            }
          ]
        },
        {
          label: '容错',
          isLeaf: false,
          key: '043',
          isSelected: false,
          isExpanded: true,
          children: [
            {
              label: '重试问题',
              key: '0.5.2.0',
              isLeaf: true
            },
            {
              label: '限流问题',
              key: '0.5.2.1',
              isLeaf: true
            },
            {
              label: '降级问题',
              key: '0.5.2.2',
              isLeaf: true
            },
            {
              label: '熔断问题',
              key: '0.5.2.3',
              isLeaf: true
            }
          ]
        },
        {
          label: '可用性',
          isLeaf: false,
          key: '044',
          isSelected: false,
          isExpanded: true,
          children: [
            {
              label: 'HA考虑不全面',
              key: '0.5.3.0',
              isLeaf: true
            },
            {
              label: '负载均衡考虑不全面',
              key: '0.5.3.1',
              isLeaf: true
            },
            {
              label: '数据不一致',
              key: '0.5.3.2',
              isLeaf: true
            }
          ]
        },
        {
          label: '可恢复性',
          isLeaf: false,
          key: '045',
          isSelected: false,
          isExpanded: true,
          children: [
            {
              label: '自愈失败',
              key: '0.5.4.0',
              isLeaf: true
            }
          ]
        },
      ]
    },
    {
      label: '轮船',
      key: '07',
      children: [
        {
          label: '客船', isLeaf: true,
          key: '08',
        },
        {
          label: '公务船', isLeaf: true,
          key: '09',
        },
        {
          label: '货船', isLeaf: true,
          key: '10',
        },
        {
          label: '客船1', isLeaf: true,
          key: '11',
        },
        {
          label: '公务船1', isLeaf: true,
          key: '12',
        },
        {
          label: '货船1', isLeaf: true,
          key: '13',
        },
        {
          label: '客船2', isLeaf: true,
          key: '14',
        },
        {
          label: '公务船2', isLeaf: true,
          key: '15',
        },
        {
          label: '货船2', isLeaf: true,
          key: '16',
        },
      ]
    },
    {
      label: '自定义下拉面板超长自定义标签，自定义下拉面板超长自定义标签，自定义下拉面板超长自定义标签', isLeaf: true,
      key: '17',
    },
  ];

}

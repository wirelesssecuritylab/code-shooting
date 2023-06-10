import {
  ChangeDetectorRef,
  Component,
  Injector, Input,
  OnDestroy,
  OnInit,
} from '@angular/core';
import {
  FieldType,
  ModalDismissReasons,
  PlxModal,
  PlxTableService,
} from 'paletx';
import { Router, ActivatedRoute } from '@angular/router';

const TEST = '练习';
const ROW = 'row';
const LANGUAGE = 'language';
const TARGETADDR = 'targetAddr';
@Component({
  selector: 'app-operation',
  templateUrl: './operation.component.html',
  styleUrls: ['./operation.component.css'],
})
export class OperationComponent implements OnInit, OnDestroy {
  public rowData: any;
  public start_at: Date;
  public end_at: Date;
  public type: any;
  public nowTime: Date;
  public timeStr: string;
  public formSetting: any;
  public modal: any;
  lansSelector: any[];
  @Input() isMyRange: boolean;

  constructor(
    private plxTableService: PlxTableService,
    private inject: Injector,
    public router: Router,
    private modalService: PlxModal,
  ) {
    this.rowData = this.plxTableService.getPositionInfo().currentRowData;
  }

  ngOnInit(): void {
    this.initSelector();
    this.formSetting = this.initFormSetting();
    this.nowTime = new Date(Date.parse(new Date().toString()));
    this.type = this.rowData.type;
    this.start_at = this.rowData.start_at;
    this.end_at = this.rowData.end_at;

    this.timeDown();
  }

  timeDown() {
    if (this.type != TEST) {
      let curTime = Math.floor(
        (this.start_at.getTime() - this.nowTime.getTime()) / 1000
      );
      const _this = this;
      if (curTime > 0) {
        const timer = setInterval(function () {
          curTime = curTime - 1;
          if (curTime === 0) {
            clearInterval(timer);
          } else {
            _this.timeStr = _this.timediff(curTime);
          }
        }, 1000);
      } else {
      }
    }
  }

  timediff(secondAll) {
    const day = Math.floor(secondAll / (60 * 60 * 24));
    const hour = Math.floor((secondAll - day * 60 * 60 * 24) / (60 * 60));
    const min = Math.floor(
      (secondAll - day * 60 * 60 * 24 - hour * 60 * 60) / 60
    );
    const sec = Math.floor(
      secondAll - day * 60 * 60 * 24 - hour * 60 * 60 - min * 60
    );
    let timeString = '';
    if (day > 0) {
      timeString = day + '天';
    }
    if (hour > 0) {
      timeString = timeString + hour + ':';
    }
    if (min > 0) {
      timeString = timeString + min + ':';
    }
    return timeString + sec;
  }

  showBtn(): boolean {
    return !(
      this.type == TEST || this.nowTime.getTime() >= this.start_at.getTime()
    );
  }

  showDdl(): boolean {
    return (
      this.type !== TEST && this.nowTime.getTime() <= this.start_at.getTime()
    );
  }

  /*action() {
    localStorage.setItem("row", JSON.stringify(this.rowData));
    this.router.navigateByUrl('main/user/action')
  }

  actionCancel() {
    localStorage.removeItem(TARGETADDR);
    localStorage.removeItem(LANGUAGE);
    localStorage.removeItem(ROW);
  }*/

  ngOnDestroy(): void {
    this.modalService.destroyModalInstance();
  }

  public isOpen = true;
  public openWithoutDestroy(content) {
    const size: 'sm' | 'lg' | 'xs' = 'sm';
    const options = {
      size: size,
      // contentClass: 'plx-modal-custom-content',
      enterEventFunc: this.func.bind(this),
      escCallback: this.escCallback.bind(this),
      destroyOnClose: false,
      modalId: 'plx-modal-1',
      openCallback: () => {
        // console.log('open');
      },
    };
    this.modal = this.modalService.open(content, options);
    this.isOpen = true;
  }

  public open(content) {
    const size: 'sm' | 'lg' | 'xs' = 'sm';
    const options = {
      size: size,
      enterEventFunc: this.func.bind(this),
      escCallback: this.escCallback.bind(this),
      openCallback: () => {
        // console.log('open');
      },
    };
    this.modal = this.modalService.open(content, options);
    this.isOpen = true;
  }

  private escRe = true;
  public escCallback(): boolean {
    // this.escRe = !this.escRe;
    console.info('escCallback');
    return this.escRe;
  }

  public func(): void {
    console.log('enter event');
    if (this.isOpen) {
      this.modal.close();
    }
    this.isOpen = !this.isOpen;
  }

  private getDismissReason(reason: any): string {
    if (reason === ModalDismissReasons.ESC) {
      return 'by pressing ESC';
    } else if (reason === ModalDismissReasons.BACKDROP_CLICK) {
      return 'by clicking on a backdrop';
    } else {
      return `with: ${reason}`;
    }
  }

  public cancel(): void {
    //this.actionCancel()
    this.modalService.destroyModalInstance();
    this.modal.close();
  }

  submit() {
    const isValid = this.formSetting.validateFields();
    if (isValid) {
      console.log('form value:', this.formSetting.formObject.value);
      let language = '';
      language = this.formSetting.formObject.value.language;
      /*let targetAddr = '';
      this.rowData.targets.forEach(function (value) {
        if (value.language == language) {
          targetAddr = value.targetAddr;
        }
      });*/
      /*localStorage.setItem("language", language);
      localStorage.setItem("targetAddr", targetAddr);
      this.action();*/
      // this.router.navigateByUrl(`main/user/action?rangeId=${this.rowData.targetId}&language=${encodeURIComponent(language)}&rangeName=${this.rowData.name}`);
      if(this.isMyRange) {
        this.router.navigateByUrl(`main/user/shoot-range?rangeId=${this.rowData.targetId}&language=${encodeURIComponent(language)}&rangeName=${this.rowData.name}&rangeType=${this.type}&isMyRange=true`);
      } else {
        this.router.navigateByUrl(`main/user/shoot-range?rangeId=${this.rowData.targetId}&language=${encodeURIComponent(language)}&rangeName=${this.rowData.name}&rangeType=${this.type}`);
      }
      this.modal.close();
    }
  }

  initSelector() {
    let s = [];
    this.rowData.languages.forEach(function (value) {
      s.push({ label: value, value: value });
    });
    this.lansSelector = s;
  }

  initFormSetting() {
    return {
      isShowHeader: false,
      isGroup: false,
      labelClass: 'col-sm-3',
      componentClass: 'col-sm-8',
      srcObj: {
        id: 1,
        authType: '0',
      },
      fieldSet: [
        {
          fields: [
            {
              name: 'name',
              label: '靶场',
              type: FieldType.TEXT,
              disabled: true,
              text: this.rowData.name,
            },
            {
              name: 'language',
              label: '选择语言',
              type: FieldType.SELECTOR,
              binding: {
                dropdownContainer: 'body',
                scrollSelectors: ['.modal-body'],
              },
              valueSet: this.lansSelector,
              required: true,
            },
          ],
        },
      ],
    };
  }
}

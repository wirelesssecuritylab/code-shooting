import {
  Component,
  OnInit,
  Input,
  ElementRef,
  ViewChild,
  AfterViewInit,
  OnDestroy,
  ViewContainerRef
} from '@angular/core';
import { Router } from '@angular/router';
import { LoginService } from './login.service';

import {
  NgxQrcodeElementTypes,
  NgxQrcodeErrorCorrectionLevels,
} from '@techiediaries/ngx-qrcode';
import { PlxMessage } from 'paletx';

const adminPrivilege: string[] = ['submitStandardAnswer', 'submitTemplate', 'createRange', 'deleteRange', 'viewRangeScore'];
const expertPrivilege: string = 'editTargetTag';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
})

export class LoginComponent implements OnInit, AfterViewInit, OnDestroy {
  @Input() size: any;
  public passwordInput;
  public userInput;
  public depart;
  public departId;
  public chName;
  public qrcode;
  public time: any;
  public qrcodeTime: any;
  @ViewChild('tabmenu') tabmenu: ElementRef;
  public elementType: NgxQrcodeElementTypes;
  public correctionLevel: NgxQrcodeErrorCorrectionLevels;
  public isHavePrivilege: boolean = true;
  public userRole: string = 'shootuser';
  constructor(
    public router: Router,
    private logService: LoginService,
    private vRef: ViewContainerRef,
    private plxMessageService: PlxMessage
  ) // private http: HttpClient
  {

  }

  ngOnDestroy(): void {
    clearInterval(this.time);
    clearInterval(this.qrcodeTime);
  }

  ngOnInit(): void {
    // this.draw();
    //localStorage.clear()

    if (localStorage.getItem('user')) {
      console.log(localStorage.getItem('user'))
      this.router.navigate(['/main']);
    }
    const portalUser = "admin"
    this.userInput = portalUser
    this.logService.verify(portalUser).subscribe(
      (res) => {
        if (res?.code === 200) {
          this.login()
        }
      }
    );
  }


  login() {
    this.logService
      .getUserInfo({ id: this.userInput })
      .subscribe((res) => {
        console.log(res)
        /**
         * 增加用户权限判断，若无任何权限，不允许登录即提示无权限
         */
        if (!res.detail || !res.detail.privileges || res.detail.privileges.length === 0) {
          return;
        }
        this.isHavePrivilege = true;
        let privileges = res.detail.privileges;
        for (let i = 0; i < privileges.length; i++) {
          if (adminPrivilege.indexOf(privileges[i]) >= 0) {
            this.userRole = 'admin';
            break;
          }
        }
        if (privileges.includes(expertPrivilege)) {
          localStorage.setItem('expertPrivilege', expertPrivilege);
        }
        localStorage.setItem('role', this.userRole);
        localStorage.setItem('user', this.userInput);
        this.depart = res?.detail?.department;
        this.chName = res?.detail?.name;
        this.departId = res?.detail?.orgID;
        localStorage.setItem('department', res?.detail?.department);
        localStorage.setItem('institute', res?.detail?.institute);
        localStorage.setItem('name', res?.detail?.name);
        localStorage.setItem('org_id', res?.detail?.orgID);
        this.router.navigateByUrl('main/user/list');


      });
  }

  ngAfterViewInit() { }

}

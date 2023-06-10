import {Component, OnInit} from '@angular/core';
import {Router} from "@angular/router";
import {HttpClient} from "@angular/common/http";
import {PlxMessage} from "paletx";

@Component({
  selector: 'app-userinfo',
  templateUrl: './userinfo.component.html',
  styleUrls: ['./userinfo.component.css']
})
export class UserinfoComponent implements OnInit {

  userId: string;
  userInfo: any;
  isEdit = false;
  department: string;
  centerName: string;
  institute: string;
  teamName: string;
  email: string;
  isEmailValid = true;

  constructor(private router: Router, private http: HttpClient,private plxMessage: PlxMessage) {
    this.userId = localStorage.getItem('user');
    if (!this.userId || this.userId == 'null') {
      localStorage.clear();
      this.router.navigateByUrl('/login');
    }
  }

  ngOnInit(): void {
    this.initPersonalData();
  }

  initPersonalData() {
    const params = {
      name: "query",
      parameters: {
        id: this.userId
      }
    };
    this.http.post(`/api/code-shooting/v1/actions/person`, params).subscribe((info:any) => {
        this.userInfo = info;
        this.department = info.department;
        this.centerName = info.centerName;
        this.institute = info.institute;
        this.teamName = info.teamName;
        this.email = info.email;
    });
  }

  editInfo() {
    this.isEdit = true;
  }

  inputEmail() {
    this.isEmailValid = true;
  }

  save() {
    if(this.isEmailValid) {
      const params = {
        name: "modify",
        parameters: {
          id: this.userId,
          department: this.department,
          centerName: this.centerName,
          institute: this.institute,
          teamName: this.teamName,
          email: this.email
        }
      };
      this.http.post(`/api/code-shooting/v1/actions/person`, params).subscribe(() => {
        this.userInfo.department = this.department;
        this.userInfo.centerName = this.centerName;
        this.userInfo.institute = this.institute;
        this.userInfo.teamName = this.teamName;
        this.userInfo.email = this.email;
        this.isEdit = false;
      }, err => {
        this.plxMessage.error('保存个人信息失败！', err.msg);
      });
    }
  }

  cancel() {
    this.isEdit = false;
  }

  validateEmail() {
    const emailPattern = /^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$/;
    if(!emailPattern.test(this.email)) {
      this.isEmailValid = false;
    }
  }

  refreshInfo() {
    const params = {
      name: "refresh",
      parameters: {
        id: this.userId
      }
    };
    this.http.post(`/api/code-shooting/v1/actions/person`, params).subscribe((info:any) => {
      this.department = info.department;
      this.centerName = info.centerName;
      this.institute = info.institute;
      this.teamName = info.teamName;
    });
  }
}

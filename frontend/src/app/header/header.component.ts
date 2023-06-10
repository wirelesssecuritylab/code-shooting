import { Component, OnInit } from '@angular/core';
import { Router } from "@angular/router";
@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {
  public showDropdown = false;
  public userName: string;
  public userId: string;
  public currentUser: string;
  constructor(private router: Router) { }

  ngOnInit(): void {
    // this.userName = localStorage.getItem('name') !== 'undefined'
    // && localStorage.getItem('name') !== '' ? localStorage.getItem('name') : localStorage.getItem('user');
    this.userName = localStorage.getItem('name');
    this.userId = localStorage.getItem('user');
    this.currentUser = this.userName + this.userId;
  }

  logout() {
    localStorage.clear();
    this.router.navigateByUrl('/login');
  }

  getCookie(key) {
    let match = document.cookie.match(new RegExp('(^|;|\\s)' + key + '=([^;]+)'));
    if (!match) {
      return '';
    }
    return decodeURIComponent(match[2]);
  }

  help() {
    this.router.navigateByUrl('/main/user/doc');
  }

  personalcenter() {
    this.router.navigateByUrl('/main/user/personalcenter');
  }

  onEnter() {
    this.showDropdown = true;
  }

  onLeave() {
    this.showDropdown = false;
  }

}

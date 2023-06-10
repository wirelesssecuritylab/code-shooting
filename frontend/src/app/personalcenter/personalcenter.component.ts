import {Component, OnInit} from '@angular/core';
import {NavigationEnd, Router} from "@angular/router";

@Component({
  selector: 'app-personalcenter',
  templateUrl: './personalcenter.component.html',
  styleUrls: ['./personalcenter.component.css']
})
export class PersonalcenterComponent implements OnInit {

  menuWidth = 'width: 200px';
  menuConfig = {
    toggleBtn: {
      text: '收起',
      collapsedText: '侧边栏展开'
    },
    menusData: [
      {
        name: '个人信息',
        iconClass: 'plx-ico-sm-user-16',
        href: '/main/user/personalcenter/userinfo',
        selected: true
      },
      {
        name: '我的靶场',
        iconClass: 'plx-ico-reporting-tool-f-24',
        href: '/main/user/personalcenter/myrange',
      },
      {
        name: '我的练习',
        iconClass: 'plx-ico-comment-f-16',
        href: '/main/user/personalcenter/freepractice',
      },
      {
        name: '我的靶子',
        iconClass: 'plx-ico-reason-code-16',
        href: '/main/user/personalcenter/mytarget',
      },
    ]
  };

  constructor(private router: Router) {
    this.router.events.subscribe(event => {
      if(event instanceof NavigationEnd){
        this.menuConfig.menusData.forEach(menu => {
          if(menu.href == event.urlAfterRedirects) {
            menu.selected = true;
          } else {
            menu.selected = false;
          }
        })
      }
    })
  }

  ngOnInit(): void {
  }

  menuClick(e) {
    this.router.navigateByUrl(e.menu.href);
  }

  collapseBtnClick(e) {
    if(e) {
      this.menuWidth = 'width: 40px';
    } else {
      this.menuWidth = 'width: 200px';
    }
  }
}

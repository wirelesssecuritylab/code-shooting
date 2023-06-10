import { Component, OnInit } from '@angular/core';
import { Router } from "@angular/router";
import { HeaderMenu } from './header-menu.class';

@Component({
  selector: 'app-header-menu',
  templateUrl: './header-menu.component.html',
  styleUrls: ['./header-menu.component.css']
})
export class HeaderMenuComponent implements OnInit {

  public menuList = new HeaderMenu().defaultMenu;
  public menuIdx: number = 0;
  public userRole: string;

  constructor(private router: Router) {
    this.userRole = localStorage.getItem('role')
  }

  ngOnInit(): void {
    this.filterMenuList();
    this.judgeCurrentRoute();
  }

  /**
   * 获取当前路由对应的菜单项的下标
   */
  judgeCurrentRoute() {
    let curRouterUrl = this.router.routerState.snapshot.url;
    this.menuIdx = this.menuList.findIndex((item: any) => item.menuLink === curRouterUrl);
  }


  /**
   * 根据用户所拥有的权限控制顶部菜单的显示
   */
  filterMenuList(): void {
    this.menuList = this.menuList.filter(filterItem => filterItem.menuRole.indexOf(this.userRole) >= 0);
  }

  /**
   * 点击顶部菜单
   * @param val 被点击菜单项在菜单列表中的下标
   */
  clickMenu(val: any): void {
    this.menuIdx = val;
    this.router.navigateByUrl(this.menuList[val].menuLink);
  }

}

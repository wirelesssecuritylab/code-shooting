import { Component, OnInit } from '@angular/core';
import {Router} from "@angular/router";
import {HttpClient} from "@angular/common/http";

@Component({
  selector: 'app-helpdoc',
  templateUrl: './help.component.html',
  styleUrls: ['./help.component.css']
})
export class HelpComponent implements OnInit{

  userRole: string;

  menuConfig = {
    toggleBtn: {
      text: '收起',
      collapsedText: '侧边栏展开'
    },
    menusData: [
    ]
  };

  mdContent: string;
  menuWidth = 'width: 200px';

  constructor(private router: Router,private http: HttpClient) {
    this.userRole = localStorage.getItem('role');
  }

  ngOnInit(): void {
    this.menuConfig.menusData.push({
      name: '用户使用指南',
      iconClass: 'plx-ico-sm-user-16',
      href: '/main/user/doc',
      selected: true,
      children: [
      ]
    });
    if(this.userRole == 'admin') {
      this.menuConfig.menusData.push({
        name: '系统管理指南',
        iconClass: 'plx-ico-sm-user-managment-16',
        href: '/main/user/doc',
        children: [
        ]
      });
    }
    this.loadMenus();
  }

  loadMenus() {
    this.http.get(`/api/code-shooting/v1/documents/catalogue/all/list`).subscribe((menus: any) => {
      if(menus && menus.children && Array.isArray(menus.children)) {
        menus.children.forEach(menu => {
          if(menu.name == 'system-manage' && this.userRole == 'admin') {
            this.menuConfig.menusData[1].children = [...this.parseMenus(menu)];
          } else if(menu.name == 'user-guide'){
            this.menuConfig.menusData[0].children = [...this.parseMenus(menu)];
          }
        });
      }
    });
  }

  parseMenus(menu: any): any[] {
    const menuDocs = menu.ducoments;
    if(menuDocs && Array.isArray(menuDocs)){
      return menuDocs.map(doc => {
        return {
          name: doc.endsWith('.md') ? doc.substring(0, doc.length -3) : doc,
          href: menu.id
        }
      });
    }
    return [];
  }

  menuClick(e) {
    this.http.post(`/api/code-shooting/v1/documents/detail/${e.menu.href}`, {filePath: e.menu.name + '.md'},{responseType: 'text'}).subscribe((res) =>{
      this.mdContent = res;
    });
  }

  collapseBtnClick(e) {
    if(e) {
      this.menuWidth = 'width: 40px';
    } else {
      this.menuWidth = 'width: 200px';
    }
  }

}

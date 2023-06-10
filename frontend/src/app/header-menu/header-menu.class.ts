/**
 * 声明顶部菜单类
 */
export class HeaderMenu {
  public defaultMenu: any = [
    {
      menuName: '靶场',
      menuRole: ['shootuser', 'admin'],
      menuLink: '/main/user/list'
    },
    {
      menuName: '自由练习',
      menuRole: ['shootuser', 'admin'],
      menuLink: '/main/user/test'
    },
    {
      menuName: '靶场管理',
      menuRole: ['admin'],
      menuLink: '/main/manage/range/list'
    },
    {
      menuName: '靶子管理',
      menuRole: ['shootuser','admin'],
      menuLink: '/main/manage/target/list'
    },
    {
      menuName: '规范管理',
      menuRole: ['shootuser','admin'],
      menuLink: '/main/manage/template'
    },
    {
      menuName: '组织看板',
      menuRole: ['admin'],
      menuLink: '/main/manage/dashboard'
    },
    {
      menuName: '个人看板',
      menuRole: ['shootuser', 'admin'],
      menuLink: '/main/user/dashboard'
    },
    {
      menuName: '关于我们',
      menuRole: ['shootuser', 'admin'],
      menuLink: '/main/user/aboutdevteam'
    }
  ];

  constructor() {
  }
}

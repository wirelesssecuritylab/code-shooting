import { Routes } from '@angular/router';
import { LoginComponent } from './login/login.component';
import { ListComponent } from "./list/list.component";
import { MainComponent } from "./main/main.component";
import { ActionTabComponent } from "./action-tab/action-tab.component";
import { RangeListComponent } from './admin/manage/range-list/range-list.component';
import { RangeAddComponent } from './admin/manage/range-add/range-add.component';
import { QueryScoreComponent } from './admin/score/query-score/query-score.component';
import { AuthGuard } from './auth/auth.guard';
import { AdminGuard } from './auth/admin.guard';
import { TargetListComponent } from './admin/target/target-list/target-list.component';
import { TargetAddComponent } from './admin/target/target-add/target-add.component';
import { RandomRangeComponent } from './user/target-range/random-range/random-range.component';
import { TestScoreComponent } from './user/test-score/test-score.component';
import { TemplateManageComponent } from './admin/template/template-manage/template-manage.component';
import { TemplateHistoryComponent } from "./admin/template/template-history/template-history.component";
import { TemplateAddComponent } from "./admin/template/template-add/template-add.component";
import { ShootTargetComponent } from './shoot/shoot-target/shoot-target.component';
import { DashboardComponent } from "./dashboard/dashboard.component";
import { HelpComponent } from "./help/help.component";
import { PersonalcenterComponent } from "./personalcenter/personalcenter.component";
import { AboutdevteamComponent } from './aboutdevteam/aboutdevteam.component';
import { ShootRangeComponent } from './action-tab/shoot-range/shoot-range.component';
import { UserinfoComponent } from "./personalcenter/userinfo/userinfo.component";
import { MyrangeComponent } from "./personalcenter/myrange/myrange.component";
import { FreePracticeComponent } from './personalcenter/freepractice/freepractice.component';
import { MyTargetComponent } from './personalcenter/mytarget/mytarget.component';

export const ROUTES: Routes = [
  {
    path: '',
    children: [
      {
        path: '',
        redirectTo: 'login',
        pathMatch: 'full'
      },
      {
        path: 'login',
        component: LoginComponent,
      },
      {
        path: 'main',
        canActivate: [AuthGuard],
        component: MainComponent,
        children: [
          {
            path: '',
            redirectTo: 'user',
            pathMatch: 'full'
          },
          {
            path: 'user',
            children: [
              {
                path: '',
                redirectTo: 'list',
                pathMatch: 'full'
              },
              {
                path: 'list',
                component: ListComponent,
              },
              {
                path: 'shoot-range',
                component: ShootRangeComponent,
              },
              {
                path: 'test',
                component: RandomRangeComponent,
              },
              {
                path: 'query',
                component: TestScoreComponent,
              },
              {
                path: 'shoot',
                component: ShootTargetComponent,
              },
              {
                path: 'dashboard',
                component: DashboardComponent
              },
              {
                path: 'doc',
                component: HelpComponent
              },
              {
                path: 'personalcenter',
                component: PersonalcenterComponent,
                children: [
                  {
                    path: '',
                    redirectTo: 'userinfo',
                    pathMatch: 'full'
                  },
                  {
                    path: 'userinfo',
                    component: UserinfoComponent
                  },
                  {
                    path: 'freepractice',
                    component: FreePracticeComponent
                  },
                  {
                    path: 'myrange',
                    component: MyrangeComponent
                  },
                  {
                    path: 'mytarget',
                    component: MyTargetComponent
                  }]
              },
              {
                path: 'aboutdevteam',
                component: AboutdevteamComponent
              }
            ]
          },
          {
            path: 'manage',
            children: [
              {
                path: '',
                redirectTo: 'range',
                pathMatch: 'full'
              },
              {
                path: 'range',
                canActivate: [AdminGuard],
                children: [
                  {
                    path: '',
                    redirectTo: 'list',
                    pathMatch: 'full'
                  },
                  {
                    path: 'list',
                    component: RangeListComponent
                  },
                  {
                    path: 'add',
                    component: RangeAddComponent
                  },
                  {
                    path: 'edit',
                    component: RangeAddComponent
                  },
                  {
                    path: 'query',
                    component: QueryScoreComponent
                  },
                ]
              },
              {
                path: 'target',
                children: [
                  {
                    path: '',
                    redirectTo: 'list',
                    pathMatch: 'full'
                  },
                  {
                    path: 'list',
                    component: TargetListComponent
                  },
                  {
                    path: 'add',
                    component: TargetAddComponent
                  },
                  {
                    path: 'edit',
                    component: TargetAddComponent
                  }
                ]
              },
              {
                path: 'template',
                children: [
                  {
                    path: '',
                    redirectTo: 'list',
                    pathMatch: 'full'
                  },
                  {
                    path: 'list',
                    component: TemplateManageComponent
                  },
                  {
                    path: 'add',
                    canActivate: [AdminGuard],
                    component: TemplateAddComponent
                  },
                  {
                    path: 'history',
                    canActivate: [AdminGuard],
                    component: TemplateHistoryComponent
                  }]
              },
              {
                path: 'dashboard',
                canActivate: [AdminGuard],
                component: DashboardComponent
              },
              {
                path: 'aboutdevteam',
                canActivate: [AdminGuard],
                component: AboutdevteamComponent
              }
            ]
          },
        ],
      }
    ]
  }
];

import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { JigsawModule, PopupService } from '@rdkmaster/jigsaw';
import { PlxI18nModule, PlxI18nService, PlxModule } from 'paletx';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LoginComponent } from './login/login.component';
import { NgxQRCodeModule } from '@techiediaries/ngx-qrcode';
import { ListComponent } from './list/list.component';
import { ROUTES } from './app.routes';
import { RouterModule } from '@angular/router';
import { AuthGuard } from './auth/auth.guard';
import { AdminGuard } from './auth/admin.guard';
import { MainComponent } from './main/main.component';
import { HeaderComponent } from './header/header.component';
import { TypeComponent } from './list/type/type.component';
import { NameComponent } from './list/name/name.component';
import { OperationComponent } from './list/operation/operation.component';
import { ActionTabComponent } from './action-tab/action-tab.component';
import { UploadComponent } from './action-tab/upload/upload.component';
import { QueryComponent } from './action-tab/query/query.component';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import {
  ResultRenderComponent,
  ScoreComponent,
} from './action-tab/query/score/score.component';
import { RangeListComponent } from './admin/manage/range-list/range-list.component';
import { RangeOperationComponent } from './admin/manage/range-operation/range-operation.component';
import { RangeAddComponent } from './admin/manage/range-add/range-add.component';
import { HeaderMenuComponent } from './header-menu/header-menu.component';
import { QueryScoreComponent } from './admin/score/query-score/query-score.component';
import { TargetListComponent } from './admin/target/target-list/target-list.component';
import { TargetOperationComponent } from './admin/target/target-operation/target-operation.component';
import { TargetAddComponent } from './admin/target/target-add/target-add.component';
import { TargetTagComponent } from './admin/target/target-tag/target-tag.component';
import { TargetConfigComponent } from './admin/manage/target-config/target-config.component';
import { RandomRangeComponent } from './user/target-range/random-range/random-range.component';
import { TemplateManageComponent } from './admin/template/template-manage/template-manage.component';
import { CommonService } from './shared/service/common.service';
import { TemplateHistoryComponent } from './admin/template/template-history/template-history.component';
import { TemplateAddComponent } from './admin/template/template-add/template-add.component';
import { ShootTargetComponent } from './shoot/shoot-target/shoot-target.component';
import { PracticeModule } from './shoot/components/practice/practice.module';
import { ShootEntranceComponent } from './action-tab/shoot-entrance/shoot-entrance.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { HelpComponent } from './help/help.component';
import { TestScoreComponent } from './user/test-score/test-score.component';
import { MarkdownModule } from 'ngx-markdown';
import { PersonalcenterComponent } from './personalcenter/personalcenter.component';
import { AboutdevteamComponent } from './aboutdevteam/aboutdevteam.component';
import { TabsComponent } from './action-tab/tabs/tabs.component';
import { ProgressBoardComponent } from './action-tab/progress-board/progress-board.component';
import { GradeBoardComponent } from './action-tab/grade-board/grade-board.component';
import { ShootRangeComponent } from './action-tab/shoot-range/shoot-range.component';
import {UserinfoComponent} from "./personalcenter/userinfo/userinfo.component";
import { MyrangeComponent } from './personalcenter/myrange/myrange.component';
import { FreePracticeComponent } from './personalcenter/freepractice/freepractice.component';
import { MyTargetComponent } from './personalcenter/mytarget/mytarget.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    ListComponent,
    MainComponent,
    HeaderComponent,
    TypeComponent,
    NameComponent,
    OperationComponent,
    ActionTabComponent,
    UploadComponent,
    QueryComponent,
    ScoreComponent,
    RangeListComponent,
    RangeOperationComponent,
    RangeAddComponent,
    HeaderMenuComponent,
    QueryScoreComponent,
    TargetListComponent,
    TargetOperationComponent,
    TargetAddComponent,
    TargetTagComponent,
    TargetConfigComponent,
    RandomRangeComponent,
    TemplateManageComponent,
    TemplateHistoryComponent,
    TemplateAddComponent,
    ShootTargetComponent,
    ShootEntranceComponent,
    DashboardComponent,
    HelpComponent,
    UploadComponent,
    TestScoreComponent,
    ResultRenderComponent,
    PersonalcenterComponent,
    AboutdevteamComponent,
    TabsComponent,
    ProgressBoardComponent,
    GradeBoardComponent,
    ShootRangeComponent,
    UserinfoComponent,
    MyrangeComponent,
    FreePracticeComponent,
    MyTargetComponent,
  ],
  entryComponents: [ResultRenderComponent,TemplateAddComponent],
  imports: [
    CommonModule,
    FormsModule,
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    NgxQRCodeModule,
    PlxModule.forRoot(),
    PlxI18nModule.forRoot(),
    HttpClientModule,
    PracticeModule,
    JigsawModule,
    TranslateModule,
    RouterModule.forRoot(ROUTES, { useHash: true }),
    MarkdownModule.forRoot(),
  ],
  providers: [
    AuthGuard,
    AdminGuard,
    CommonService,
    PopupService,
    TranslateService,
    PlxI18nService,
  ],
  bootstrap: [AppComponent],
})
export class AppModule {
  // 代码来自教练集成成果，好像没有用处。暂且注释掉
  // constructor(translateService: TranslateService) {
  //   translateService.setTranslation('zh', {
  //       'get-started': '马上开始',
  //       'give-star': '给 Jigsaw 点个星星'
  //   }, true);
  //   translateService.setTranslation('en', {
  //       'get-started': 'Get started',
  //       'give-star': 'Give us a star on Github.com'
  //   }, true);
  //   const lang: string = translateService.getBrowserLang();
  //   translateService.setDefaultLang(lang);
  //   translateService.use(lang);
  // }
}

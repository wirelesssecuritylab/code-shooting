import { Component, ViewContainerRef } from '@angular/core';
import { PlxMessage } from 'paletx';
@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'Code Shooting';
  constructor(private plxMessageService: PlxMessage,private vRef: ViewContainerRef,){
    this.plxMessageService.setRootViewContainerRef(this.vRef);
    this.plxMessageService.setLanguage('zh_CN');
  }
}

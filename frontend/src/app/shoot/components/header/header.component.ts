import {Component} from '@angular/core';

@Component({
    selector: 'tp-header',
    templateUrl: './header.component.html',
    styleUrls: ['./header.component.scss']
})
export class HeaderComponent {
    public _$userName: string = 'Anna';
}

import {NgModule} from '@angular/core';
import {CommonModule} from "@angular/common";
import {JigsawBadgeModule, JigsawButtonModule} from "@rdkmaster/jigsaw";
import {HeaderComponent} from "./header.component";

@NgModule({
    imports: [CommonModule, JigsawBadgeModule, JigsawButtonModule],
    declarations: [HeaderComponent],
    exports: [HeaderComponent]
})
export class HeaderModule {
}

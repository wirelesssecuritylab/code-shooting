import {NgModule} from '@angular/core';
import {CommonModule} from "@angular/common";
import {PerfectScrollbarModule} from "ngx-perfect-scrollbar";
import {JigsawTreeExtModule} from "@rdkmaster/jigsaw";
import {FilesTreeComponent} from "./files-tree.component";

@NgModule({
    imports: [CommonModule, JigsawTreeExtModule, PerfectScrollbarModule],
    declarations: [FilesTreeComponent],
    exports: [FilesTreeComponent]
})
export class FilesTreeModule {
}

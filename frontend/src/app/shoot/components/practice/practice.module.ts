import {NgModule} from '@angular/core';
import {CommonModule} from "@angular/common";
import {PerfectScrollbarModule} from "ngx-perfect-scrollbar";
import {MarkdownModule} from "ngx-markdown";
import {JigsawModule, PopupService} from "@rdkmaster/jigsaw";
import {PracticeComponent} from "./practice.component";
import {MonacoCodeEditorModule} from "../code-editor/code-editor";
import {EditComponent} from "../comment-tools/edit.component";
import {DOMService} from "../../services/dom-service";
import {AddComponent} from "../comment-tools/add.component";
import {DefectService} from "../../services/defect-service";
import {DefectSelectDialog} from "../comment-tools/defect-select/defect-select.dialog";
import {FilesTreeModule} from "../files-tree/files-tree.module";
import {MarkedInfoService} from "../../services/marked-info-service";
import {CommentComponent} from "../comment-tools/comment.component";
import { AnswerComponent } from '../comment-tools/answer.component';
import {ShowDefectsDialog} from "../show-defects/show-defects.dialog";
import { PlxModule } from 'paletx';

@NgModule({
    imports: [
        CommonModule, MarkdownModule.forChild(), PerfectScrollbarModule, JigsawModule, MonacoCodeEditorModule, FilesTreeModule, PlxModule
    ],
    declarations: [PracticeComponent, EditComponent, AddComponent, DefectSelectDialog, CommentComponent, AnswerComponent, ShowDefectsDialog],
    exports: [PracticeComponent, EditComponent, AddComponent, DefectSelectDialog, CommentComponent, AnswerComponent, ShowDefectsDialog],
    providers: [PopupService, DOMService, DefectService, MarkedInfoService]
})
export class PracticeModule {
}

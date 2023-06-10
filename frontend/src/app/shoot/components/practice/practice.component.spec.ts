import {async, ComponentFixture, TestBed} from '@angular/core/testing';
import {CUSTOM_ELEMENTS_SCHEMA} from "@angular/core";
import {
    JigsawButtonModule,
    JigsawDialogModule,
    JigsawIconModule,
    JigsawNotificationModule,
    JigsawSelectModule,
    JigsawTabsModule,
    PopupService
} from "@rdkmaster/jigsaw";
import {PracticeComponent} from "./practice.component";
import {MonacoCodeEditorModule} from "../code-editor/code-editor";
import {DOMService} from "../../services/dom-service";
import {DefectService} from "../../services/defect-service";
import {TranslateModule} from "@ngx-translate/core";
import {MarkedInfoService} from "../../services/marked-info-service";

describe('PracticeComponent', () => {
    let component: PracticeComponent;
    let fixture: ComponentFixture<PracticeComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            imports: [
                MonacoCodeEditorModule, JigsawButtonModule, JigsawTabsModule, JigsawIconModule, JigsawSelectModule,
                JigsawDialogModule, JigsawNotificationModule, TranslateModule.forRoot()
            ],
            declarations: [PracticeComponent],
            providers: [PopupService, DOMService, DefectService, MarkedInfoService],
            schemas: [CUSTOM_ELEMENTS_SCHEMA]
        }).compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(PracticeComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });

    it('should render submit button', () => {
        fixture.detectChanges();
        const compiled = fixture.nativeElement;
        expect(compiled.querySelector('.jigsaw-button-text').textContent).toContain('提交答案');
    });

    it('should render tp-code', () => {
        fixture.detectChanges();
        const compiled = fixture.nativeElement;
        expect(compiled.querySelector('tp-code')).toBeTruthy();
    });
});

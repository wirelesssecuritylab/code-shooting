import {ComponentFixture, TestBed} from '@angular/core/testing';
import {MonacoCodeEditor, MonacoCodeEditorModule} from './code-editor';
import {PopupService} from "@rdkmaster/jigsaw";

describe('MonacoCodeEditor unit-test suites:', () => {
    let fixture: ComponentFixture<MonacoCodeEditor>;
    let component: MonacoCodeEditor;

    beforeEach(function () {
        TestBed.configureTestingModule({
            imports: [MonacoCodeEditorModule],
            providers: [PopupService]
        });
        fixture = TestBed.createComponent(MonacoCodeEditor);
        component = fixture.componentInstance;
    });

    afterEach(function () {
        fixture = null;
        component = null;
    });

    it('should create the MonacoCodeEditor', function () {
        expect(component).toBeTruthy();
    });

    it('test code property', function () {
        component.code = "content";
        expect(component.code).toBe("content");
        component.code = null;
        expect(component.code).toBe("");
        ////
        component.editor = {
            setValue: function (code) {
            }
        }; //mock editor
        let setValueSpy = spyOn(component.editor, "setValue");
        component.code = "other content";
        expect(setValueSpy).toHaveBeenCalledWith("other content");
    });

    it('test _getLanguageValue method', function () {
        expect(component["_getLanguageValue"]("html")).toBe("html");
        expect(component["_getLanguageValue"]("css")).toBe("css");
        expect(component["_getLanguageValue"]("js")).toBe("javascript");
        expect(component["_getLanguageValue"]("javascript")).toBe("javascript");
        expect(component["_getLanguageValue"]("json")).toBe("json");
        expect(component["_getLanguageValue"]("ts")).toBe("typescript");
        expect(component["_getLanguageValue"]("xml")).toBe("xml");
    });
});

import {TestBed} from '@angular/core/testing';
import {DefectService} from "./defect-service";

describe('DefectService', () => {
    let defectService: DefectService;

    beforeEach(() => {
        TestBed.configureTestingModule({providers: [DefectService]});
        defectService = TestBed.inject(DefectService);
        defectService.targets = [];
    });

    it('test targets', () => {
        expect(defectService.targets.length).toEqual(0);
    });

    it('test toRange', () => {
        expect(typeof defectService.toRange).toBe("function");
        let target;
        expect(defectService.toRange(target)).toEqual({startLineNumber: 0, endLineNumber: 0, startColumn: 0, endColumn: 0});
        target = {StartLineNum: 1, EndLineNum: 2};
        expect(defectService.toRange(target)).toEqual({startLineNumber: 1, endLineNumber: 2, startColumn: 0, endColumn: 0});
        target = {startLineNumber: 1, endLineNumber: 2};
        expect(defectService.toRange(target)).toEqual({startLineNumber: 0, endLineNumber: 0, startColumn: 0, endColumn: 0});
        target = {StartLineNum: 1, StartColNum: 1, EndLineNum: 2, EndColNum: 4};
        expect(defectService.toRange(target)).toEqual({startLineNumber: 1, endLineNumber: 2, startColumn: 1, endColumn: 4});
    });

    it('test getTarget', () => {
        expect(typeof defectService.getTarget).toBe("function");
        let range;
        expect(defectService.getTarget(range, '')).toBeUndefined();
        range = {
            startLineNumber: 3,
            endLineNumber: 5
        };
        expect(defectService.getTarget(range, '')).toBeUndefined();
        const target = {
            StartLineNum: 3,
            EndLineNum: 5,
            FileName: 'test.ts',
            DefectClass: '',
            DefectSubClass: '',
            DefectDescribe: ''
        };
        defectService.targets = [target];
        expect(defectService.getTarget(range, target.FileName)).toEqual(target);
        range = {
            startLineNumber: 1,
            endLineNumber: 4
        };
        expect(defectService.getTarget(range, target.FileName)).toBeUndefined();
        range = {
            startLineNumber: 4,
            endLineNumber: 6
        };
        expect(defectService.getTarget(range, target.FileName)).toBeUndefined();
        range = {
            startLineNumber: 7,
            endLineNumber: 8
        };
        expect(defectService.getTarget(range, target.FileName)).toBeUndefined();
    });

    it('test editTarget', () => {
        expect(typeof defectService.editTarget).toBe("function");
        let target;
        defectService.editTarget(target);
        expect(defectService.targets.length).toEqual(0);

        target = {
            StartLineNum: 3,
            StartColNum: 1,
            EndLineNum: 5,
            EndColNum: 4,
            FileName: 'test.ts',
            DefectClass: '111',
            DefectSubClass: '222',
            DefectDescribe: '333'
        };
        defectService.targets = [target];
        expect(defectService.targets[0]).toEqual(target);

        const newTarget1 = {
            StartLineNum: 3,
            StartColNum: 1,
            EndLineNum: 5,
            EndColNum: 4,
            FileName: 'test.ts',
            DefectClass: '444',
            DefectSubClass: '555',
            DefectDescribe: '666'
        };
        defectService.editTarget(newTarget1);
        expect(defectService.targets[0]).toEqual(newTarget1);

        const newTarget2 = {
            StartLineNum: 6,
            EndLineNum: 7,
            StartColNum: 1,
            EndColNum: 4,
            FileName: 'test.ts',
            DefectClass: '111',
            DefectSubClass: '222',
            DefectDescribe: '333'
        };
        defectService.editTarget(newTarget2);
        expect(defectService.targets[0]).toEqual(newTarget1);
    });

    it('test deleteTarget', () => {
        expect(typeof defectService.deleteTarget).toBe("function");

        let range;
        expect(defectService.deleteTarget(range, '')).toBeUndefined();
        range = {
            startLineNumber: 3,
            endLineNumber: 5
        };
        expect(defectService.deleteTarget(range, '')).toBeUndefined();

        let target = {
            StartLineNum: 3,
            EndLineNum: 5,
            FileName: 'test.ts',
            DefectClass: '111',
            DefectSubClass: '222',
            DefectDescribe: '333'
        };
        defectService.targets = [target];
        expect(defectService.targets.length).toEqual(1);
        expect(defectService.deleteTarget(range, target.FileName)).toEqual(target);
        expect(defectService.targets.length).toEqual(0);

        defectService.targets = [target];
        expect(defectService.targets.length).toEqual(1);
        range = {
            startLineNumber: 6,
            endLineNumber: 7
        };
        expect(defectService.deleteTarget(range, target.FileName)).toBeUndefined()
        expect(defectService.targets.length).toEqual(1);
        expect(defectService.targets[0]).toEqual(target);
    });

    it('test addTarget', () => {
        expect(typeof defectService.addTarget).toBe("function");
        let target, range;
        defectService.addTarget(target, range);
        expect(defectService.targets.length).toEqual(0);

        range = {
            startLineNumber: 3,
            endLineNumber: 5
        };
        target = {
            StartLineNum: 3,
            EndLineNum: 5,
            FileName: 'test.ts',
            DefectClass: '111',
            DefectSubClass: '222',
            DefectDescribe: '333'
        };
        defectService.addTarget(target, range);
        expect(defectService.targets.length).toEqual(1);
        expect(defectService.targets[0]).toEqual(target);
    });
});

import {Utils} from "./utils";
import {PopupEffect} from "@rdkmaster/jigsaw";

describe("src/app/misc/utils.ts test-suite:", function () {
    let utilInstance;

    beforeEach(function () {
        utilInstance = new Utils();
    });

    it('test method isDefined', function () {
        expect(utilInstance.isDefined).toBeUndefined();
        expect(typeof Utils.isDefined).toBe("function");
        expect(Utils.isDefined(undefined)).toBeFalsy();
        expect(Utils.isDefined(null)).toBeFalsy();
        expect(Utils.isDefined(1)).toBe(true);
        expect(Utils.isDefined("a")).toBeTruthy();
        expect(Utils.isDefined([])).toBeTruthy();
        expect(Utils.isDefined({})).toBeTruthy();
        expect(Utils.isDefined(true)).toBeTruthy();
    });

    it('test method getCssValue', function () {
        expect(utilInstance.getCssValue).toBeUndefined();
        expect(typeof Utils.getCssValue).toBe("function");
        expect(Utils.getCssValue(null)).toBeNull();
        expect(Utils.getCssValue(undefined)).toBeNull();
        expect(Utils.getCssValue(12.56)).toBe("12.56px");
        expect(Utils.getCssValue(" 12  ")).toBe("12px");
        expect(Utils.getCssValue(" ab12  ")).toBe("ab12");
        expect(Utils.getCssValue(.02435)).toBe("0.02435px");
        expect(Utils.getCssValue(" 0242.34 ")).toBe("0242.34px");
    });

    it('test method stripBlank', function () {
        expect(utilInstance.stripBlank).toBeUndefined();
        expect(typeof Utils.stripBlank).toBe("function");
        expect(Utils.getCssValue(`
            expect(utilInstance.stripBlank).toBeUndefined();
        `)).toBe("expect(utilInstance.stripBlank).toBeUndefined();");
    });

    it('test method getModalOptions', function () {
        expect(utilInstance.getModalOptions).toBeUndefined();
        expect(typeof Utils.getModalOptions).toBe("function");
        expect(Utils.getModalOptions()).toEqual({
            modal: true,
            showEffect: PopupEffect.bubbleIn,
            hideEffect: PopupEffect.bubbleOut
        });
    });
});

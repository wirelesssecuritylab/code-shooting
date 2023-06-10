import {PopupEffect, PopupOptions} from "@rdkmaster/jigsaw";

export class Utils {
    public static isDefined(value: any): boolean {
        return value !== undefined && value !== null;
    }

    public static getCssValue(value: string | number): string {
        if (!this.isDefined(value)) {
            return null;
        }
        value = typeof value === 'string' ? value.trim() : String(value);
        const match = value ? value.match(/^\s*\d+\.*\d*\s*$/) : null;
        return match ? (value + 'px') : value;
    }

    // 去除默认值的空格
    public static stripBlank(script: string) {
        if (!this.isDefined(script)) {
            return undefined;
        }
        if (typeof script != 'string') {
            return script;
        }
        const match = script.match(/^( *)\S.*$/m);
        const prefixWhiteSpaceLength = match ? match[1].length : 0;
        if (prefixWhiteSpaceLength == 0) {
            return script;
        }
        const mdLines = script.trim().split(/\r?\n/g);
        script = '';
        const reg = new RegExp('^\\s{' + prefixWhiteSpaceLength + '}');
        mdLines.forEach(line => {
            script += line.replace(reg, '') + '\n';
        });
        return script;
    }

    public static getModalOptions(): PopupOptions {
        return {
            modal: true,
            showEffect: PopupEffect.bubbleIn,
            hideEffect: PopupEffect.bubbleOut
        };
    }
}

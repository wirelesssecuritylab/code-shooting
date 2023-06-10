import {ApplicationRef, ComponentFactoryResolver, ComponentRef, EmbeddedViewRef, Injectable, Injector} from '@angular/core';
import {BaseOperate} from "../misc/BaseOperate";
import {TargetInfo} from "../misc/target-types";
import {InitData} from "../misc/types";

@Injectable()
export class DOMService {
    constructor(
        private componentFactoryResolver: ComponentFactoryResolver,
        private applicationRef: ApplicationRef,
        private injector: Injector,
    ) {
    }

    public getComponentRef(component: any, initData: InitData, callback?: Function): ComponentRef<any> {
        const componentRef = this.componentFactoryResolver.resolveComponentFactory(component).create(this.injector);
        const instance = <BaseOperate>componentRef.instance;
        instance.initData = initData;
        instance.answer.subscribe((target: TargetInfo) => {
            if (typeof callback == 'function') {
                callback(target);
            }
        })
        this.applicationRef.attachView(componentRef.hostView);
        return componentRef;
    }

    public getDomElement(componentRef: any) {
        return (componentRef.hostView as EmbeddedViewRef<any>).rootNodes[0] as HTMLElement;
    }

    // todo 清理
    public removeComponentFromBody(componentRef: ComponentRef<any>) {
        this.applicationRef.detachView(componentRef.hostView);
        componentRef.destroy();
    }
}

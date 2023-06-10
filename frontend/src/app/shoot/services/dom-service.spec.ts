import {TestBed} from '@angular/core/testing';
import {DOMService} from "./dom-service";

describe('DOMService', () => {
    let domService: DOMService;

    beforeEach(() => {
        TestBed.configureTestingModule({providers: [DOMService]});
        domService = TestBed.inject(DOMService);
    });
});

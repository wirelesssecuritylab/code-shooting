import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RangeOperationComponent } from './range-operation.component';

describe('RangeOperationComponent', () => {
  let component: RangeOperationComponent;
  let fixture: ComponentFixture<RangeOperationComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ RangeOperationComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(RangeOperationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

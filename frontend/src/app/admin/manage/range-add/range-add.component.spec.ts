import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RangeAddComponent } from './range-add.component';

describe('RangeAddComponent', () => {
  let component: RangeAddComponent;
  let fixture: ComponentFixture<RangeAddComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ RangeAddComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(RangeAddComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
